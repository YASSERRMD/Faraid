# Faraid

Faraid is an enterprise-grade Islamic inheritance (ilm al-faraid) calculation
system. It computes the distribution of an estate among heirs according to the
rules of the four Sunni schools (Hanafi, Maliki, Shafi'i, Hanbali).

## Design principles

1. **Determinism first.** Every share calculation, blocking rule, awl
   (proportional reduction), and radd (proportional return) is computed by
   pure, deterministic Go code using exact rational arithmetic. There is no
   floating point and no LLM in the legal core.
2. **Exact rational arithmetic.** Shares are fractions of the estate, built on
   `math/big.Rat`, never decimals.
3. **Full audit trail.** Every result emits a structured, step-by-step
   derivation: which heirs are present, who blocks whom, base share
   assignment, asl al-mas'ala, tashih, awl or radd if applied, and the final
   per-heir fraction and amount.
4. **Madhhab-aware.** Rules that differ by school are encoded as data, not
   hardcoded branches.
5. **LLM is trial-tier only.** The LLM is used solely for non-authoritative
   convenience features (natural-language case entry, plain-language
   explanation drafting). Its output is always validated against the
   deterministic engine, sits behind a feature flag, defaults off, and is
   clearly labeled experimental. It is never the source of a legal result.

## Repository layout

```
cmd/faraidd/        server entrypoint
internal/
  core/             deterministic legal engine (no I/O, no LLM)
    rational/       math/big.Rat helpers
    heir/           heir types and relationship model
    rules/          share tables, blocking lattice, per-madhhab data
    solver/         asl, awl, radd, tashih, special cases
    derivation/     structured audit-trail emitter
  api/              HTTP handlers, OpenAPI
  store/            postgres, migrations
  llm/              provider-agnostic adapters (trial tier)
  config/
migrations/
openapi/            API contract
web/                SvelteKit app
testdata/classical/ golden worked examples from fiqh texts
docs/
```

## Status

Under active construction. This project is being built in phases; see the
build plan for the current roadmap.

## Build

```
go build ./...
go test ./...
```

## License

See [LICENSE](LICENSE).
