package useany

import (
	"fmt"
	"go/ast"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "useany"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#use-any"
	msg  = "Because it is an alias, `any` is equivalent to `interface{}` in many situations and in others it is easily interchangeable via an explicit conversion. Prefer to use `any` in new code. (ref: https://google.github.io/styleguide/go/decisions#use-any)"
	iff  = "interface{}"
)

var (
	disable          bool
	includeGenerated bool
)

// Analyzer based on https://google.github.io/styleguide/go/decisions#use-any
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://google.github.io/styleguide/go/decisions#use-any
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
		includeGenerated = c.AnalyzersSettings.Useany.IncludeGenerated
	}
	if disable {
		return nil, nil
	}
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
		(*ast.ValueSpec)(nil),
		(*ast.Field)(nil),
		(*ast.CompositeLit)(nil),
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
		switch nn := n.(type) {
		case *ast.ValueSpec:
			if !hasInterfaceType(nn.Type) {
				return
			}
			r.Append(nn.Pos(), fmt.Sprintf("%s: %s", msg, iff))
		case *ast.Field:
			ft, ok := nn.Type.(*ast.InterfaceType)
			if !ok {
				return
			}
			if len(ft.Methods.List) > 0 {
				return
			}
			r.Append(nn.Pos(), fmt.Sprintf("%s: %s", msg, iff))
		case *ast.CompositeLit:
			if !hasInterfaceType(nn.Type) {
				return
			}
			r.Append(nn.Pos(), fmt.Sprintf("%s: %s", msg, iff))
		}
	})
	r.Report()
	return nil, nil
}

func init() {
	Analyzer.Flags.BoolVar(&disable, "disable", false, "disable "+name+" analyzer")
	Analyzer.Flags.BoolVar(&includeGenerated, "include-generated", false, "include generated codes")
}

func hasInterfaceType(e ast.Expr) bool {
	switch tt := e.(type) {
	case *ast.InterfaceType:
		if len(tt.Methods.List) > 0 {
			return false
		}
		return true
	case *ast.ArrayType:
		return hasInterfaceType(tt.Elt)
	case *ast.MapType:
		return hasInterfaceType(tt.Value)
	}
	return false
}
