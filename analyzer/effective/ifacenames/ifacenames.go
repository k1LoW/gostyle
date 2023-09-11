package ifacenames

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "https://go.dev/doc/effective_go#interface-names"

// Analyzer for https://go.dev/doc/effective_go#interface-names
var Analyzer = &analysis.Analyzer{
	Name: "ifacenames",
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
		(*ast.InterfaceType)(nil),
	}

	var ii *ast.Ident
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.InterfaceType:
			if len(n.Methods.List) == 1 {
				mn := n.Methods.List[0].Names[0].Name
				if !strings.HasPrefix(ii.Name, mn) || !strings.HasSuffix(ii.Name, "er") {
					pass.Reportf(n.Pos(), "By convention, one-method interfaces are named by the method name plus an -er suffix or similar modification to construct an agent noun. (ref: https://go.dev/doc/effective_go#interface-names)")
				}
			}
		case *ast.Ident:
			ii = n
		}
	})

	return nil, nil
}
