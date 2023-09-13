package varnames

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	smallVarnameMax = 3
	mediumVarnameMax = 5
	largeVarnameMax = 10
	veryLargeVarnameMax = 15

	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, Analyzer, "a")
}
