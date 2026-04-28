// Package nojsonencode implements an analyzer that flags
// json.NewEncoder(w).Encode(...) calls in code under internal/api,
// where w flows from an http.ResponseWriter parameter.
package nojsonencode

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "nojsonencode",
	Doc: `disallow json.NewEncoder(w).Encode(...) on http.ResponseWriter in internal/api

Use chikit.SetResponse / chikit.SetError. Direct encoding produces inconsistent
error envelopes, skips canonical logging, and races with chikit.Handler's
deferred response writer.`,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	if !strings.Contains(pass.Pkg.Path(), "/internal/api") {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		// Looking for: <something>.Encode(<arg>) where <something> is
		// json.NewEncoder(w) and w's type is http.ResponseWriter (or implements it).
		encodeCall := n.(*ast.CallExpr)
		encodeSel, ok := encodeCall.Fun.(*ast.SelectorExpr)
		if !ok || encodeSel.Sel.Name != "Encode" {
			return
		}
		newEnc, ok := encodeSel.X.(*ast.CallExpr)
		if !ok {
			return
		}
		newEncSel, ok := newEnc.Fun.(*ast.SelectorExpr)
		if !ok || newEncSel.Sel.Name != "NewEncoder" {
			return
		}
		pkgIdent, ok := newEncSel.X.(*ast.Ident)
		if !ok {
			return
		}
		pkgName, ok := pass.TypesInfo.Uses[pkgIdent].(*types.PkgName)
		if !ok || pkgName.Imported().Path() != "encoding/json" {
			return
		}
		if len(newEnc.Args) != 1 {
			return
		}
		if !isResponseWriter(pass, newEnc.Args[0]) {
			return
		}
		pass.Reportf(encodeCall.Pos(),
			"use chikit.SetResponse / chikit.SetError instead of json.NewEncoder(w).Encode")
	})
	return nil, nil
}

func isResponseWriter(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	named, ok := t.(*types.Named)
	if !ok {
		// http.ResponseWriter is an interface — check the underlying type's name path.
		if iface, ok := t.Underlying().(*types.Interface); ok {
			return interfaceMatchesResponseWriter(iface)
		}
		return false
	}
	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}
	return obj.Pkg().Path() == "net/http" && obj.Name() == "ResponseWriter"
}

func interfaceMatchesResponseWriter(iface *types.Interface) bool {
	// Heuristic: http.ResponseWriter has Header(), Write([]byte) (int, error), WriteHeader(int).
	hasHeader, hasWrite, hasWriteHeader := false, false, false
	for i := 0; i < iface.NumMethods(); i++ {
		switch iface.Method(i).Name() {
		case "Header":
			hasHeader = true
		case "Write":
			hasWrite = true
		case "WriteHeader":
			hasWriteHeader = true
		}
	}
	return hasHeader && hasWrite && hasWriteHeader
}
