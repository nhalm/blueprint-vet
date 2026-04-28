// Package nowriteheader implements an analyzer that flags direct
// http.ResponseWriter.WriteHeader calls in code under internal/api.
package nowriteheader

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "nowriteheader",
	Doc: `disallow direct WriteHeader calls in internal/api

Handlers must use chikit.SetResponse or chikit.SetError. The chikit.Handler
middleware owns deferred response writing — calling WriteHeader directly
bypasses canonical logging, error envelope formatting, and the response-state
mutex.`,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	if !strings.Contains(pass.Pkg.Path(), "/internal/api") {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "WriteHeader" {
			return
		}
		pass.Reportf(call.Pos(),
			"use chikit.SetResponse or chikit.SetError; chikit.Handler middleware owns response writing")
	})
	return nil, nil
}
