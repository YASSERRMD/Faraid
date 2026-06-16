package llm

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// stubCompleter records how many times Complete was called and what MaxTokens
// it received.
type stubCompleter struct {
	calls     int
	lastMax   int
	returnErr error
}

func (s *stubCompleter) Complete(_ context.Context, req Request) (Response, error) {
	s.calls++
	s.lastMax = req.MaxTokens
	if s.returnErr != nil {
		return Response{}, s.returnErr
	}
	return Response{Text: "ok"}, nil
}

func TestControlsTokenCap(t *testing.T) {
	stub := &stubCompleter{}
	c := NewControls(stub, ControlsOptions{TokenCap: 100})

	// Request below cap passes through unchanged.
	_, _ = c.Complete(context.Background(), Request{MaxTokens: 50})
	if stub.lastMax != 50 {
		t.Errorf("below cap: got MaxTokens %d, want 50", stub.lastMax)
	}

	// Request above cap is silently reduced.
	_, _ = c.Complete(context.Background(), Request{MaxTokens: 500})
	if stub.lastMax != 100 {
		t.Errorf("above cap: got MaxTokens %d, want 100", stub.lastMax)
	}

	// Exactly at cap is allowed unchanged.
	_, _ = c.Complete(context.Background(), Request{MaxTokens: 100})
	if stub.lastMax != 100 {
		t.Errorf("at cap: got MaxTokens %d, want 100", stub.lastMax)
	}
}

func TestControlsTokenCapZeroMeansNoLimit(t *testing.T) {
	stub := &stubCompleter{}
	c := NewControls(stub, ControlsOptions{TokenCap: 0})
	_, _ = c.Complete(context.Background(), Request{MaxTokens: 99999})
	if stub.lastMax != 99999 {
		t.Errorf("zero cap should not restrict; got %d", stub.lastMax)
	}
}

func TestControlsRateLimit(t *testing.T) {
	stub := &stubCompleter{}
	// Limit to 2 requests per minute with a frozen clock so the window never
	// resets during the test.
	frozen := time.Now()
	c := NewControls(stub, ControlsOptions{RequestsPerMinute: 2})
	// Override the rate limiter's clock with a frozen value.
	c.rate.now = func() time.Time { return frozen }
	c.rate.windowStart = frozen

	if _, err := c.Complete(context.Background(), Request{}); err != nil {
		t.Fatalf("call 1 should succeed: %v", err)
	}
	if _, err := c.Complete(context.Background(), Request{}); err != nil {
		t.Fatalf("call 2 should succeed: %v", err)
	}
	_, err := c.Complete(context.Background(), Request{})
	if !errors.Is(err, ErrRateLimited) {
		t.Errorf("call 3 should be rate-limited; got %v", err)
	}
	if stub.calls != 2 {
		t.Errorf("inner should have been called 2 times, got %d", stub.calls)
	}
}

func TestControlsRateLimitWindowReset(t *testing.T) {
	stub := &stubCompleter{}
	var now time.Time
	c := NewControls(stub, ControlsOptions{RequestsPerMinute: 1})
	c.rate.now = func() time.Time { return now }
	c.rate.windowStart = now

	// First call in the window succeeds.
	if _, err := c.Complete(context.Background(), Request{}); err != nil {
		t.Fatalf("first call: %v", err)
	}
	// Second call in same window is rejected.
	if _, err := c.Complete(context.Background(), Request{}); !errors.Is(err, ErrRateLimited) {
		t.Fatalf("second call in window should be rate-limited")
	}
	// Advance clock past one minute.
	now = now.Add(61 * time.Second)
	// Should succeed again in the new window.
	if _, err := c.Complete(context.Background(), Request{}); err != nil {
		t.Errorf("after window reset should succeed: %v", err)
	}
}

func TestControlsRateLimitZeroMeansNoLimit(t *testing.T) {
	stub := &stubCompleter{}
	c := NewControls(stub, ControlsOptions{RequestsPerMinute: 0})
	for i := range 100 {
		if _, err := c.Complete(context.Background(), Request{}); err != nil {
			t.Fatalf("call %d should succeed with no rate limit: %v", i, err)
		}
	}
}

func TestControlsUsageLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	stub := &stubCompleter{}
	c := NewControls(stub, ControlsOptions{Logger: logger})

	_, _ = c.Complete(context.Background(), Request{Prompt: "hello", MaxTokens: 10})

	logged := buf.String()
	for _, want := range []string{"llm call", "prompt_len=5", "max_tokens=10"} {
		if !strings.Contains(logged, want) {
			t.Errorf("log missing %q; got: %s", want, logged)
		}
	}
}

func TestControlsNoLoggerIsNoOp(t *testing.T) {
	stub := &stubCompleter{}
	c := NewControls(stub, ControlsOptions{}) // Logger is nil
	if _, err := c.Complete(context.Background(), Request{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestControlsLoggingOnError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	stub := &stubCompleter{returnErr: errors.New("boom")}
	c := NewControls(stub, ControlsOptions{Logger: logger})

	_, _ = c.Complete(context.Background(), Request{})
	if !strings.Contains(buf.String(), "boom") {
		t.Errorf("error should appear in log; got: %s", buf.String())
	}
}

func TestControlsKillSwitch(t *testing.T) {
	stub := &stubCompleter{}
	c := NewControls(stub, ControlsOptions{})

	// Before kill: works normally.
	if _, err := c.Complete(context.Background(), Request{}); err != nil {
		t.Fatalf("before kill: %v", err)
	}
	c.Kill()

	// After kill: always ErrDisabled.
	for range 3 {
		_, err := c.Complete(context.Background(), Request{})
		if !errors.Is(err, ErrDisabled) {
			t.Errorf("after kill should return ErrDisabled; got %v", err)
		}
	}
	// Inner was called only once (before kill).
	if stub.calls != 1 {
		t.Errorf("inner calls after kill: got %d, want 1", stub.calls)
	}
}

func TestControlsKillTakesPrecedenceOverRateLimit(t *testing.T) {
	stub := &stubCompleter{}
	frozen := time.Now()
	c := NewControls(stub, ControlsOptions{RequestsPerMinute: 100})
	c.rate.now = func() time.Time { return frozen }
	c.rate.windowStart = frozen

	c.Kill()
	_, err := c.Complete(context.Background(), Request{})
	if !errors.Is(err, ErrDisabled) {
		t.Errorf("kill should take precedence over rate limit; got %v", err)
	}
	if stub.calls != 0 {
		t.Errorf("inner should not be called after kill")
	}
}
