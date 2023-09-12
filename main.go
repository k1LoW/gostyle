package main

import (
	"github.com/k1LoW/gostyle/analyzer/decisions/pkgnames"
	"github.com/k1LoW/gostyle/analyzer/decisions/recvnames"
	"github.com/k1LoW/gostyle/analyzer/effective/ifacenames"
	"github.com/k1LoW/gostyle/analyzer/guide/mixedcaps"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(
		ifacenames.Analyzer,
		pkgnames.Analyzer,
		mixedcaps.Analyzer,
		recvnames.Analyzer,
	)
}
