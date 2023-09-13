package ifacenames

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	tests := []struct {
		all bool // -all flag
		pkg string
	}{
		{false, "a"},
		{true, "b"},
	}
	for _, tt := range tests {
		all = tt.all
		td := testutil.WithModules(t, analysistest.TestData(), nil)
		analysistest.Run(t, td, Analyzer, tt.pkg)
	}
}
