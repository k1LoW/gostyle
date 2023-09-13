package repetition

import (
	"fmt"
	"go/ast"
	"slices"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "repetition"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#repetition"
	msg  = "a piece of Go source code should avoid unnecessary repetition. (ref: https://google.github.io/styleguide/go/decisions#repetition)"
)

var (
	disable          bool
	includeGenerated bool
	exclude          string
)

// Analyzer based on https://google.github.io/styleguide/go/decisions#repetition
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
	words := strings.Split(exclude, ",")
	if c != nil {
		disable = c.IsDisabled(name)
		includeGenerated = c.AnalyzersSettings.Recvnames.IncludeGenerated
		words = c.AnalyzersSettings.Repetition.Exclude
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
		(*ast.FuncDecl)(nil),
	}

	opts := []reporter.Option{}
	if includeGenerated {
		opts = append(opts, reporter.IncludeGenerated())
	}
	r, err := reporter.New(name, pass, opts...)
	if err != nil {
		return nil, err
	}

	pkg := pass.Pkg.Name()
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ValueSpec:
			if n.Names == nil {
				return
			}
			if !n.Names[0].IsExported() {
				return
			}
			// Package vs. exported symbol name
			for _, name := range n.Names {
				splitted := camelcase.Split(name.Name)
				for _, s := range splitted {
					if len(s) == 1 {
						continue
					}
					if slices.Contains(words, s) {
						continue
					}
					if strings.Contains(pkg, strings.ToLower(s)) {
						r.Append(n.Pos(), fmt.Sprintf("%s: %s<-[%s]->%s", msg, pkg, s, name.Name))
					}
				}
			}
		case *ast.FuncDecl:
			if n.Recv != nil {
				return
			}
			if !n.Name.IsExported() {
				return
			}
			// Package vs. exported symbol name
			splitted := camelcase.Split(n.Name.Name)
			for _, s := range splitted {
				if len(s) == 1 {
					continue
				}
				if slices.Contains(words, s) {
					continue
				}
				if strings.Contains(pkg, strings.ToLower(s)) {
					r.Append(n.Pos(), fmt.Sprintf("%s: %s<-[%s]->%s", msg, pkg, s, n.Name.Name))
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
