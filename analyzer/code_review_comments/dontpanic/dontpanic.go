package dontpanic

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "dontpanic"
	doc  = "Analyzer based on https://github.com/golang/go/wiki/CodeReviewComments#dont-panic"
	msg  = "Don't use panic for normal error handling. Use error and multiple return values. (ref: https://github.com/golang/go/wiki/CodeReviewComments#dont-panic )"
)

var (
	disable          bool
	includeGenerated bool
	excludeTest      bool
)

// Analyzer based on https://github.com/golang/go/wiki/CodeReviewComments#dont-panic
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://github.com/golang/go/wiki/CodeReviewComments#dont-panic
var AnalyzerWithConfig = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		config.Loader,
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	c, err := config.Load(pass)
	if err != nil {
		return nil, err
	}
	if c != nil {
		disable = c.IsDisabled(name)
		includeGenerated = c.AnalyzersSettings.Dontpanic.IncludeGenerated
		excludeTest = c.AnalyzersSettings.Dontpanic.ExcludeTest
	}
	if disable {
		return nil, nil
	}
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	var opts []reporter.Option
	if includeGenerated {
		opts = append(opts, reporter.IncludeGenerated())
	}
	r, err := reporter.New(name, pass, opts...)
	if err != nil {
		return nil, err
	}
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch e := n.(type) {
		case *ast.CallExpr:
			id, ok := e.Fun.(*ast.Ident)
			if !ok {
				return
			}
			if id.Name != "panic" {
				return
			}
			if excludeTest {
				if strings.HasSuffix(pass.Fset.File(e.Pos()).Name(), "_test.go") {
					return
				}
			}
			r.Append(e.Pos(), msg)
		}
	})
	r.Report()
	return nil, nil
}

func init() {
	Analyzer.Flags.BoolVar(&disable, "disable", false, "disable "+name+" analyzer")
	Analyzer.Flags.BoolVar(&includeGenerated, "include-generated", false, "include generated codes")
	Analyzer.Flags.BoolVar(&excludeTest, "exclude-test", false, "exclude test files")
}
