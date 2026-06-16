package llm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// Request is a single completion request to the trial-tier LLM.
type Request struct {
	System    string
	Prompt    string
	MaxTokens int
}

// Response is the model's reply.
type Response struct {
	Text string
}

// Completer is the provider-agnostic interface the trial tier depends on. No
// calling code references a specific vendor, so there is no lock-in.
type Completer interface {
	Complete(ctx context.Context, req Request) (Response, error)
}

// ErrDisabled is returned by New when the trial tier is turned off, which is
// the default.
var ErrDisabled = errors.New("llm: trial tier is disabled")

// Options configure the trial-tier LLM. Enabled defaults to false so the
// feature is off unless explicitly turned on.
type Options struct {
	Enabled    bool
	Provider   string // "openai" or "anthropic"
	BaseURL    string
	APIKey     string
	Model      string
	MaxTokens  int
	HTTPClient *http.Client // optional; defaults to http.DefaultClient, injectable for tests
}

// New builds a Completer for the configured provider, or returns ErrDisabled
// when the trial tier is off. The result is never the source of a legal
// outcome: it powers only the non-authoritative convenience features.
func New(opts Options) (Completer, error) {
	if !opts.Enabled {
		return nil, ErrDisabled
	}
	switch opts.Provider {
	case "openai":
		return newOpenAI(opts), nil
	case "anthropic":
		return newAnthropic(opts), nil
	default:
		return nil, fmt.Errorf("llm: unknown provider %q", opts.Provider)
	}
}

// httpClient returns the configured client or the default.
func (o Options) httpClient() *http.Client {
	if o.HTTPClient != nil {
		return o.HTTPClient
	}
	return http.DefaultClient
}

// maxTokens returns the configured cap or a safe default.
func (o Options) maxTokens() int {
	if o.MaxTokens > 0 {
		return o.MaxTokens
	}
	return 1024
}
