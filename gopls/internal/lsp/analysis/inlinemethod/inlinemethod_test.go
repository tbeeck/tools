package inlinemethod_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/gopls/internal/lsp/analysis/inlinemethod"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, inlinemethod.Analyzer, "basic")
}

func TestParams(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, inlinemethod.Analyzer, "params")
}

func TestReturn(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, inlinemethod.Analyzer, "return")
}
