package recvnames

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
	name = "recvnames"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#receiver-names"
	msg  = "Receiver variable names must be short (usually one or two letters in length)"
	msga = "Receiver variable names must be abbreviations for the type itself"
)

var (
	disable          bool
	includeGenerated bool
)

// Analyzer based on https://google.github.io/styleguide/go/decisions#receiver-names
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
		includeGenerated = c.AnalyzersSettings.Recvnames.IncludeGenerated
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
			var sn string
			for _, l := range n.Recv.List {
				switch t := l.Type.(type) {
				case *ast.StarExpr:
					sn = t.X.(*ast.Ident).Name
				case *ast.Ident:
					sn = t.Name
				}
				sn = strings.ToLower(sn)
				for _, n := range l.Names {
					if len(n.Name) > 2 {
						r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name))
					}
					for _, c := range n.Name {
						if !strings.ContainsRune(sn, c) {
							r.Append(n.Pos(), fmt.Sprintf("%s: %s", msga, n.Name))
							return
						}
					}
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
