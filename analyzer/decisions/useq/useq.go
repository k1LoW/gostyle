package useq

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "useq"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#use-q"
	msg  = "Using %%q is recommended in output intended for humans where the input value could possibly be empty or contain control characters. (ref: https://google.github.io/styleguide/go/decisions#use-q )"

	badDQ  = "\"%s\""
	badDQ2 = "\\\"%s\\\""
	badSQ  = "'%s'"
)

var (
	disable          bool
	includeGenerated bool
)

// Analyzer based on https://google.github.io/styleguide/go/decisions#use-q
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://google.github.io/styleguide/go/decisions#use-q
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
		includeGenerated = c.AnalyzersSettings.Useq.IncludeGenerated
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
		switch e := n.(type) {
		case *ast.CallExpr:
			if len(e.Args) < 2 {
				return
			}
			if !isString(pass, e.Args[0]) {
				return
			}
			bl, ok := e.Args[0].(*ast.BasicLit)
			if !ok {
				return
			}
			if bl.Kind != token.STRING {
				return
			}

			format := bl.Value
			if (strings.Contains(format, badDQ) && format != badDQ) || strings.Contains(format, badDQ2) || strings.Contains(format, badSQ) {
				format = strings.ReplaceAll(format, "%", "%%")
				r.Append(e.Pos(), fmt.Sprintf("%s: %s", msg, format))
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

func isString(pass *analysis.Pass, e ast.Expr) bool {
	typ := pass.TypesInfo.TypeOf(e)
	if typ == nil {
		return false
	}
	if typ.String() == "string" {
		return true
	}
	return false
}
