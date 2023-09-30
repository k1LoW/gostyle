package funcfmt

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
	name = "funcfmt"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#function-formatting"
	msgs = "The signature of a function or method declaration should remain on a single line to avoid indentation confusion. (ref: https://google.github.io/styleguide/go/decisions#function-formatting )"
	msgc = "Function and method calls should not be separated based solely on line length. (ref: https://google.github.io/styleguide/go/decisions#function-formatting )"
)

var (
	disable          bool
	includeGenerated bool
)

// Analyzer based on https://google.github.io/styleguide/go/decisions#function-formatting
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://google.github.io/styleguide/go/decisions#function-formatting
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
		includeGenerated = c.AnalyzersSettings.Funcfmt.IncludeGenerated
	}

	if disable {
		return nil, nil
	}
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
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
		switch nn := n.(type) {
		case *ast.FuncDecl:
			if len(nn.Type.Params.List) == 0 {
				return
			}
			for _, f := range nn.Type.Params.List {
				if len(f.Names) == 0 {
					continue
				}
				for _, id := range f.Names {
					if pass.Fset.Position(id.Pos()).Line != pass.Fset.Position(nn.Pos()).Line {
						r.Append(nn.Pos(), msgs)
						return
					}
				}
			}
		case *ast.CallExpr:
			if len(nn.Args) == 0 {
				return
			}
			for _, arg := range nn.Args {
				if pass.Fset.Position(arg.Pos()).Line != pass.Fset.Position(nn.Pos()).Line {
					r.Append(nn.Pos(), msgc)
					return
				}
			}
		}
	})
	r.Report()
	return nil, nil
}

func init() {
	Analyzer.Flags.BoolVar(&disable, "disable", false, "disable "+name+" analyzer")
	Analyzer.Flags.BoolVar(&includeGenerated, "include-generated", false, "include generated codes")
}
