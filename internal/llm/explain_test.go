package llm

import (
	"context"
	"errors"
	"testing"
)

const groundTruth = "[result] husband = 1/4\n[result] son = 1/2\n[result] daughter = 1/4"

func TestExplainConsistent(t *testing.T) {
	prose := "The husband receives 1/4 of the estate, the son 1/2, and the daughter 1/4."
	exp, err := Explain(context.Background(), fakeCompleter{text: prose},
		groundTruth, []string{"1/4", "1/2", "1/4"})
	if err != nil {
		t.Fatal(err)
	}
	if !exp.Consistent || exp.Text != prose || !exp.Experimental {
		t.Errorf("expected consistent prose, got %+v", exp)
	}
}

func TestExplainEquivalentFractionStillConsistent(t *testing.T) {
	// 2/8 equals 1/4 by value, so it passes the guard.
	prose := "The husband receives 2/8 of the estate."
	exp, _ := Explain(context.Background(), fakeCompleter{text: prose}, groundTruth, []string{"1/4"})
	if !exp.Consistent {
		t.Error("2/8 should match 1/4 by value")
	}
}

func TestExplainDriftRejected(t *testing.T) {
	// The prose invents 1/3, which the engine never produced.
	prose := "The husband receives 1/3 of the estate."
	exp, err := Explain(context.Background(), fakeCompleter{text: prose}, groundTruth, []string{"1/4", "1/2"})
	if err != nil {
		t.Fatal(err)
	}
	if exp.Consistent {
		t.Error("drifting prose should be rejected")
	}
	if exp.Text != groundTruth {
		t.Error("on rejection, the raw derivation must be returned")
	}
}

func TestExplainCompleterError(t *testing.T) {
	if _, err := Explain(context.Background(), fakeCompleter{err: errors.New("boom")}, groundTruth, nil); err == nil {
		t.Error("completer error should propagate")
	}
}
