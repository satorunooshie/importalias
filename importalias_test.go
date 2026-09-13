package importalias

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "a")
}

func TestImportPathImpliedName(t *testing.T) {
	tests := map[string]string{
		"encoding/json":                 "json",
		"example.com/project/schema/v2": "schema",
		"example.com/project/go-thing":  "thing",
		"example.com/project/a.thing":   "a",
	}
	for importPath, want := range tests {
		if got := importPathImpliedName(importPath); got != want {
			t.Errorf("importPathImpliedName(%q) = %q, want %q", importPath, got, want)
		}
	}
}
