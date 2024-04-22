package contexts

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
	name = "contexts"
	doc  = "Analyzer based on https://go.dev/wiki/CodeReviewComments#contexts"
	msgp = "Most functions that use a Context should accept it as their first parameter. (ref: https://go.dev/wiki/CodeReviewComments#contexts )"
	msgs = "Don't add a Context member to a struct type; instead add a ctx parameter to each method on that type that needs to pass it along. The one exception is for methods whose signature must match an interface in the standard library or in a third party library. (ref: https://go.dev/wiki/CodeReviewComments#contexts )"
)

var (
	disable          bool
	includeGenerated bool
	excludeTest      bool
)

// Analyzer based on https://go.dev/wiki/CodeReviewComments#contexts
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://go.dev/wiki/CodeReviewComments#contexts
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
	var opts []reporter.Option
	if c != nil {
		disable = c.IsDisabled(name)
		includeGenerated = c.AnalyzersSettings.Contexts.IncludeGenerated
		excludeTest = c.AnalyzersSettings.Contexts.ExcludeTest
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
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
		(*ast.StructType)(nil),
	}

	if includeGenerated {
		opts = append(opts, reporter.IncludeGenerated())
	}
	r, err := reporter.New(name, pass, opts...)
	if err != nil {
		return nil, err
	}

	i.Preorder(nodeFilter, func(n ast.Node) {
		switch nn := n.(type) {
		case *ast.FuncDecl:
			if excludeTest {
				if strings.HasSuffix(pass.Fset.File(nn.Pos()).Name(), "_test.go") {
					return
				}
			}
			if len(nn.Type.Params.List) < 2 {
				return
			}
			for _, p := range nn.Type.Params.List[1:] {
				e, ok := p.Type.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				id, ok := e.X.(*ast.Ident)
				if !ok {
					continue
				}
				if id.Name == "context" && e.Sel.Name == "Context" {
					r.Append(nn.Pos(), fmt.Sprintf("%s: %s", msgp, nn.Name.Name))
				}
			}
		case *ast.FuncLit:
			if excludeTest {
				if strings.HasSuffix(pass.Fset.File(nn.Pos()).Name(), "_test.go") {
					return
				}
			}
			if len(nn.Type.Params.List) < 2 {
				return
			}
			for _, p := range nn.Type.Params.List[1:] {
				e, ok := p.Type.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				id, ok := e.X.(*ast.Ident)
				if !ok {
					continue
				}
				if id.Name == "context" && e.Sel.Name == "Context" {
					r.Append(nn.Pos(), msgp)
				}
			}
		case *ast.StructType:
			if excludeTest {
				if strings.HasSuffix(pass.Fset.File(nn.Pos()).Name(), "_test.go") {
					return
				}
			}
			if len(nn.Fields.List) == 0 {
				return
			}
			for _, f := range nn.Fields.List {
				e, ok := f.Type.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				id, ok := e.X.(*ast.Ident)
				if !ok {
					continue
				}
				if id.Name == "context" && e.Sel.Name == "Context" {
					r.Append(e.Pos(), fmt.Sprintf("%s: %s", msgs, f.Names[0].Name))
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
	Analyzer.Flags.BoolVar(&excludeTest, "exclude-test", false, "exclude test files")
}
