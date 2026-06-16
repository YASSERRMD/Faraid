package llm

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type captured struct {
	url     string
	headers http.Header
	body    []byte
}

// mockClient returns an http.Client whose transport records the request and
// replies with the given status and body.
func mockClient(status int, body string, cap *captured) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if cap != nil {
			cap.url = r.URL.String()
			cap.headers = r.Header.Clone()
			if r.Body != nil {
				cap.body, _ = io.ReadAll(r.Body)
			}
		}
		return &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})}
}

func TestNewDisabledAndUnknown(t *testing.T) {
	if _, err := New(Options{Enabled: false}); !errors.Is(err, ErrDisabled) {
		t.Errorf("disabled should return ErrDisabled, got %v", err)
	}
	if _, err := New(Options{Enabled: true, Provider: "cohere"}); err == nil {
		t.Error("unknown provider should error")
	}
}

func TestOpenAIComplete(t *testing.T) {
	cap := &captured{}
	c, err := New(Options{
		Enabled:    true,
		Provider:   "openai",
		BaseURL:    "https://api.example.com/v1",
		APIKey:     "sk-test",
		Model:      "gpt-test",
		HTTPClient: mockClient(200, `{"choices":[{"message":{"role":"assistant","content":"hi there"}}]}`, cap),
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Complete(context.Background(), Request{System: "be brief", Prompt: "hello", MaxTokens: 50})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text != "hi there" {
		t.Errorf("text = %q, want 'hi there'", resp.Text)
	}
	if cap.url != "https://api.example.com/v1/chat/completions" {
		t.Errorf("url = %q", cap.url)
	}
	if cap.headers.Get("Authorization") != "Bearer sk-test" {
		t.Errorf("missing bearer auth: %q", cap.headers.Get("Authorization"))
	}
	var sent openAIRequest
	if err := json.Unmarshal(cap.body, &sent); err != nil {
		t.Fatal(err)
	}
	if sent.MaxTokens != 50 || len(sent.Messages) != 2 || sent.Messages[0].Role != "system" {
		t.Errorf("request body wrong: %+v", sent)
	}
}

func TestAnthropicComplete(t *testing.T) {
	cap := &captured{}
	c, err := New(Options{
		Enabled:    true,
		Provider:   "anthropic",
		BaseURL:    "https://api.anthropic.com",
		APIKey:     "sk-ant",
		HTTPClient: mockClient(200, `{"content":[{"type":"text","text":"hello world"}]}`, cap),
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Complete(context.Background(), Request{System: "ground truth", Prompt: "explain"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text != "hello world" {
		t.Errorf("text = %q", resp.Text)
	}
	if cap.url != "https://api.anthropic.com/v1/messages" {
		t.Errorf("url = %q", cap.url)
	}
	if cap.headers.Get("x-api-key") != "sk-ant" || cap.headers.Get("anthropic-version") != anthropicVersion {
		t.Errorf("missing anthropic headers: %v", cap.headers)
	}
	var sent anthropicRequest
	if err := json.Unmarshal(cap.body, &sent); err != nil {
		t.Fatal(err)
	}
	if sent.Model != defaultAnthropicModel || sent.System != "ground truth" || sent.MaxTokens != 1024 {
		t.Errorf("request body wrong: %+v", sent)
	}
}

func TestProviderErrors(t *testing.T) {
	ctx := context.Background()

	// Non-2xx status.
	c, _ := New(Options{Enabled: true, Provider: "openai", HTTPClient: mockClient(500, "boom", nil)})
	if _, err := c.Complete(ctx, Request{Prompt: "x"}); err == nil {
		t.Error("non-2xx should error")
	}

	// OpenAI with no choices.
	c, _ = New(Options{Enabled: true, Provider: "openai", HTTPClient: mockClient(200, `{"choices":[]}`, nil)})
	if _, err := c.Complete(ctx, Request{Prompt: "x"}); err == nil {
		t.Error("empty choices should error")
	}

	// Anthropic with no text block.
	c, _ = New(Options{Enabled: true, Provider: "anthropic", HTTPClient: mockClient(200, `{"content":[]}`, nil)})
	if _, err := c.Complete(ctx, Request{Prompt: "x"}); err == nil {
		t.Error("no text block should error")
	}

	// Malformed JSON.
	c, _ = New(Options{Enabled: true, Provider: "anthropic", HTTPClient: mockClient(200, `{bad`, nil)})
	if _, err := c.Complete(ctx, Request{Prompt: "x"}); err == nil {
		t.Error("bad json should error")
	}
}

func TestTransportError(t *testing.T) {
	failing := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("network down")
	})}
	c, _ := New(Options{Enabled: true, Provider: "openai", HTTPClient: failing})
	if _, err := c.Complete(context.Background(), Request{Prompt: "x"}); err == nil {
		t.Error("transport error should propagate")
	}
}

func TestRequestMaxTokensPreferred(t *testing.T) {
	cap := &captured{}
	c, _ := New(Options{Enabled: true, Provider: "anthropic", MaxTokens: 100, HTTPClient: mockClient(200, `{"content":[{"type":"text","text":"x"}]}`, cap)})
	if _, err := c.Complete(context.Background(), Request{Prompt: "x", MaxTokens: 42}); err != nil {
		t.Fatal(err)
	}
	var sent anthropicRequest
	_ = json.Unmarshal(cap.body, &sent)
	if sent.MaxTokens != 42 {
		t.Errorf("per-request max tokens should win, got %d", sent.MaxTokens)
	}
}
