package mixedcaps

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "mixedcaps"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/guide#mixed-caps"
	msg  = "Go source code uses MixedCaps or mixedCaps (camel case) rather than underscores (snake case) when writing multi-word names. (ref: https://google.github.io/styleguide/go/guide#mixed-caps)"
)

// Analyzer based on https://google.github.io/styleguide/go/guide#mixed-caps
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.Ident)(nil),
	}

	var pkg bool
	r, err := reporter.New(name, pass)
	if err != nil {
		return nil, err
	}
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.File:
			// Package names are not checked.
			pkg = true
			return
		case *ast.Ident:
			if strings.Contains(n.Name, "_") && !strings.HasPrefix(n.Name, "_") && !pkg {
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name))
			}
		}
		pkg = false
	})
	r.Report()
	return nil, nil
}
