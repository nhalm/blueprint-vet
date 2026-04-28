// Package nofmtprint implements an analyzer that flags fmt.Print*,
// fmt.Fprint*, and fmt.Sprint* calls outside whitelisted paths.
package nofmtprint

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "nofmtprint",
	Doc: `runtime logging must go through canonlog, not fmt.Print*

fmt.Print* lands on stdout, bypassing Datadog. Allowed only in:
  - cmd/<app>/         (CLI feedback, after the canonlog event)
  - internal/config/   (pre-canonlog config-load errors)

Bad:

	fmt.Println("user created")

Good:

	canonlog.InfoAdd(ctx, "event", "user_created")`,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var bannedNames = map[string]bool{
	"Print":    true,
	"Println":  true,
	"Printf":   true,
	"Fprint":   true,
	"Fprintln": true,
	"Fprintf":  true,
	"Sprint":   true,
	"Sprintln": true,
	"Sprintf":  true,
}

func run(pass *analysis.Pass) (any, error) {
	path := pass.Pkg.Path()
	if strings.Contains(path, "/cmd/") || strings.HasSuffix(path, "/internal/config") || strings.Contains(path, "/internal/config/") {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		if !bannedNames[sel.Sel.Name] {
			return
		}
		pkgIdent, ok := sel.X.(*ast.Ident)
		if !ok {
			return
		}
		pkgName, ok := pass.TypesInfo.Uses[pkgIdent].(*types.PkgName)
		if !ok || pkgName.Imported().Path() != "fmt" {
			return
		}
		pass.Reportf(call.Pos(),
			"use canonlog instead of fmt.%s; runtime logs must reach Datadog", sel.Sel.Name)
	})
	return nil, nil
}
