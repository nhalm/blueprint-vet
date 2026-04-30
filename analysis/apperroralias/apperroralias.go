// Package apperroralias implements an analyzer that flags imports of
// internal/errors packages that don't use the alias `apperrors`.
package apperroralias

import (
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "apperroralias",
	Doc: `internal/errors must be imported as apperrors

The package name "errors" collides with stdlib. The blueprint codified
"apperrors" as the alias. Bare imports shadow stdlib errors; other aliases
(apierrors, domerr, myerrors) produce drift across files in the same codebase.

Bad:

	import "myapp/internal/errors"
	import apierrors "myapp/internal/errors"

Good:

	import apperrors "myapp/internal/errors"`,
	Run: run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, imp := range file.Imports {
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				continue
			}
			if !strings.HasSuffix(path, "/internal/errors") {
				continue
			}
			if imp.Name == nil {
				pass.Reportf(imp.Pos(),
					"import %q must use alias `apperrors` (bare import shadows stdlib errors)", path)
				continue
			}
			if imp.Name.Name != "apperrors" {
				pass.Reportf(imp.Name.Pos(),
					"import %q must use alias `apperrors`, not %q", path, imp.Name.Name)
			}
		}
	}
	return nil, nil
}
