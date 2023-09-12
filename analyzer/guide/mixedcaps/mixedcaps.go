package mixedcaps

import (
	"fmt"
	"go/ast"
	"slices"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/detector"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "mixedcaps"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/guide#mixed-caps"
	msg  = "Go source code uses MixedCaps or mixedCaps (camel case) rather than underscores (snake case) when writing multi-word names. (ref: https://google.github.io/styleguide/go/guide#mixed-caps)"
)

var (
	disable          bool
	includeGenerated bool
	excludeWords     string
)

// Analyzer based on https://google.github.io/styleguide/go/guide#mixed-caps
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	if disable {
		return nil, nil
	}
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}
	words := strings.Split(excludeWords, ",")

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.ImportSpec)(nil),
		(*ast.ValueSpec)(nil),
		(*ast.TypeSpec)(nil),
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
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.File:
			pkgname := n.Name.Name
			if slices.Contains(words, pkgname) {
				return
			}
			if !detector.IsMixedCaps(strings.TrimSuffix(pkgname, "_test")) {
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, pkgname))
			}
		case *ast.ImportSpec:
			if n.Name == nil {
				return
			}
			if slices.Contains(words, n.Name.Name) {
				return
			}
			if !detector.IsMixedCaps(n.Name.Name) {
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name.Name))
			}
		case *ast.ValueSpec:
			for _, ident := range n.Names {
				if slices.Contains(words, ident.Name) {
					continue
				}
				if !detector.IsMixedCaps(ident.Name) {
					r.Append(ident.Pos(), fmt.Sprintf("%s: %s", msg, ident.Name))
				}
			}
		case *ast.TypeSpec:
			if slices.Contains(words, n.Name.Name) {
				return
			}
			if !detector.IsMixedCaps(n.Name.Name) {
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name.Name))
			}
		case *ast.FuncDecl:
			if !slices.Contains(words, n.Name.Name) {
				if !detector.IsMixedCaps(n.Name.Name) {
					r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name.Name))
				}
			}
			if n.Recv == nil {
				return
			}
			for _, field := range n.Recv.List {
				for _, ident := range field.Names {
					if slices.Contains(words, ident.Name) {
						continue
					}
					if !detector.IsMixedCaps(ident.Name) {
						r.Append(ident.Pos(), fmt.Sprintf("%s: %s", msg, ident.Name))
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
	Analyzer.Flags.StringVar(&excludeWords, "exclude-words", "", "exclude words (comma separated)")
}
