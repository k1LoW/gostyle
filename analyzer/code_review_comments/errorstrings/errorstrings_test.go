package errorstrings

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	excludeTest = true
	td := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, td, Analyzer, "a")
}
