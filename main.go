package main

import (
	"github.com/k1LoW/gostyle/analyzer/decisions/getters"
	"github.com/k1LoW/gostyle/analyzer/decisions/nilslices"
	"github.com/k1LoW/gostyle/analyzer/decisions/pkgnames"
	"github.com/k1LoW/gostyle/analyzer/decisions/recvnames"
	"github.com/k1LoW/gostyle/analyzer/decisions/repetition"
	"github.com/k1LoW/gostyle/analyzer/decisions/underscores"
	"github.com/k1LoW/gostyle/analyzer/decisions/varnames"
	"github.com/k1LoW/gostyle/analyzer/effective/ifacenames"
	"github.com/k1LoW/gostyle/analyzer/guide/mixedcaps"
	"github.com/k1LoW/gostyle/config"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(
		config.Loader,
		getters.AnalyzerWithConfig,
		ifacenames.AnalyzerWithConfig,
		pkgnames.AnalyzerWithConfig,
		mixedcaps.AnalyzerWithConfig,
		nilslices.AnalyzerWithConfig,
		recvnames.AnalyzerWithConfig,
		repetition.AnalyzerWithConfig,
		underscores.AnalyzerWithConfig,
		varnames.AnalyzerWithConfig,
	)
}
