# Changelog

All notable changes to Faraid are documented here.

## v1.0.0 (2026-06-16)

### Summary

First production release of Faraid: a deterministic Islamic inheritance
calculator covering all four Sunni schools (Hanafi, Maliki, Shafii, Hanbali).

### Core engine (Phases 1-20)

- Solver for the full quranic fixed-share table (fard), residuary (asaba),
  total exclusion (hajb hirman), share reduction (hajb nuqsan), awl (pro-rata
  reduction), and radd (return of surplus)
- School-parameterised divergence: grandfather-with-siblings (Abu Hanifa vs Zayd
  ibn Thabit view), mushtaraka (al-Himariyya), and distant-kindred routing
- Rational arithmetic with big.Int denominators; every fraction is exact
- Estate pre-distribution: funeral, debts, and wasiyyah (one-third cap enforced)
- Step-by-step derivation trail attached to every result

### Heir coverage (Phases 5-15)

- Spouses (husband, up to four wives sharing equally)
- Lineal descendants: son, daughter, son's son, son's daughter (with asaba
  tabsir and completion-to-two-thirds rules)
- Lineal ascendants: father, mother, paternal grandfather, paternal and maternal
  grandmothers
- Collaterals: full, consanguine, and uterine brothers and sisters, and the
  sons-of-brothers tier
- Tasib (asaba ma'a ghayrihi): full and consanguine sisters with female descendants

### Special cases (Phases 16-20)

- Gharrawain (Umariyyatan): mother takes one third of residue, not estate
- Mushtaraka (al-Himariyya): full brothers join uterine siblings in the shared
  one-third under Maliki and Shafii
- Jadd wa ikhwa: grandfather-with-siblings under all four school views
- Distant kindred (dhawu al-arham): structured distribution under Hanafi and
  Hanbali; excluded under Maliki and Shafii

### API and persistence (Phases 21-35)

- Go HTTP server with chi router under `/api/v1/`
- Endpoints: `POST /solve`, `POST /compare`, `GET/POST/DELETE /cases`,
  `POST /export` (PDF and CSV), `GET /readyz`, `GET /debug/vars`
- OpenAPI 3.1 specification at `openapi/faraid.yaml`
- In-memory store (default) and PostgreSQL store via pgx v5
- Auto-schema migration (`CREATE TABLE IF NOT EXISTS`) on startup

### Security and reliability (Phase 48)

- Secure response headers middleware (CSP, X-Frame-Options, Referrer-Policy)
- Per-IP fixed-window rate limiter (120 req/min default)
- Request body size limit (512 KB)
- Graceful shutdown with 10-second drain on SIGINT/SIGTERM

### Observability (Phase 49)

- expvar counters: `faraid_solve_total`, `faraid_solve_errors`,
  `faraid_compare_total`, `faraid_explain_total`, `faraid_parse_total`
- Readiness probe at `GET /api/v1/readyz` (calls `store.Ping`)
- nginx access log forwarded to stdout

### SvelteKit frontend (Phases 25-47)

- Svelte 5 runes mode; Tailwind CSS 4 logical properties for RTL
- Heir entry form with bilingual (Arabic/English) labels
- Derivation panel with collapsible stage-by-stage audit trail
- Result table with per-heir fraction, parts, and amount
- PDF download and link-sharing via saved cases
- Accessibility: skip-link, ARIA roles, aria-live regions, descriptive labels

### LLM features (Phases 39, 46) - trial tier only

- `POST /explain`: narrative explanation of a result (gated by feature flag)
- `POST /parse`: natural-language case entry parsed to structured heirs
- Both features disabled by default (`PUBLIC_LLM_ENABLED` unset)
- LLM results are never the source of a legal ruling; consistency guard rejects
  any explanation that contradicts the computed fractions

### Deployment (Phase 49)

- Multi-stage Dockerfiles for Go API and SvelteKit frontend
- docker-compose.yml: db (postgres:17), api, web, proxy (nginx:1.27)
- `@sveltejs/adapter-node` for production SvelteKit builds

### Documentation (Phase 50)

- `docs/deployment.md`: docker compose quick-start and environment reference
- `docs/user-guide.md`: determinism guarantee, how-to, LLM limitations
- `docs/scholar-review.md`: generated scholar review pack (fixed-share rules,
  blocking rules, nuqsan rules, divergence matrix, known limitations)

### Known limitations

See `docs/scholar-review.md` section 7 for the full list. Key items:

- Munasakha (sequential deaths before distribution) is not supported
- Dhawu al-arham distribution covers common patterns only
- Bequests UI does not expose the `HeirsConsentToExcessBequest` flag
