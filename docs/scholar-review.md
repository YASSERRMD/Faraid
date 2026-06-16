# Faraid v1.0.0 Scholar Review Pack

Generated from live rule data by `cmd/review`. Do not edit by hand.
Corrections belong in the rule source files, not here.

---

## 1. Fixed-Share Rules (Fard)

Rules sourced from `internal/core/rules/`. Each row reflects one `FixedShareRule` struct.

| Heir | Share | Condition | Reference |
|------|-------|-----------|-----------|
| husband | 1/2 | no inheriting descendant | Quran 4:12 |
| husband | 1/4 | an inheriting descendant is present | Quran 4:12 |
| wife | 1/4 | no inheriting descendant | Quran 4:12 |
| wife | 1/8 | an inheriting descendant is present | Quran 4:12 |
| father | 1/6 | an inheriting descendant is present | Quran 4:11 |
| mother | 1/6 | an inheriting descendant is present, or two or more siblings | Quran 4:11 |
| mother | 1/3 | no inheriting descendant and fewer than two siblings | Quran 4:11 |
| paternal grandfather | 1/6 | an inheriting descendant is present | Quran 4:11 by analogy to the father |
| paternal grandmother | 1/6 | grandmother present and not blocked | Sunnah (grandmother one sixth) |
| maternal grandmother | 1/6 | grandmother present and not blocked | Sunnah (grandmother one sixth) |
| daughter | 1/2 | one daughter and no son | Quran 4:11 |
| daughter | 2/3 | two or more daughters and no son | Quran 4:11 |
| son's daughter | 1/2 | one son's daughter; no son, no daughter, no son's son | Quran 4:11 by analogy to the daughter |
| son's daughter | 2/3 | two or more son's daughters; no son, no daughter, no son's son | Quran 4:11 by analogy to the daughter |
| son's daughter | 1/6 | exactly one daughter present, completing the group to two thirds; no son, no son's son | Sunnah, the ruling of Ibn Mas'ud (completion to two thirds) |
| full sister | 1/2 | one full sister; no descendant, father, or full brother | Quran 4:176 |
| full sister | 2/3 | two or more full sisters; no descendant, father, or full brother | Quran 4:176 |
| consanguine sister | 1/2 | one consanguine sister; no full sister and no blocking heir | Quran 4:176 by analogy to the full sister |
| consanguine sister | 2/3 | two or more consanguine sisters; no full sister and no blocking heir | Quran 4:176 by analogy to the full sister |
| consanguine sister | 1/6 | exactly one full sister present, completing the group to two thirds | Sunnah, completion to two thirds (takmila) |
| uterine brother | 1/6 | a single uterine sibling; no descendant, father, or grandfather | Quran 4:12 |
| uterine brother | 1/3 | two or more uterine siblings sharing equally; no descendant, father, or grandfather | Quran 4:12 |
| uterine sister | 1/6 | a single uterine sibling; no descendant, father, or grandfather | Quran 4:12 |
| uterine sister | 1/3 | two or more uterine siblings sharing equally; no descendant, father, or grandfather | Quran 4:12 |

---

## 2. Total Exclusion Rules (Hajb Hirman)

Rules sourced from `registerBlock` calls across all per-heir files.

| Blocked Heir | Blocked By | Condition | Reference |
|--------------|------------|-----------|-----------|
| paternal grandfather | father | father is present | the father excludes his ascendants and the siblings |
| paternal grandmother | father | father is present | the father excludes his ascendants and the siblings |
| paternal grandmother | mother | mother is present | the mother excludes the grandmothers |
| maternal grandmother | mother | mother is present | the mother excludes the grandmothers |
| son's son | son | son is present | a son excludes lower descendants and all siblings |
| son's daughter | son | son is present | a son excludes lower descendants and all siblings |
| son's daughter | daughter, son's son | two or more daughters and no son's son | the daughters take the whole two thirds |
| full brother | father | father is present | the father excludes his ascendants and the siblings |
| full brother | son | son is present | a son excludes lower descendants and all siblings |
| full brother | son's son | son's son is present | a son's son excludes the siblings |
| full sister | father | father is present | the father excludes his ascendants and the siblings |
| full sister | son | son is present | a son excludes lower descendants and all siblings |
| full sister | son's son | son's son is present | a son's son excludes the siblings |
| consanguine brother | father | father is present | the father excludes his ascendants and the siblings |
| consanguine brother | son | son is present | a son excludes lower descendants and all siblings |
| consanguine brother | son's son | son's son is present | a son's son excludes the siblings |
| consanguine brother | full brother | full brother is present | a full brother excludes the consanguine siblings |
| consanguine brother | full sister | a full sister inherits as residuary with a female descendant | asaba ma'a ghayrihi takes the residue |
| consanguine sister | father | father is present | the father excludes his ascendants and the siblings |
| consanguine sister | son | son is present | a son excludes lower descendants and all siblings |
| consanguine sister | son's son | son's son is present | a son's son excludes the siblings |
| consanguine sister | full brother | full brother is present | a full brother excludes the consanguine siblings |
| consanguine sister | full sister | two or more full sisters and no consanguine brother | the full sisters take the whole two thirds |
| consanguine sister | full sister | a full sister inherits as residuary with a female descendant | asaba ma'a ghayrihi takes the residue |
| uterine brother | father | father is present | the father excludes his ascendants and the siblings |
| uterine brother | paternal grandfather | paternal grandfather is present | the grandfather excludes uterine siblings |
| uterine brother | son | son is present | a son excludes lower descendants and all siblings |
| uterine brother | son's son | son's son is present | a son's son excludes the siblings |
| uterine brother | daughter | daughter is present | an inheriting descendant excludes uterine siblings |
| uterine brother | son's daughter | son's daughter is present | an inheriting descendant excludes uterine siblings |
| uterine sister | father | father is present | the father excludes his ascendants and the siblings |
| uterine sister | paternal grandfather | paternal grandfather is present | the grandfather excludes uterine siblings |
| uterine sister | son | son is present | a son excludes lower descendants and all siblings |
| uterine sister | son's son | son's son is present | a son's son excludes the siblings |
| uterine sister | daughter | daughter is present | an inheriting descendant excludes uterine siblings |
| uterine sister | son's daughter | son's daughter is present | an inheriting descendant excludes uterine siblings |

---

## 3. Share-Reduction Rules (Hajb Nuqsan)

Rules sourced from `hajbNuqsan` in `internal/core/rules/nuqsan.go`.

| Reduced Heir | Reduced By | Condition | Reference |
|--------------|------------|-----------|-----------|
| husband | son, daughter, son's son, son's daughter | an inheriting descendant lowers the husband from one half to one quarter | Quran 4:12 |
| wife | son, daughter, son's son, son's daughter | an inheriting descendant lowers the wife from one quarter to one eighth | Quran 4:12 |
| mother | son, daughter, son's son, son's daughter, full brother, full sister, consanguine brother, consanguine sister, uterine brother, uterine sister | an inheriting descendant, or two or more siblings, lower the mother from one third to one sixth | Quran 4:11 |

---

## 4. School Divergence Parameters

The four parameters that differ by school. All other rules are shared.

| Parameter | Hanafi | Maliki | Shafii | Hanbali |
|-----------|--------|--------|--------|---------|
| Grandfather with brothers | Abu Hanifa: grandfather excludes brothers | Zayd: grandfather and brothers share | Zayd: grandfather and brothers share | Zayd: grandfather and brothers share |
| Mushtaraka (al-Himariyya) | Full brothers do not share uterine third | Full brothers share uterine third | Full brothers share uterine third | Full brothers do not share uterine third |
| Radd to spouse | Spouse excluded from radd | Spouse excluded from radd | Spouse excluded from radd | Spouse excluded from radd |
| Distant kindred (dhawu al-arham) | Inherit via structured route | Excluded; residue to treasury | Excluded; residue to treasury | Inherit via structured route |

---

## 5. Divergence Matrix

Canonical cases run through all four schools by the live solver.
Rows with identical results across all schools are marked **agree**.

### Husband only

Deceased: **female**

All four schools **agree**: husband: 1/2

### Wife only

Deceased: **male**

All four schools **agree**: wife: 1/4

### Son and daughter

Deceased: **male**

All four schools **agree**: son: 2/3; daughter: 1/3

### Husband, mother, father

Deceased: **female**

All four schools **agree**: husband: 1/2; father: 1/3; mother: 1/6

### Wife, mother, father (gharrawain)

Deceased: **male**

All four schools **agree**: wife: 1/4; father: 1/2; mother: 1/4

### Wife and daughter (radd)

Deceased: **male**

All four schools **agree**: wife: 1/8; daughter: 7/8

### Husband and daughter (radd)

Deceased: **female**

All four schools **agree**: husband: 1/4; daughter: 3/4

### Grandfather and full brothers (jadd)

Deceased: **male**

| School | Shares | Flags |
|--------|--------|-------|
| Hanafi | paternal grandfather: 1 |  |
| Maliki | paternal grandfather: 1/3; full brother: 2/3 |  |
| Shafii | paternal grandfather: 1/3; full brother: 2/3 |  |
| Hanbali | paternal grandfather: 1/3; full brother: 2/3 |  |

### Mushtaraka: husband, mother, 2 uterine brothers, 2 full brothers

Deceased: **female**

| School | Shares | Flags |
|--------|--------|-------|
| Hanafi | husband: 1/2; mother: 1/6; uterine brother: 1/3 |  |
| Maliki | husband: 1/2; mother: 1/6; full brother: 1/6; uterine brother: 1/6 |  |
| Shafii | husband: 1/2; mother: 1/6; full brother: 1/6; uterine brother: 1/6 |  |
| Hanbali | husband: 1/2; mother: 1/6; uterine brother: 1/3 |  |

### Daughter and son's daughter

Deceased: **male**

All four schools **agree**: daughter: 3/4; son's daughter: 1/4

### Two daughters and son's son

Deceased: **male**

All four schools **agree**: daughter: 2/3; son's son: 1/3

### Full sisters with daughters (tasib)

Deceased: **male**

All four schools **agree**: daughter: 2/3; full sister: 1/3

### Consanguine sister with full sister blocked

Deceased: **male**

All four schools **agree**: full sister: 3/4; consanguine sister: 1/4

### Mother with two siblings (sibling count)

Deceased: **male**

All four schools **agree**: mother: 1/6; full brother: 5/6

---

## 6. Classical Test Coverage

The table below summarises the `testdata/classical/` corpus used as the
regression suite. Test files are run by `go test ./internal/core/solver/...`.

| File | Focus |
|------|-------|
| `awl.json` | Proportional reduction (awl) when fixed shares exceed the estate |
| `blocking.json` | Hajb hirman: total exclusion chains |
| `descendants.json` | Sons, daughters, sons' sons, sons' daughters |
| `jadd.json` | Grandfather with siblings (all four school views) |
| `mixed.json` | Multi-category configurations |
| `parents.json` | Father, mother, grandparents |
| `radd.json` | Return (radd) when fixed shares do not exhaust the estate |
| `residuary.json` | Asaba bi'l-nafs and asaba bi'l-ghair |
| `siblings.json` | Full, consanguine, and uterine siblings |
| `special.json` | Gharrawain, mushtaraka, and other named special cases |
| `spouses.json` | Husband and wife in various configurations |

---

## 7. Known Limitations

1. **Dhawu al-arham (distant kindred)** are handled by a separate `DistributeDhawuArham` function and are not included in the main `Solve` result. The school flag `DhawuArham` is encoded and tested, but the distribution logic covers only the most common patterns.

2. **Munasakha (cases with deceased heirs)** are not supported. If an heir dies before distribution is complete, the case must be decomposed into sequential sub-cases manually.

3. **Bequests exceeding one third** are silently capped to one third unless `HeirsConsentToExcessBequest` is set. There is no UI control for this field in v1.0.

4. **Non-monetary assets** (land, livestock, jewellery) are not modelled. All amounts are treated as fungible integers in the smallest currency unit.

5. **Waqf and conditional estates** are outside scope. The engine assumes a freely distributable estate.

6. **The LLM explain/parse features** are trial-tier only, gated behind `PUBLIC_LLM_ENABLED=true`, and are never the source of a legal result. Every LLM response carries an explicit disclaimer. The consistency guard rejects any explanation that contradicts the computed fractions.

7. **`NeedsReview` flag**: the solver sets this flag on results it cannot fully resolve under the current rules (for example, some edge cases in the grandfather-with-siblings calculation). Such results should be audited manually.

---

*End of scholar review pack.*
