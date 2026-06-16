package llm

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrRateLimited is returned by Controls.Complete when the caller has exceeded
// the configured requests-per-minute limit.
var ErrRateLimited = errors.New("llm: rate limit exceeded")

// ControlsOptions configures the Controls wrapper.
type ControlsOptions struct {
	// TokenCap is the maximum MaxTokens allowed per request. When a request
	// exceeds this value it is silently reduced to the cap. 0 disables the cap.
	TokenCap int
	// RequestsPerMinute limits LLM calls within a 60-second window. 0 disables
	// the rate limit.
	RequestsPerMinute int
}

// Controls wraps a Completer with cost, rate, and audit guards. It is
// goroutine-safe. Build it with NewControls, then pass it wherever a Completer
// is expected.
type Controls struct {
	inner Completer
	cap   int
	rate  *fixedWindowLimiter
}

// NewControls wraps inner with the given controls.
func NewControls(inner Completer, opts ControlsOptions) *Controls {
	c := &Controls{inner: inner, cap: opts.TokenCap}
	if opts.RequestsPerMinute > 0 {
		c.rate = newFixedWindowLimiter(opts.RequestsPerMinute, time.Now)
	}
	return c
}

// Complete applies all active controls, then delegates to the underlying
// Completer.
func (c *Controls) Complete(ctx context.Context, req Request) (Response, error) {
	if c.cap > 0 && req.MaxTokens > c.cap {
		req.MaxTokens = c.cap
	}
	if c.rate != nil {
		if err := c.rate.Allow(); err != nil {
			return Response{}, err
		}
	}
	return c.inner.Complete(ctx, req)
}

// fixedWindowLimiter enforces a count limit within a rolling 60-second window.
type fixedWindowLimiter struct {
	mu          sync.Mutex
	limit       int
	count       int
	windowStart time.Time
	now         func() time.Time
}

func newFixedWindowLimiter(perMinute int, now func() time.Time) *fixedWindowLimiter {
	return &fixedWindowLimiter{limit: perMinute, windowStart: now(), now: now}
}

func (l *fixedWindowLimiter) Allow() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	n := l.now()
	if n.Sub(l.windowStart) >= time.Minute {
		l.count = 0
		l.windowStart = n
	}
	if l.count >= l.limit {
		return ErrRateLimited
	}
	l.count++
	return nil
}
