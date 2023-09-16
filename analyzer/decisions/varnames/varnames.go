package varnames

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/k1LoW/gostyle/config"
	"github.com/k1LoW/gostyle/reporter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "varnames"
	doc  = "Analyzer based on https://google.github.io/styleguide/go/decisions#variable-names"
)

const (
	scopeSmall = iota
	scopeMedium
	scopeLarge
	scopeVeryLarge
)

var (
	disable             bool
	includeGenerated    bool
	smallScopeMax       int
	smallVarnameMax     int
	mediumScopeMax      int
	mediumVarnameMax    int
	largeScopeMax       int
	largeVarnameMax     int
	veryLargeVarnameMax int
)

// Analyzer based on https://google.github.io/styleguide/go/decisions#variable-names
var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://google.github.io/styleguide/go/decisions#variable-names
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
		includeGenerated = c.AnalyzersSettings.Varnames.IncludeGenerated
		smallScopeMax = c.AnalyzersSettings.Varnames.SmallScopeMax
		smallVarnameMax = c.AnalyzersSettings.Varnames.SmallVarnameMax
		mediumScopeMax = c.AnalyzersSettings.Varnames.MediumScopeMax
		mediumVarnameMax = c.AnalyzersSettings.Varnames.MediumVarnameMax
		largeScopeMax = c.AnalyzersSettings.Varnames.LargeScopeMax
		largeVarnameMax = c.AnalyzersSettings.Varnames.LargeVarnameMax
		veryLargeVarnameMax = c.AnalyzersSettings.Varnames.VeryLargeVarnameMax
	}
	if disable {
		return nil, nil
	}
	if smallVarnameMax <= 0 && mediumVarnameMax <= 0 && largeVarnameMax <= 0 && veryLargeVarnameMax <= 0 {
		return nil, nil
	}

	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from inspect: %T", pass.ResultOf[inspect.Analyzer])
	}

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
		(*ast.ValueSpec)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.RangeStmt)(nil),
	}

	var opts []reporter.Option
	if includeGenerated {
		opts = append(opts, reporter.IncludeGenerated())
	}
	r, err := reporter.New(name, pass, opts...)
	if err != nil {
		return nil, err
	}

	conf, err := config.NewTypesConfig(pass)
	if err != nil {
		return nil, err
	}
	pkg, err := conf.Check("varnames", pass.Fset, pass.Files, nil)
	if err != nil {
		return nil, err
	}
	sr := &scopeReporter{
		pkg:  pkg,
		r:    r,
		pass: pass,
	}
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ValueSpec:
			for _, id := range n.Names {
				sr.report(id.Pos(), id.Name)
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
				sr.report(id.Pos(), id.Name)
			}
		case *ast.RangeStmt:
			idk, ok := n.Key.(*ast.Ident)
			if ok {
				sr.report(idk.Pos(), idk.Name)
			}
			idv, ok := n.Value.(*ast.Ident)
			if ok {
				sr.report(idv.Pos(), idv.Name)
			}
		}
	})
	r.Report()
	return nil, nil
}

type scopeReporter struct {
	pkg  *types.Package
	r    *reporter.Reporter
	pass *analysis.Pass
}

func (sr *scopeReporter) report(pos token.Pos, varname string) {
	s := sr.pkg.Scope().Innermost(pos)
	o := s.Lookup(varname)
	if o == nil {
		s, o = s.LookupParent(varname, pos)
		if s == nil || o == nil {
			return
		}
	}
	switch o.(type) {
	case *types.Var, *types.Const:
		switch sr.scope(s) {
		case scopeSmall:
			if smallVarnameMax > 0 && len(varname) > smallVarnameMax {
				sr.r.Append(pos, fmt.Sprintf("'%s' is small scope. Variable name length of small scope should be less than or equal to %d. (THIS IS NOT IN Go Style)", varname, smallVarnameMax))
			}
		case scopeMedium:
			if mediumVarnameMax > 0 && len(varname) > mediumVarnameMax {
				sr.r.Append(pos, fmt.Sprintf("'%s' is medium scope. Variable name length of medium scope should be less than or equal to %d. (THIS IS NOT IN Go Style)", varname, mediumVarnameMax))
			}
		case scopeLarge:
			if largeVarnameMax > 0 && len(varname) > largeVarnameMax {
				sr.r.Append(pos, fmt.Sprintf("'%s' is large scope. Variable name length of large scope should be less than or equal to %d. (THIS IS NOT IN Go Style)", varname, largeVarnameMax))
			}
		case scopeVeryLarge:
			if veryLargeVarnameMax > 0 && len(varname) > veryLargeVarnameMax {
				sr.r.Append(pos, fmt.Sprintf("'%s' is very large scope. Variable name length of very large scope should be less than or equal to %d. (THIS IS NOT IN Go Style)", varname, veryLargeVarnameMax))
			}
		}
	}
}

func (sr *scopeReporter) scope(s *types.Scope) int {
	start := sr.pass.Fset.Position(s.Pos()).Line
	end := sr.pass.Fset.Position(s.End()).Line
	if start == 0 && end == 0 {
		return scopeVeryLarge
	}
	scope := end - start
	if scope <= smallScopeMax {
		return scopeSmall
	}
	if scope <= mediumScopeMax {
		return scopeMedium
	}
	if scope <= largeScopeMax {
		return scopeLarge
	}
	return scopeVeryLarge
}

func init() {
	Analyzer.Flags.BoolVar(&disable, "disable", false, "disable "+name+" analyzer")
	Analyzer.Flags.BoolVar(&includeGenerated, "include-generated", false, "include generated codes")
	Analyzer.Flags.IntVar(&smallScopeMax, "small-scope-max", config.DefaultSmallScopeMax, "max lines for small scope")
	Analyzer.Flags.IntVar(&smallVarnameMax, "small-varname-max", config.DefaultSmallVarnameMax, "max length of variable name for small scope")
	Analyzer.Flags.IntVar(&mediumScopeMax, "medium-scope-max", config.DefaultMediumScopeMax, "max lines for medium scope")
	Analyzer.Flags.IntVar(&mediumVarnameMax, "medium-varname-max", config.DefaultMediumVarnameMax, "max length of variable name for medium scope")
	Analyzer.Flags.IntVar(&largeScopeMax, "large-scope-max", config.DefaultLargeScopeMax, "max lines for large scope")
	Analyzer.Flags.IntVar(&largeVarnameMax, "large-varname-max", config.DefaultLargeVarnameMax, "max length of variable name for large scope")
	Analyzer.Flags.IntVar(&veryLargeVarnameMax, "very-large-varname-max", config.DefaultVeryLargeVarnameMax, "max length of variable name for very large scope")
}
