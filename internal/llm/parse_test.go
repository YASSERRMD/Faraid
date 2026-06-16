package llm

import (
	"context"
	"errors"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
)

// fakeCompleter returns a canned response, for testing the parser without a
// network.
type fakeCompleter struct {
	text string
	err  error
}

func (f fakeCompleter) Complete(context.Context, Request) (Response, error) {
	return Response{Text: f.text}, f.err
}

func TestParseCaseValid(t *testing.T) {
	resp := `{"deceasedSex":"female","estateTotal":2400,"heirs":{"husband":1,"son":1,"daughter":1},"notes":"ok"}`
	p, err := ParseCase(context.Background(), fakeCompleter{text: resp}, "a woman left a husband, a son, and a daughter")
	if err != nil {
		t.Fatal(err)
	}
	if p.DeceasedSex != "female" || p.Heirs["son"] != 1 || p.EstateTotal != 2400 {
		t.Errorf("unexpected proposal: %+v", p)
	}
	// The proposal converts to a valid case.
	c, err := p.ToCase()
	if err != nil || c.Heirs.Count(heir.Son) != 1 {
		t.Errorf("ToCase failed: %v", err)
	}
}

func TestParseCaseStripsFences(t *testing.T) {
	resp := "Here is the case:\n```json\n{\"deceasedSex\":\"male\",\"heirs\":{\"son\":2}}\n```\n"
	p, err := ParseCase(context.Background(), fakeCompleter{text: resp}, "two sons")
	if err != nil {
		t.Fatalf("should tolerate fences and prose: %v", err)
	}
	if p.Heirs["son"] != 2 {
		t.Errorf("proposal = %+v", p)
	}
}

func TestParseCaseFailures(t *testing.T) {
	cases := map[string]string{
		"not json":     `I cannot help with that`,
		"unknown heir": `{"deceasedSex":"male","heirs":{"cousin":1}}`,
		"bad sex":      `{"deceasedSex":"other","heirs":{"son":1}}`,
		"impossible":   `{"deceasedSex":"male","heirs":{"husband":1}}`,
	}
	for name, resp := range cases {
		if _, err := ParseCase(context.Background(), fakeCompleter{text: resp}, "x"); err == nil {
			t.Errorf("%s: expected an error so the caller falls back to manual entry", name)
		}
	}
}

func TestParseCaseCompleterError(t *testing.T) {
	if _, err := ParseCase(context.Background(), fakeCompleter{err: errors.New("boom")}, "x"); err == nil {
		t.Error("completer error should propagate")
	}
}
