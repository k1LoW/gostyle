package main

import (
	"github.com/k1LoW/gostyle/analyzer/effective/ifacenames"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(
		ifacenames.Analyzer,
	)
}
