package mixedcaps

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "https://go.dev/doc/effective_go#mixed-caps"

// Analyzer for https://go.dev/doc/effective_go#mixed-caps
var Analyzer = &analysis.Analyzer{
	Name: "mixedcaps",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	i.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			if strings.Contains(n.Name, "_") {
				pass.Reportf(n.Pos(), "The convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names. (ref: https://go.dev/doc/effective_go#mixed-caps)")
			}
		}
	})

	return nil, nil
}
