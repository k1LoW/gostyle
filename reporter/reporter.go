package reporter

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
)

const (
	defaultPrefixKey         = "gostyle"
	sep                      = ":"
	NoLintCommentAnnotation  = "nolint"
	NoStyleCommentAnnotation = "nostyle"
	LintIgnore               = "lint:ignore"
	IgnoreAll                = "all"
)

type Reporter struct {
	name              string
	pass              *analysis.Pass
	cm                comment.Maps
	reports           []*report
	prefix            string
	ignoreAnotation   string
	disableLintIgnore bool
	disableNoLint     bool
}

type report struct {
	pos token.Pos
	msg string
}

type ReporterOption func(*Reporter)

func IgnoreAnotation(s string) ReporterOption {
	return func(r *Reporter) {
		r.ignoreAnotation = s
	}
}

func DisableLintIgnore() ReporterOption {
	return func(r *Reporter) {
		r.disableLintIgnore = true
	}
}

func DisableNoLint() ReporterOption {
	return func(r *Reporter) {
		r.disableNoLint = true
	}
}

func Prefix(s string) ReporterOption {
	return func(r *Reporter) {
		r.prefix = s
	}
}

func New(name string, pass *analysis.Pass, opts ...ReporterOption) (*Reporter, error) {
	cm, ok := pass.ResultOf[commentmap.Analyzer].(comment.Maps)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from commentmap: %T", pass.ResultOf[commentmap.Analyzer])
	}
	r := &Reporter{
		name:            name,
		pass:            pass,
		cm:              cm,
		prefix:          fmt.Sprintf("[%s.%s] ", defaultPrefixKey, name),
		ignoreAnotation: NoStyleCommentAnnotation,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r, nil
}

func (r *Reporter) Append(pos token.Pos, msg string) {
	r.reports = append(r.reports, &report{pos: pos, msg: msg})
}

func (r *Reporter) Report() {
	for _, rr := range r.reports {
		if r.IgnoreReport(rr.pos) {
			continue
		}
		r.pass.Reportf(rr.pos, fmt.Sprintf("%s%s", r.prefix, rr.msg))
	}
}

func (r *Reporter) IgnoreReport(pos token.Pos) bool {
	f1 := r.pass.Fset.File(pos)
	for i := range r.cm {
		for n, cgs := range r.cm[i] {
			f2 := r.pass.Fset.File(n.Pos())
			if f1 != f2 {
				// different file
				continue
			}

			for _, cg := range cgs {
				if f1.Line(pos) != f2.Line(cg.Pos()) {
					// different line
					continue
				}

				for _, c := range cg.List {
					t := c.Text
					if !strings.HasPrefix(t, "//") {
						continue
					}
					if !r.disableLintIgnore {
						// '//lint:ignore'
						if strings.HasPrefix(t, fmt.Sprintf("//%s", LintIgnore)) {
							return true
						}
					}
					if !r.disableNoLint {
						// 'nolint:all'
						if strings.Contains(t, fmt.Sprintf("%s%s%s", NoLintCommentAnnotation, sep, IgnoreAll)) {
							return true
						}
					}
					// 'nostyle:all'
					if strings.Contains(t, fmt.Sprintf("%s%s%s", r.ignoreAnotation, sep, IgnoreAll)) {
						return true
					}
					// 'nostyle:' and r.name
					if strings.Contains(t, r.ignoreAnotation+sep) && strings.Contains(t, r.name) {
						return true
					}
				}
			}
		}
	}
	return false
}
