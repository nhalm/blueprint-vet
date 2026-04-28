// Package mockgendirective implements an analyzer that flags
// *_interface.go files missing a //go:generate mockgen directive.
package mockgendirective

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "mockgendirective",
	Doc: `*_interface.go files must contain a //go:generate mockgen directive

Without the directive, go generate ./... does not regenerate the mock when the
interface changes. Tests then pass against a stale interface signature
indefinitely — the worst kind of false-green CI.`,
	Run: run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		name := pass.Fset.Position(file.Pos()).Filename
		base := filepath.Base(name)
		if !strings.HasSuffix(base, "_interface.go") {
			continue
		}
		if hasMockgenDirective(file.Comments) {
			continue
		}
		pass.Reportf(file.Package,
			"%s: missing //go:generate mockgen directive (interface files must regenerate mocks via go generate)",
			base)
	}
	return nil, nil
}

func hasMockgenDirective(groups []*ast.CommentGroup) bool {
	for _, g := range groups {
		for _, c := range g.List {
			text := strings.TrimSpace(c.Text)
			if strings.HasPrefix(text, "//go:generate") && strings.Contains(text, "mockgen") {
				return true
			}
		}
	}
	return false
}
