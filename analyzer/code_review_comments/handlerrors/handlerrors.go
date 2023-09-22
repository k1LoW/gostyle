package handlerrors

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "handlerrors"
	doc  = "Analyzer based on https://github.com/golang/go/wiki/CodeReviewComments#handle-errors"
	msg  = "Do not discard errors using `_` variables. If a function returns an error, check it to make sure the function succeeded. Handle the error, return it, or, in truly exceptional situations, panic. (ref: https://github.com/golang/go/wiki/CodeReviewComments#handle-errors)"
)

var (
	disable          bool
	includeGenerated bool
	excludeTest      bool
)

// Analyzer based on https://github.com/golang/go/wiki/CodeReviewComments#handle-errors
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://github.com/golang/go/wiki/CodeReviewComments#handle-errors
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

var errTyp = func() *types.Interface {
	return types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}()

func run(pass *analysis.Pass) (any, error) {
	c, err := config.Load(pass)
	if err != nil {
		return nil, err
	}
	if c != nil {
		disable = c.IsDisabled(name)
		includeGenerated = c.AnalyzersSettings.Handlerrors.IncludeGenerated
		excludeTest = c.AnalyzersSettings.Handlerrors.ExcludeTest
	}
	if disable {
		return nil, nil
	}
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
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

	br := &blankErrReporter{
		r:    r,
		pass: pass,
	}

	i.Preorder(nodeFilter, func(n ast.Node) {
		switch nn := n.(type) {
		case *ast.AssignStmt:
			if len(nn.Rhs) == 0 {
				return
			}
			e, ok := nn.Rhs[0].(*ast.CallExpr)
			if !ok {
				return
			}
			for i, l := range nn.Lhs {
				id, ok := l.(*ast.Ident)
				if !ok {
					continue
				}
				if id.Name != "_" {
					continue
				}
				if excludeTest {
					if strings.HasSuffix(pass.Fset.File(e.Pos()).Name(), "_test.go") {
						return
					}
				}
				br.report(e, i)
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

type blankErrReporter struct {
	r    *reporter.Reporter
	pass *analysis.Pass
}

func (br *blankErrReporter) report(e *ast.CallExpr, i int) {
	typ, ok := br.pass.TypesInfo.Types[e]
	if !ok {
		return
	}
	switch t := typ.Type.(type) {
	case *types.Named:
		if types.Implements(t, errTyp) {
			br.r.Append(e.Pos(), msg)
		}
	case *types.Pointer:
		if types.Implements(t, errTyp) {
			br.r.Append(e.Pos(), msg)
		}
	case *types.Tuple:
		if t.Len() <= i {
			return
		}
		if types.Implements(t.At(i).Type(), errTyp) {
			br.r.Append(e.Pos(), msg)
		}
	}
}
