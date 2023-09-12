package mixedcaps

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/detector"
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

var (
	disable          bool
	includeGenerated bool
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
	if disable {
		return nil, nil
	}
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.Ident)(nil),
	}

	opts := []reporter.Option{}
	if includeGenerated {
		opts = append(opts, reporter.IncludeGenerated())
	}
	r, err := reporter.New(name, pass, opts...)
	if err != nil {
		return nil, err
	}
	var pkg bool
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.File:
			pkg = true
			return
		case *ast.Ident:
			if pkg {
				// skip _test package
				if !detector.IsMixedCaps(strings.TrimSuffix(n.Name, "_test")) {
					r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name))
				}
				pkg = false
				return
			}
			if !detector.IsMixedCaps(n.Name) {
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name))
			}
		}
		pkg = false
	})
	r.Report()
	return nil, nil
}

func init() {
	Analyzer.Flags.BoolVar(&disable, "disable", false, "disable "+name+" analyzer")
	Analyzer.Flags.BoolVar(&includeGenerated, "include-generated", false, "include generated codes")
}
