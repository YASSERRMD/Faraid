// Package llm provides the provider-agnostic LLM adapters used only by the
// non-authoritative trial tier (natural-language case parsing and plain
// language explanation drafting).
//
// Nothing under internal/core may import this package. The LLM is never the
// source of a legal result: its output is always validated against the
// deterministic engine, sits behind a feature flag, and defaults off. The
// adapters are introduced in Phase 37.
package llm
