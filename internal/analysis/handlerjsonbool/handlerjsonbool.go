// Package handlerjsonbool implements an analyzer that flags calls to
// chikit.JSON / chikit.Query whose bool return is discarded.
package handlerjsonbool

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "handlerjsonbool",
	Doc: `chikit.JSON / chikit.Query results must be checked

The bool return signals whether validation/binding failed. If the call result
is discarded, validation errors silently set the response but the handler
keeps running — calling the service with a zero-value request struct.

Use:
    if !chikit.JSON(r, &dst) { return }`,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	if !strings.Contains(pass.Pkg.Path(), "/internal/api") {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}
		call := n.(*ast.CallExpr)
		if !isChikitGuardCall(pass, call) {
			return true
		}
		if len(stack) < 2 {
			return true
		}
		parent := stack[len(stack)-2]
		switch p := parent.(type) {
		case *ast.ExprStmt:
			pass.Reportf(call.Pos(),
				"chikit.%s result discarded; use `if !chikit.%s(...) { return }`",
				call.Fun.(*ast.SelectorExpr).Sel.Name,
				call.Fun.(*ast.SelectorExpr).Sel.Name)
		case *ast.AssignStmt:
			if allBlank(p.Lhs) {
				pass.Reportf(call.Pos(),
					"chikit.%s result discarded via _; use `if !chikit.%s(...) { return }`",
					call.Fun.(*ast.SelectorExpr).Sel.Name,
					call.Fun.(*ast.SelectorExpr).Sel.Name)
			}
		}
		return true
	})
	return nil, nil
}

func isChikitGuardCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if sel.Sel.Name != "JSON" && sel.Sel.Name != "Query" {
		return false
	}
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	pkgName, ok := pass.TypesInfo.Uses[pkgIdent].(*types.PkgName)
	if !ok {
		return false
	}
	// Match any chikit import path (real lives at github.com/nhalm/chikit; tests
	// stub it locally). Last path segment is enough.
	imported := pkgName.Imported().Path()
	return imported == "chikit" || strings.HasSuffix(imported, "/chikit")
}

func allBlank(lhs []ast.Expr) bool {
	for _, e := range lhs {
		ident, ok := e.(*ast.Ident)
		if !ok || ident.Name != "_" {
			return false
		}
	}
	return len(lhs) > 0
}
