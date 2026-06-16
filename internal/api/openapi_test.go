package api

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

const specPath = "../../openapi/faraid.yaml"

// loadSpec loads and validates the OpenAPI contract.
func loadSpec(t *testing.T) *openapi3.T {
	t.Helper()
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}
	if err := doc.Validate(loader.Context); err != nil {
		t.Fatalf("spec is not valid: %v", err)
	}
	return doc
}

func TestOpenAPILintsClean(t *testing.T) {
	doc := loadSpec(t)

	if doc.OpenAPI == "" || doc.Info == nil || doc.Info.Title != "Faraid API" {
		t.Errorf("unexpected spec header: %q %+v", doc.OpenAPI, doc.Info)
	}

	// Every endpoint the engine exposes must be declared.
	for _, p := range []string{"/solve", "/compare", "/cases", "/cases/{id}", "/madhahib", "/export", "/explain", "/healthz"} {
		if doc.Paths.Find(p) == nil {
			t.Errorf("missing path %q", p)
		}
	}

	// Key schemas must be present.
	for _, s := range []string{"CaseInput", "SolveRequest", "SolveResult", "HeirShare", "DerivationStep", "Comparison", "SavedCase", "Error"} {
		if _, ok := doc.Components.Schemas[s]; !ok {
			t.Errorf("missing schema %q", s)
		}
	}
}
