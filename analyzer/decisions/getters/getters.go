package getters

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/detector"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "getters"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#getters"
	msg  = "Function and method names should not use a Get or get prefix, unless the underlying concept uses the word “get” (e.g. an HTTP GET). Prefer starting the name with the noun directly, for example use Counts over GetCounts. (ref: https://google.github.io/styleguide/go/decisions#getters)"
)

var (
	disable          bool
	includeGenerated bool
	exclude          string
)

// Analyzer based on https://google.github.io/styleguide/go/decisions#getters
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://google.github.io/styleguide/go/decisions#getters
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
	words := strings.Split(exclude, ",")
	if c != nil {
		disable = c.IsDisabled(name)
		words = c.AnalyzersSettings.Getters.Exclude
		includeGenerated = c.AnalyzersSettings.Getters.IncludeGenerated
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
		(*ast.InterfaceType)(nil),
		(*ast.FuncDecl)(nil),
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
				_, ok := v.(*ast.FuncLit)
				if !ok {
					continue
				}
				id := n.Names[i]
				if slices.Contains(words, id.Name) {
					continue
				}
				if detector.HasGetPrefix(id.Name) {
					r.Append(id.Pos(), fmt.Sprintf("%s: %s", msg, id.Name))
				}
			}
		case *ast.InterfaceType:
			if n.Methods == nil {
				return
			}
			for _, field := range n.Methods.List {
				for _, id := range field.Names {
					if slices.Contains(words, id.Name) {
						continue
					}
					if detector.HasGetPrefix(id.Name) {
						r.Append(id.Pos(), fmt.Sprintf("%s: %s", msg, id.Name))
					}
				}
			}
		case *ast.FuncDecl:
			if slices.Contains(words, n.Name.Name) {
				return
			}
			if detector.HasGetPrefix(n.Name.Name) {
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name.Name))
			}
		case *ast.AssignStmt:
			if n.Tok != token.DEFINE {
				return
			}
			for i, e := range n.Rhs {
				_, ok := e.(*ast.FuncLit)
				if !ok {
					continue
				}
				id, ok := n.Lhs[i].(*ast.Ident)
				if !ok {
					continue
				}
				if slices.Contains(words, id.Name) {
					continue
				}
				if detector.HasGetPrefix(id.Name) {
					r.Append(id.Pos(), fmt.Sprintf("%s: %s", msg, id.Name))
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
	Analyzer.Flags.StringVar(&exclude, "exclude", "", "exclude words (comma separated)")
}
