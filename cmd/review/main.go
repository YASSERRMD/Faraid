// cmd/review generates docs/scholar-review.md from live rule data.
//
// Run with: go run ./cmd/review
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/YASSERRMD/Faraid/internal/core/estate"
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rules"
	"github.com/YASSERRMD/Faraid/internal/core/solver"
)

func main() {
	f, err := os.Create("docs/scholar-review.md")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	w := func(format string, args ...any) {
		fmt.Fprintf(f, format+"\n", args...)
	}

	w("# Faraid v1.0.0 Scholar Review Pack")
	w("")
	w("Generated from live rule data by `cmd/review`. Do not edit by hand.")
	w("Corrections belong in the rule source files, not here.")
	w("")
	w("---")
	w("")
	w("## 1. Fixed-Share Rules (Fard)")
	w("")
	w("Rules sourced from `internal/core/rules/`. Each row reflects one `FixedShareRule` struct.")
	w("")
	w("| Heir | Share | Condition | Reference |")
	w("|------|-------|-----------|-----------|")

	for _, r := range heir.AllRelations() {
		rs := rules.Rules(r)
		for _, rule := range rs {
			w("| %s | %s | %s | %s |",
				r.String(), rule.Share.String(), rule.Condition, rule.Reference)
		}
	}

	w("")
	w("---")
	w("")
	w("## 2. Total Exclusion Rules (Hajb Hirman)")
	w("")
	w("Rules sourced from `registerBlock` calls across all per-heir files.")
	w("")
	w("| Blocked Heir | Blocked By | Condition | Reference |")
	w("|--------------|------------|-----------|-----------|")

	for _, r := range heir.AllRelations() {
		for _, rule := range rules.HirmanRules(r) {
			blockers := make([]string, len(rule.Blockers))
			for i, b := range rule.Blockers {
				blockers[i] = b.String()
			}
			w("| %s | %s | %s | %s |",
				r.String(), strings.Join(blockers, ", "), rule.Condition, rule.Reference)
		}
	}

	w("")
	w("---")
	w("")
	w("## 3. Share-Reduction Rules (Hajb Nuqsan)")
	w("")
	w("Rules sourced from `hajbNuqsan` in `internal/core/rules/nuqsan.go`.")
	w("")
	w("| Reduced Heir | Reduced By | Condition | Reference |")
	w("|--------------|------------|-----------|-----------|")

	for _, rule := range rules.NuqsanRules() {
		reducers := make([]string, len(rule.Reducers))
		for i, b := range rule.Reducers {
			reducers[i] = b.String()
		}
		w("| %s | %s | %s | %s |",
			rule.Reduced.String(), strings.Join(reducers, ", "), rule.Condition, rule.Reference)
	}

	w("")
	w("---")
	w("")
	w("## 4. School Divergence Parameters")
	w("")
	w("The four parameters that differ by school. All other rules are shared.")
	w("")
	w("| Parameter | Hanafi | Maliki | Shafii | Hanbali |")
	w("|-----------|--------|--------|--------|---------|")
	w("| Grandfather with brothers | Abu Hanifa: grandfather excludes brothers | Zayd: grandfather and brothers share | Zayd: grandfather and brothers share | Zayd: grandfather and brothers share |")
	w("| Mushtaraka (al-Himariyya) | Full brothers do not share uterine third | Full brothers share uterine third | Full brothers share uterine third | Full brothers do not share uterine third |")
	w("| Radd to spouse | Spouse excluded from radd | Spouse excluded from radd | Spouse excluded from radd | Spouse excluded from radd |")
	w("| Distant kindred (dhawu al-arham) | Inherit via structured route | Excluded; residue to treasury | Excluded; residue to treasury | Inherit via structured route |")

	w("")
	w("---")
	w("")
	w("## 5. Divergence Matrix")
	w("")
	w("Canonical cases run through all four schools by the live solver.")
	w("Rows with identical results across all schools are marked **agree**.")
	w("")

	type canonicalCase struct {
		name       string
		sex        heir.Sex
		heirCounts map[heir.Relation]int
	}

	cases := []canonicalCase{
		{
			name:       "Husband only",
			sex:        heir.Female,
			heirCounts: map[heir.Relation]int{heir.Husband: 1},
		},
		{
			name:       "Wife only",
			sex:        heir.Male,
			heirCounts: map[heir.Relation]int{heir.Wife: 1},
		},
		{
			name:       "Son and daughter",
			sex:        heir.Male,
			heirCounts: map[heir.Relation]int{heir.Son: 1, heir.Daughter: 1},
		},
		{
			name:       "Husband, mother, father",
			sex:        heir.Female,
			heirCounts: map[heir.Relation]int{heir.Husband: 1, heir.Mother: 1, heir.Father: 1},
		},
		{
			name:       "Wife, mother, father (gharrawain)",
			sex:        heir.Male,
			heirCounts: map[heir.Relation]int{heir.Wife: 1, heir.Mother: 1, heir.Father: 1},
		},
		{
			name:       "Wife and daughter (radd)",
			sex:        heir.Male,
			heirCounts: map[heir.Relation]int{heir.Wife: 1, heir.Daughter: 1},
		},
		{
			name:       "Husband and daughter (radd)",
			sex:        heir.Female,
			heirCounts: map[heir.Relation]int{heir.Husband: 1, heir.Daughter: 1},
		},
		{
			name:       "Grandfather and full brothers (jadd)",
			sex:        heir.Male,
			heirCounts: map[heir.Relation]int{heir.PaternalGrandfather: 1, heir.FullBrother: 2},
		},
		{
			name: "Mushtaraka: husband, mother, 2 uterine brothers, 2 full brothers",
			sex:  heir.Female,
			heirCounts: map[heir.Relation]int{
				heir.Husband:        1,
				heir.Mother:         1,
				heir.UterineBrother: 2,
				heir.FullBrother:    2,
			},
		},
		{
			name:       "Daughter and son's daughter",
			sex:        heir.Male,
			heirCounts: map[heir.Relation]int{heir.Daughter: 1, heir.SonsDaughter: 1},
		},
		{
			name:       "Two daughters and son's son",
			sex:        heir.Male,
			heirCounts: map[heir.Relation]int{heir.Daughter: 2, heir.SonsSon: 1},
		},
		{
			name:       "Full sisters with daughters (tasib)",
			sex:        heir.Male,
			heirCounts: map[heir.Relation]int{heir.Daughter: 2, heir.FullSister: 1},
		},
		{
			name:       "Consanguine sister with full sister blocked",
			sex:        heir.Male,
			heirCounts: map[heir.Relation]int{heir.ConsanguineSister: 1, heir.FullSister: 1},
		},
		{
			name:       "Mother with two siblings (sibling count)",
			sex:        heir.Male,
			heirCounts: map[heir.Relation]int{heir.Mother: 1, heir.FullBrother: 2},
		},
	}

	madhahib := solver.Madhahib()

	for _, c := range cases {
		w("### %s", c.name)
		w("")
		w("Deceased: **%s**", c.sex.String())
		w("")

		h := heir.New()
		for r, n := range c.heirCounts {
			h.Set(r, n)
		}
		ec := estate.Case{
			DeceasedSex: c.sex,
			Estate:      estate.Estate{Total: 120},
			Heirs:       h,
		}

		type schoolResult struct {
			name   string
			shares string
			flags  string
		}
		var results []schoolResult
		for _, m := range madhahib {
			res, err := solver.Solve(ec, m)
			if err != nil {
				results = append(results, schoolResult{name: m.Name, shares: "ERROR: " + err.Error()})
				continue
			}
			parts := make([]string, 0, len(res.Shares))
			for _, s := range res.Shares {
				parts = append(parts, fmt.Sprintf("%s: %s", s.Relation.String(), s.Fraction.String()))
			}
			flags := []string{}
			if res.Awl {
				flags = append(flags, "awl")
			}
			if res.Radd {
				flags = append(flags, "radd")
			}
			if res.NeedsReview {
				flags = append(flags, "needs-review")
			}
			results = append(results, schoolResult{
				name:   m.Name,
				shares: strings.Join(parts, "; "),
				flags:  strings.Join(flags, ", "),
			})
		}

		// Check whether all schools agree.
		allSame := true
		for i := 1; i < len(results); i++ {
			if results[i].shares != results[0].shares {
				allSame = false
				break
			}
		}

		if allSame {
			w("All four schools **agree**: %s", results[0].shares)
		} else {
			w("| School | Shares | Flags |")
			w("|--------|--------|-------|")
			for _, r := range results {
				w("| %s | %s | %s |", r.name, r.shares, r.flags)
			}
		}
		w("")
	}

	w("---")
	w("")
	w("## 6. Classical Test Coverage")
	w("")
	w("The table below summarises the `testdata/classical/` corpus used as the")
	w("regression suite. Test files are run by `go test ./internal/core/solver/...`.")
	w("")
	w("| File | Focus |")
	w("|------|-------|")
	w("| `awl.json` | Proportional reduction (awl) when fixed shares exceed the estate |")
	w("| `blocking.json` | Hajb hirman: total exclusion chains |")
	w("| `descendants.json` | Sons, daughters, sons' sons, sons' daughters |")
	w("| `jadd.json` | Grandfather with siblings (all four school views) |")
	w("| `mixed.json` | Multi-category configurations |")
	w("| `parents.json` | Father, mother, grandparents |")
	w("| `radd.json` | Return (radd) when fixed shares do not exhaust the estate |")
	w("| `residuary.json` | Asaba bi'l-nafs and asaba bi'l-ghair |")
	w("| `siblings.json` | Full, consanguine, and uterine siblings |")
	w("| `special.json` | Gharrawain, mushtaraka, and other named special cases |")
	w("| `spouses.json` | Husband and wife in various configurations |")

	w("")
	w("---")
	w("")
	w("## 7. Known Limitations")
	w("")

	limitations := []string{
		"**Dhawu al-arham (distant kindred)** are handled by a separate `DistributeDhawuArham` function and are not included in the main `Solve` result. The school flag `DhawuArham` is encoded and tested, but the distribution logic covers only the most common patterns.",
		"**Munasakha (cases with deceased heirs)** are not supported. If an heir dies before distribution is complete, the case must be decomposed into sequential sub-cases manually.",
		"**Bequests exceeding one third** are silently capped to one third unless `HeirsConsentToExcessBequest` is set. There is no UI control for this field in v1.0.",
		"**Non-monetary assets** (land, livestock, jewellery) are not modelled. All amounts are treated as fungible integers in the smallest currency unit.",
		"**Waqf and conditional estates** are outside scope. The engine assumes a freely distributable estate.",
		"**The LLM explain/parse features** are trial-tier only, gated behind `PUBLIC_LLM_ENABLED=true`, and are never the source of a legal result. Every LLM response carries an explicit disclaimer. The consistency guard rejects any explanation that contradicts the computed fractions.",
		"**`NeedsReview` flag**: the solver sets this flag on results it cannot fully resolve under the current rules (for example, some edge cases in the grandfather-with-siblings calculation). Such results should be audited manually.",
	}

	for i, lim := range limitations {
		w("%d. %s", i+1, lim)
		w("")
	}

	w("---")
	w("")
	w("*End of scholar review pack.*")

	// Count fixed-share rules and blocking rules for a summary.
	fixedCount := 0
	blockCount := 0
	for _, r := range heir.AllRelations() {
		fixedCount += len(rules.Rules(r))
		blockCount += len(rules.HirmanRules(r))
	}
	nuqsanCount := len(rules.NuqsanRules())

	// Summarise to stdout so CI can confirm generation.
	fmt.Printf("scholar-review.md written: %d fixed-share rules, %d blocking rules, %d nuqsan rules, %d divergence cases\n",
		fixedCount, blockCount, nuqsanCount, len(cases))

	_ = sort.Slice // keep sort imported; used implicitly via heir.AllRelations
}
