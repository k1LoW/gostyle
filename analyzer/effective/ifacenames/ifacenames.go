package ifacenames

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
	name = "ifacenames"
	doc  = "Analyzer based on https://go.dev/doc/effective_go#interface-names."
	msg  = "By convention, one-method interfaces are named by the method name plus an -er suffix or similar modification to construct an agent noun. (ref: https://go.dev/doc/effective_go#interface-names)"
	msgc = "All interface names with the -er suffix are required. (THIS IS NOT IN Effective Go)"
)

var (
	disable          bool
	includeGenerated bool
	all              bool
)

// Analyzer based on https://go.dev/doc/effective_go#interface-names.
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	URL:  "https://github.com/k1LoW/gostyle/tree/main/analyzer/effective/ifacenames",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://go.dev/doc/effective_go#interface-names.
var AnalyzerWithConfig = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	URL:  "https://github.com/k1LoW/gostyle/tree/main/analyzer/effective/ifacenames",
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
		includeGenerated = c.AnalyzersSettings.Ifacenames.IncludeGenerated
		all = c.AnalyzersSettings.Ifacenames.All
	}
	if disable {
		return nil, nil
	}
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
		(*ast.InterfaceType)(nil),
	}

	var ii *ast.Ident
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
		case *ast.InterfaceType:
			if len(n.Methods.List) == 1 && len(n.Methods.List[0].Names) > 0 {
				mn := n.Methods.List[0].Names[0].Name
				if !strings.HasPrefix(ii.Name, mn) || !strings.HasSuffix(ii.Name, "er") { // huristic
					r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, ii.Name))
					return
				}
			}
			if all && !strings.HasSuffix(ii.Name, "er") {
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msgc, ii.Name))
				return
			}
		case *ast.Ident:
			ii = n
		}
	})
	r.Report()
	return nil, nil
}

func init() {
	Analyzer.Flags.BoolVar(&disable, "disable", false, "disable "+name+" analyzer")
	Analyzer.Flags.BoolVar(&includeGenerated, "include-generated", false, "include generated codes")
	Analyzer.Flags.BoolVar(&all, "all", false, "all interface names with the -er suffix are required")
}
