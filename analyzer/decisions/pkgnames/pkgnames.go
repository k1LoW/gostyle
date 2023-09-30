package pkgnames

import (
	"fmt"
	"go/ast"
	"slices"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "pkgnames"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#package-names"
	msg  = "Go package names should be short and contain only lowercase letters. A package name composed of multiple words should be left unbroken in all lowercase. (ref: https://google.github.io/styleguide/go/decisions#package-names )"
	msg2 = "Avoid uninformative package names like util, utility, common, helper, and so on. (ref: https://google.github.io/styleguide/go/decisions#package-names )"
)

var (
	disable          bool
	includeGenerated bool
	uninformatives   = []string{
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

// AnalyzerWithConfig based on https://google.github.io/styleguide/go/decisions#package-names
var AnalyzerWithConfig = &analysis.Analyzer{
	Name: "pkgnames",
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
		includeGenerated = c.AnalyzersSettings.Pkgnames.IncludeGenerated
	}
	if disable {
		return nil, nil
	}
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.ImportSpec)(nil),
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
		var pkgname string
		switch n := n.(type) {
		case *ast.ImportSpec:
			if n.Name != nil {
				pkgname = n.Name.Name
			}
		case *ast.File:
			pkgname = n.Name.Name
		}
		if pkgname == "" || pkgname == "_" {
			return
		}
		if strings.Contains(strings.TrimSuffix(pkgname, "_test"), "_") {
			r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, pkgname))
		}
		if strings.ToLower(pkgname) != pkgname {
			r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, pkgname))
		}
		if slices.Contains(uninformatives, strings.ToLower(pkgname)) {
			r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg2, pkgname))
		}
	})
	r.Report()
	return nil, nil
}

func init() {
	Analyzer.Flags.BoolVar(&disable, "disable", false, "disable "+name+" analyzer")
	Analyzer.Flags.BoolVar(&includeGenerated, "include-generated", false, "include generated codes")
}
