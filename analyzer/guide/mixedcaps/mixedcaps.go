package mixedcaps

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
	name = "mixedcaps"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/guide#mixed-caps"
	msg  = "Go source code uses MixedCaps or mixedCaps (camel case) rather than underscores (snake case) when writing multi-word names. (ref: https://google.github.io/styleguide/go/guide#mixed-caps )"
)

var (
	disable          bool
	includeGenerated bool
	exclude          string
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

// AnalyzerWithConfig based on https://google.github.io/styleguide/go/guide#mixed-caps
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
	var opts []reporter.Option
	if c != nil {
		disable = c.IsDisabled(name)
		words = c.AnalyzersSettings.Mixedcaps.Exclude
		includeGenerated = c.AnalyzersSettings.Mixedcaps.IncludeGenerated
		opts = append(opts, reporter.ExcludeFiles(c.ConfigDir, c.ExcludeFiles))
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
		(*ast.ValueSpec)(nil),
		(*ast.TypeSpec)(nil),
		(*ast.InterfaceType)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.RangeStmt)(nil),
	}

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
			pkg := n.Name.Name
			if slices.Contains(words, pkg) {
				return
			}
			if !detector.IsMixedCaps(strings.TrimSuffix(pkg, "_test")) {
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, pkg))
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
			for _, id := range n.Names {
				if slices.Contains(words, id.Name) {
					continue
				}
				if !detector.IsMixedCaps(id.Name) {
					r.Append(id.Pos(), fmt.Sprintf("%s: %s", msg, id.Name))
				}
			}
		case *ast.TypeSpec:
			if slices.Contains(words, n.Name.Name) {
				return
			}
			if !detector.IsMixedCaps(n.Name.Name) {
				r.Append(n.Pos(), fmt.Sprintf("%s: %s", msg, n.Name.Name))
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
					if !detector.IsMixedCaps(id.Name) {
						r.Append(id.Pos(), fmt.Sprintf("%s: %s", msg, id.Name))
					}
				}
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
				for _, id := range field.Names {
					if slices.Contains(words, id.Name) {
						continue
					}
					if !detector.IsMixedCaps(id.Name) {
						r.Append(id.Pos(), fmt.Sprintf("%s: %s", msg, id.Name))
					}
				}
			}
		case *ast.AssignStmt:
			if n.Tok != token.DEFINE {
				return
			}
			for _, e := range n.Lhs {
				id, ok := e.(*ast.Ident)
				if !ok {
					continue
				}
				if slices.Contains(words, id.Name) {
					continue
				}
				if !detector.IsMixedCaps(id.Name) {
					r.Append(id.Pos(), fmt.Sprintf("%s: %s", msg, id.Name))
				}
			}
		case *ast.RangeStmt:
			idk, ok := n.Key.(*ast.Ident)
			if ok && !slices.Contains(words, idk.Name) && !detector.IsMixedCaps(idk.Name) {
				r.Append(idk.Pos(), fmt.Sprintf("%s: %s", msg, idk.Name))
			}
			idv, ok := n.Value.(*ast.Ident)
			if ok && !slices.Contains(words, idv.Name) && !detector.IsMixedCaps(idv.Name) {
				r.Append(idv.Pos(), fmt.Sprintf("%s: %s", msg, idv.Name))
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
