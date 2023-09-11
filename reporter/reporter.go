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
	NoLintCommentAnnotation  = "nolint:"
	NoStyleCommentAnnotation = "nostyle:"
	LintIgnore               = "lint:ignore"
	IgnoreAll                = "all"
)

type Reporter struct {
	name    string
	pass    *analysis.Pass
	cm      comment.Maps
	reports []*Report
}

type Report struct {
	Pos token.Pos
	Msg string
}

func New(name string, pass *analysis.Pass) (*Reporter, error) {
	cm, ok := pass.ResultOf[commentmap.Analyzer].(comment.Maps)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from commentmap: %T", pass.ResultOf[commentmap.Analyzer])
	}
	return &Reporter{name: name, pass: pass, cm: cm}, nil
}

func (r *Reporter) Append(pos token.Pos, msg string) {
	r.reports = append(r.reports, &Report{Pos: pos, Msg: msg})
}

func (r *Reporter) Report() {
	for _, report := range r.reports {
		if r.IgnoreReport(report.Pos) {
			continue
		}
		r.pass.Reportf(report.Pos, report.Msg)
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
					continue
				}

				for _, c := range cg.List {
					t := c.Text
					if !strings.HasPrefix(t, "//") {
						continue
					}
					// '//lint:ignore'
					if strings.HasPrefix(t, fmt.Sprintf("//%s", LintIgnore)) {
						return true
					}
					// '//nolint:' or '//nostyle:'
					if !strings.HasPrefix(t, fmt.Sprintf("//%s", NoLintCommentAnnotation)) && !strings.HasPrefix(t, fmt.Sprintf("//%s", NoStyleCommentAnnotation)) {
						continue
					}
					// 'nolint:all'
					if strings.Contains(t, fmt.Sprintf("%s%s", NoLintCommentAnnotation, IgnoreAll)) {
						return true
					}
					// 'nostyle:all'
					if strings.Contains(t, fmt.Sprintf("%s%s", NoStyleCommentAnnotation, IgnoreAll)) {
						return true
					}
					// 'nostyle:' and r.name
					if strings.Contains(t, NoStyleCommentAnnotation) && strings.Contains(t, r.name) {
						return true
					}
				}
			}
		}
	}
	return false
}
