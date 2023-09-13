package nilslices

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "nilslices"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#nil-slices"
	msg  = "if you declare an empty slice as a local variable (especially if it can be the source of a return value), prefer the nil initialization to reduce the risk of bugs by callers. (ref: https://google.github.io/styleguide/go/decisions#nil-slices)"
)

var (
	disable          bool
	includeGenerated bool
)

// Analyzer based on https://google.github.io/styleguide/go/guide#nil-slices
var Analyzer = &analysis.Analyzer{
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
		includeGenerated = c.AnalyzersSettings.Nilslices.IncludeGenerated
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
		(*ast.AssignStmt)(nil),
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
		switch n := n.(type) {
		case *ast.ValueSpec:
			for i, v := range n.Values {
				c, ok := v.(*ast.CompositeLit)
				if !ok {
					continue
				}
				if c.Elts != nil {
					continue
				}
				if _, ok := c.Type.(*ast.ArrayType); !ok {
					continue
				}
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Names[i].Name))
			}
		case *ast.AssignStmt:
			if n.Tok != token.DEFINE {
				return
			}
			for i, e := range n.Rhs {
				c, ok := e.(*ast.CompositeLit)
				if !ok {
					continue
				}
				if c.Elts != nil {
					continue
				}
				if _, ok := c.Type.(*ast.ArrayType); !ok {
					continue
				}
				id, ok := n.Lhs[i].(*ast.Ident)
				if !ok {
					continue
				}
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, id.Name))
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
