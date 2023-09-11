package pkgnames

import (
	"fmt"
	"go/ast"
	"slices"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "pkgnames"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#package-names"
	msg  = "Go package names should be short and contain only lowercase letters. A package name composed of multiple words should be left unbroken in all lowercase. (ref: https://google.github.io/styleguide/go/decisions#package-names)"
	msg2 = "Avoid uninformative package names like util, utility, common, helper, and so on. (ref: https://google.github.io/styleguide/go/decisions#package-names)"
)

var (
	disable        bool
	uninformatives = []string{
		"util",
		"utility",
		"common",
		"helper",
	}
)

// Analyzer based on https://google.github.io/styleguide/go/decisions#package-names
var Analyzer = &analysis.Analyzer{
	Name: "pkgnames",
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

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.Ident)(nil),
	}

	var pkg bool
	r, err := reporter.New(name, pass)
	if err != nil {
		return nil, err
	}
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.File:
			// Package names are not checked.
			pkg = true
			return
		case *ast.Ident:
			if pkg {
				if strings.Contains(strings.TrimSuffix(n.Name, "_test"), "_") {
					r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name))
				}
				if strings.ToLower(n.Name) != n.Name {
					r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name))
				}
				if slices.Contains(uninformatives, strings.ToLower(n.Name)) {
					r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg2, n.Name))
				}
			}
		}
		pkg = false
	})
	r.Report()
	return nil, nil
}

func init() {
	Analyzer.Flags.BoolVar(&disable, "disable", false, "disable "+name+" analyzer")
}
