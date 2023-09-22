package repetition

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
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
	msgp = "When naming exported symbols, the name of the package is always visible outside your package, so redundant information between the two should be reduced or eliminated. (ref: https://google.github.io/styleguide/go/decisions#package-vs-exported-symbol-name)"
	msgt = "The compiler always knows the type of a variable, and in most cases it is also clear to the reader what type a variable is by how it is used. It is only necessary to clarify the type of a variable if its value appears twice in the same scope. (ref: https://google.github.io/styleguide/go/decisions#variable-name-vs-type)"
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
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

// AnalyzerWithConfig based on https://google.github.io/styleguide/go/decisions#repetition
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
		(*ast.AssignStmt)(nil),
		(*ast.RangeStmt)(nil),
		(*ast.FuncDecl)(nil),
	}

	var opts []reporter.Option
	if includeGenerated {
		opts = append(opts, reporter.IncludeGenerated())
	}
	r, err := reporter.New(name, pass, opts...)
	if err != nil {
		return nil, err
	}

	tr := &typeVarReporter{
		r:       r,
		pass:    pass,
		exclude: words,
	}
	pkgn := pass.Pkg.Name()
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ValueSpec:
			if len(n.Names) == 0 {
				return
			}
			// Package vs. exported symbol name
			for _, id := range n.Names {
				tr.report(id.Pos(), id.Name)
				if !id.IsExported() {
					continue
				}
				splitted := camelcase.Split(id.Name)
				if slices.Contains(words, id.Name) {
					continue
				}
				for _, s := range splitted {
					if len(s) == 1 {
						continue
					}
					if strings.HasSuffix(pkgn, "_test") && strings.ToLower(s) == "test" {
						continue
					}
					if strings.Contains(pkgn, strings.ToLower(s)) {
						r.Append(n.Pos(), fmt.Sprintf("%s: %s<-[%s]->%s", msgp, pkgn, s, id.Name))
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
				tr.report(id.Pos(), id.Name)
			}
		case *ast.RangeStmt:
			idk, ok := n.Key.(*ast.Ident)
			if ok {
				tr.report(idk.Pos(), idk.Name)
			}
			idv, ok := n.Value.(*ast.Ident)
			if ok {
				tr.report(idv.Pos(), idv.Name)
			}
		case *ast.FuncDecl:
			if n.Recv != nil {
				return
			}
			if !n.Name.IsExported() {
				return
			}
			if slices.Contains(words, n.Name.Name) {
				return
			}
			if strings.HasPrefix(n.Name.Name, "Test") {
				if n.Type != nil && n.Type.Params != nil && len(n.Type.Params.List) > 0 {
					switch t := n.Type.Params.List[0].Type.(type) {
					case *ast.StarExpr:
						s, ok := t.X.(*ast.SelectorExpr)
						if ok {
							id, ok := s.X.(*ast.Ident)
							if ok && id.Name == "testing" {
								return
							}
						}
					case *ast.SelectorExpr:
						id, ok := t.X.(*ast.Ident)
						if ok && id.Name == "testing" {
							return
						}
					}
				}
			}

			// Package vs. exported symbol name
			splitted := camelcase.Split(n.Name.Name)
			for _, s := range splitted {
				if len(s) == 1 {
					continue
				}
				if strings.HasSuffix(pkgn, "_test") && strings.ToLower(s) == "test" {
					continue
				}
				if strings.Contains(pkgn, strings.ToLower(s)) {
					r.Append(n.Pos(), fmt.Sprintf("%s: %s<-[%s]->%s", msgp, pkgn, s, n.Name.Name))
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

type typeVarReporter struct {
	r       *reporter.Reporter
	pass    *analysis.Pass
	exclude []string
}

func (tr *typeVarReporter) report(pos token.Pos, varname string) {
	if slices.Contains(tr.exclude, varname) {
		return
	}
	// Variable name vs. type.
	s := tr.pass.Pkg.Scope().Innermost(pos)
	o := s.Lookup(varname)
	if o == nil {
		s, o = s.LookupParent(varname, pos)
		if s == nil || o == nil {
			return
		}
	}
	switch o.(type) {
	case *types.Var, *types.Const:
		switch o.Type() {
		case types.Typ[types.Int], types.Typ[types.Int8], types.Typ[types.Int16], types.Typ[types.Int32], types.Typ[types.Int64]:
			if strings.Contains(strings.ToLower(varname), "int") {
				tr.r.Append(pos, fmt.Sprintf("%s: %s<-[%s]->%s", msgt, varname, "int", o.Type().String()))
			} else if strings.Contains(strings.ToLower(varname), "num") {
				tr.r.Append(pos, fmt.Sprintf("%s: %s<-[%s]->%s", msgt, varname, "num", o.Type().String()))
			}
		case types.Typ[types.Uint], types.Typ[types.Uint8], types.Typ[types.Uint16], types.Typ[types.Uint32], types.Typ[types.Uint64]:
			if strings.Contains(strings.ToLower(varname), "uint") {
				tr.r.Append(pos, fmt.Sprintf("%s: %s<-[%s]->%s", msgt, varname, "uint", o.Type().String()))
			} else if strings.Contains(strings.ToLower(varname), "num") {
				tr.r.Append(pos, fmt.Sprintf("%s: %s<-[%s]->%s", msgt, varname, "num", o.Type().String()))
			}
		case types.Typ[types.Float32], types.Typ[types.Float64]:
			if strings.Contains(strings.ToLower(varname), "float") {
				tr.r.Append(pos, fmt.Sprintf("%s: %s<-[%s]->%s", msgt, varname, "float", o.Type().String()))
			} else if strings.Contains(strings.ToLower(varname), "num") {
				tr.r.Append(pos, fmt.Sprintf("%s: %s<-[%s]->%s", msgt, varname, "num", o.Type().String()))
			}
		case types.Typ[types.String]:
			if strings.Contains(strings.ToLower(varname), "string") {
				tr.r.Append(pos, fmt.Sprintf("%s: %s<-[%s]->%s", msgt, varname, "string", o.Type().String()))
			} else if strings.Contains(strings.ToLower(varname), "str") {
				tr.r.Append(pos, fmt.Sprintf("%s: %s<-[%s]->%s", msgt, varname, "str", o.Type().String()))
			}
		case types.Typ[types.Bool]:
			if strings.Contains(strings.ToLower(varname), "bool") {
				tr.r.Append(pos, fmt.Sprintf("%s: %s<-[%s]->%s", msgt, varname, "bool", o.Type().String()))
			}
		default:
			if strings.Contains(strings.ToLower(varname), strings.ToLower(o.Type().String())) {
				tr.r.Append(pos, fmt.Sprintf("%s: %s<-[%s]->%s", msgt, tr.pass.Pkg.Name(), o.Type().String(), varname))
			}
		}
	}
}
