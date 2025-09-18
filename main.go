package main

import (
	_ "embed"
	"fmt"
	"os"
	"slices"

	"github.com/k1LoW/gostyle/analyzer"
	"github.com/k1LoW/gostyle/cmd"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err) //nostyle:handlerrors
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) == 1 || (len(os.Args) > 1 && slices.Contains([]string{"run", "init", "completion", "-v", "help", "-h"}, os.Args[1])) {
		cmd.Execute()
		return nil
	}

	unitchecker.Main(analyzer.Analyzers...)
	return nil
}
