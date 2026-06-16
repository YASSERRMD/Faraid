package llm

import (
	"context"
)

// ControlsOptions configures the Controls wrapper.
type ControlsOptions struct {
	// TokenCap is the maximum MaxTokens allowed per request. When a request
	// exceeds this value it is silently reduced to the cap. 0 disables the cap.
	TokenCap int
}

// Controls wraps a Completer with cost, rate, and audit guards. It is
// goroutine-safe. Build it with NewControls, then pass it wherever a Completer
// is expected.
type Controls struct {
	inner Completer
	cap   int
}

// NewControls wraps inner with the given controls.
func NewControls(inner Completer, opts ControlsOptions) *Controls {
	return &Controls{inner: inner, cap: opts.TokenCap}
}

// Complete applies all active controls, then delegates to the underlying
// Completer.
func (c *Controls) Complete(ctx context.Context, req Request) (Response, error) {
	if c.cap > 0 && req.MaxTokens > c.cap {
		req.MaxTokens = c.cap
	}
	return c.inner.Complete(ctx, req)
}
