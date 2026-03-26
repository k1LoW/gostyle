package funcfmt

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	tests := []struct {
		calls bool
		pkg   string
	}{
		{false, "a"},
		{true, "b"},
	}
	for _, tt := range tests {
		checkCalls = tt.calls
		td := testutil.WithModules(t, analysistest.TestData(), nil)
		analysistest.Run(t, td, Analyzer, tt.pkg)
	}
}
