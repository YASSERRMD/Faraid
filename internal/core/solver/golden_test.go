package solver

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/estate"
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// goldenCase is one classical worked example loaded from testdata/classical.
type goldenCase struct {
	Name        string            `json:"name"`
	Reference   string            `json:"reference"`
	DeceasedSex string            `json:"deceasedSex"`
	Estate      int64             `json:"estate"`
	Madhhab     string            `json:"madhhab"`
	Heirs       map[string]int    `json:"heirs"`
	Expected    map[string]string `json:"expected"`
	SpecialCase string            `json:"specialCase"`
	Base        int64             `json:"base"`
	Residue     string            `json:"residue"`
}

func parseSex(t *testing.T, s string) heir.Sex {
	t.Helper()
	switch s {
	case "male":
		return heir.Male
	case "female":
		return heir.Female
	default:
		t.Fatalf("unknown deceased sex %q", s)
		return heir.SexUnknown
	}
}

func buildHeirs(t *testing.T, m map[string]int) *heir.Heirs {
	t.Helper()
	h := heir.New()
	for name, count := range m {
		r, ok := heir.ParseRelation(name)
		if !ok {
			t.Fatalf("unknown relation %q", name)
		}
		h.Set(r, count)
	}
	return h
}

func mustFraction(t *testing.T, s string) rational.Fraction {
	t.Helper()
	f, err := rational.Parse(s)
	if err != nil {
		t.Fatalf("bad fraction %q: %v", s, err)
	}
	return f
}

func runGolden(t *testing.T, gc goldenCase) {
	m, ok := MadhhabByName(gc.Madhhab)
	if !ok {
		t.Fatalf("unknown madhhab %q", gc.Madhhab)
	}
	c := estate.Case{
		DeceasedSex: parseSex(t, gc.DeceasedSex),
		Estate:      estate.Estate{Total: gc.Estate},
		Heirs:       buildHeirs(t, gc.Heirs),
	}
	r, err := Solve(c, m)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}

	for name, want := range gc.Expected {
		rel, ok := heir.ParseRelation(name)
		if !ok {
			t.Fatalf("unknown relation %q", name)
		}
		if got := share(r, rel); !got.Equal(mustFraction(t, want)) {
			t.Errorf("%s: %s = %s, want %s", gc.Name, name, got, want)
		}
	}
	if gc.SpecialCase != "" && r.SpecialCase != gc.SpecialCase {
		t.Errorf("%s: special case = %q, want %q", gc.Name, r.SpecialCase, gc.SpecialCase)
	}
	if gc.Base != 0 && r.Base.Int64() != gc.Base {
		t.Errorf("%s: base = %s, want %d", gc.Name, r.Base, gc.Base)
	}
	if gc.Residue != "" && !r.Residue.Equal(mustFraction(t, gc.Residue)) {
		t.Errorf("%s: residue = %s, want %s", gc.Name, r.Residue, gc.Residue)
	}
}

func TestClassicalCorpus(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "..", "testdata", "classical", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no classical corpus files found")
	}

	total := 0
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		var cases []goldenCase
		if err := json.Unmarshal(data, &cases); err != nil {
			t.Fatalf("parse %s: %v", f, err)
		}
		for _, gc := range cases {
			t.Run(gc.Name, func(t *testing.T) { runGolden(t, gc) })
			total++
		}
	}
	t.Logf("ran %d classical golden cases from %d files", total, len(files))
}
