package main

import (
	_ "embed"
	"os"
	"slices"

	"github.com/k1LoW/gostyle/analyzer"
	"github.com/k1LoW/gostyle/cmd"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	if len(os.Args) == 1 || (len(os.Args) > 1 && slices.Contains([]string{"run", "init", "completion", "-v", "help", "-h"}, os.Args[1])) {
		cmd.Execute()
		return
	}

	unitchecker.Main(analyzer.Analyzers...)
}
