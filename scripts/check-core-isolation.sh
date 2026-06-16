#!/usr/bin/env bash
# Enforce the provider-agnostic invariant: the deterministic legal core under
# internal/core must never depend on the LLM layer or any provider SDK, either
# directly or transitively. This keeps every legal result free of any
# nondeterministic dependency. Extend the forbidden pattern as provider SDKs
# are introduced under internal/llm.
set -euo pipefail

forbidden='github.com/YASSERRMD/Faraid/internal/llm'

deps=$(go list -deps ./internal/core/... 2>/dev/null || true)

if echo "$deps" | grep -qE "$forbidden"; then
  echo "ERROR: internal/core depends on a forbidden package:"
  echo "$deps" | grep -E "$forbidden"
  exit 1
fi

echo "OK: core isolation invariant holds (no llm or provider SDK dependency)"
