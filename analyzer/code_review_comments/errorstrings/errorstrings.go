package errorstrings

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "errorstrings"
	doc  = "Analyzer based on https://github.com/golang/go/wiki/CodeReviewComments#error-strings"
	msg  = "Error strings should not be capitalized (unless beginning with proper nouns or acronyms) or end with punctuation, since they are usually printed following other context. (ref: https://github.com/golang/go/wiki/CodeReviewComments#error-strings )"
)

var (
	disable          bool
	includeGenerated bool
	excludeTest      bool
)

// Analyzer based on https://github.com/golang/go/wiki/CodeReviewComments#error-strings
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://github.com/golang/go/wiki/CodeReviewComments#error-strings
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
	if c != nil {
		disable = c.IsDisabled(name)
		includeGenerated = c.AnalyzersSettings.Errorstrings.IncludeGenerated
		excludeTest = c.AnalyzersSettings.Errorstrings.ExcludeTest
	}
	if disable {
		return nil, nil
	}
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
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
		if excludeTest {
			if strings.HasSuffix(pass.Fset.File(n.Pos()).Name(), "_test.go") {
				return
			}
		}
		switch e := n.(type) {
		case *ast.CallExpr:
			if len(e.Args) == 0 {
				return
			}
			fn, ok := e.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}
			if fn.Sel.Name == "Errorf" {
				bl, ok := e.Args[0].(*ast.BasicLit)
				if !ok {
					return
				}
				if isNG(bl.Value) {
					r.Append(e.Pos(), fmt.Sprintf("%s: %s", msg, bl.Value))
				}
				return
			}
			id, ok := fn.X.(*ast.Ident)
			if !ok {
				return
			}
			if id.Name == "errors" && fn.Sel.Name == "New" {
				bl, ok := e.Args[0].(*ast.BasicLit)
				if !ok {
					return
				}
				if isNG(bl.Value) {
					r.Append(e.Pos(), fmt.Sprintf("%s: %s", msg, bl.Value))
				}
			}
		}
	})
	r.Report()
	return nil, nil
}

func isNG(in string) bool {
	f := strings.Trim(in, "\"'`")
	return unicode.IsUpper(rune(f[0])) || strings.HasSuffix(f, ".")
}

func init() {
	Analyzer.Flags.BoolVar(&disable, "disable", false, "disable "+name+" analyzer")
	Analyzer.Flags.BoolVar(&includeGenerated, "include-generated", false, "include generated codes")
	Analyzer.Flags.BoolVar(&excludeTest, "exclude-test", false, "exclude test files")
}
