package recvtype

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "recvtype"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#receiver-type"
	msg  = "When in doubt, use a pointer receiver. (GOSTYLE MEMO: It's a strong check, so read Go Style and decide if it should be ignored or not proactively) (ref: https://google.github.io/styleguide/go/decisions#receiver-type )"
	msgm = "If the receiver is a map, function, or channel, use a value rather than a pointer. (ref: https://google.github.io/styleguide/go/decisions#receiver-type )"
)

var (
	disable          bool
	includeGenerated bool
)

// Analyzer based on https://google.github.io/styleguide/go/decisions#receiver-type
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://google.github.io/styleguide/go/decisions#receiver-type
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
		includeGenerated = c.AnalyzersSettings.Recvtype.IncludeGenerated
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
		case *ast.FuncDecl:
			if n.Recv == nil {
				return
			}
			for _, f := range n.Recv.List {
				switch e := f.Type.(type) {
				case *ast.StarExpr:
					typ := pass.TypesInfo.TypeOf(e.X)
					if typ == nil {
						return
					}
					if _, ok := typ.Underlying().(*types.Map); ok {
						r.Append(n.Pos(), fmt.Sprintf("%s: %s", msgm, f.Names[0].Name))
					}
					if _, ok := typ.Underlying().(*types.Signature); ok {
						r.Append(n.Pos(), fmt.Sprintf("%s: %s", msgm, f.Names[0].Name))
					}
					if _, ok := typ.Underlying().(*types.Chan); ok {
						r.Append(n.Pos(), fmt.Sprintf("%s: %s", msgm, f.Names[0].Name))
					}
				case *ast.Ident:
					typ := pass.TypesInfo.TypeOf(f.Type)
					if typ != nil {
						if _, ok := typ.Underlying().(*types.Map); ok {
							return
						}
						if _, ok := typ.Underlying().(*types.Signature); ok {
							return
						}
						if _, ok := typ.Underlying().(*types.Chan); ok {
							return
						}
					}
					r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name.Name))
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
