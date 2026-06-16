#!/usr/bin/env bash
# Forbid the em dash (U+2014) and en dash (U+2013) characters anywhere in the
# project source. The build standard requires commas, colons, or rewording
# instead. This check is part of CI and can also be run locally.
set -uo pipefail

matches=$(grep -rlP "[\x{2014}\x{2013}]" \
  --exclude-dir=.git \
  --exclude-dir=node_modules \
  --exclude-dir=temp \
  . 2>/dev/null || true)

if [ -n "$matches" ]; then
  echo "ERROR: em dash or en dash character found in:"
  echo "$matches"
  exit 1
fi

echo "OK: no em dash or en dash characters"
