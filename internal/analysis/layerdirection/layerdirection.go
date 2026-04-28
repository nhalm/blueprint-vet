// Package layerdirection implements an analyzer that enforces the
// blueprint's layer-direction rules across internal/ packages.
//
//   internal/models cannot import internal/repository, internal/service, or internal/api.
//   internal/api cannot import internal/repository (consume domain via the service interface).
package layerdirection

import (
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "layerdirection",
	Doc: `enforce blueprint layer direction across internal/ packages

internal/models is the deepest layer and cannot import upward.
internal/api consumes domain through the service interface, not the repository directly.`,
	Run: run,
}

type rule struct {
	source string   // package path substring that activates this rule
	deny   []string // import path substrings forbidden when the rule is active
	desc   string
}

var rules = []rule{
	{
		source: "/internal/models",
		deny:   []string{"/internal/repository", "/internal/service", "/internal/api"},
		desc:   "models cannot import repository, service, or api",
	},
	{
		source: "/internal/api",
		deny:   []string{"/internal/repository"},
		desc:   "api cannot import repository directly; consume via the service interface",
	},
}

func run(pass *analysis.Pass) (any, error) {
	pkgPath := pass.Pkg.Path()
	for _, r := range rules {
		if !strings.Contains(pkgPath, r.source) {
			continue
		}
		for _, file := range pass.Files {
			for _, imp := range file.Imports {
				path, err := strconv.Unquote(imp.Path.Value)
				if err != nil {
					continue
				}
				for _, deny := range r.deny {
					if strings.Contains(path, deny) {
						pass.Reportf(imp.Pos(), "%s (imports %q)", r.desc, path)
						break
					}
				}
			}
		}
	}
	return nil, nil
}
