// Package errortranslate implements an analyzer that flags repository
// methods returning errors from generated repo calls without wrapping
// them through translateError.
package errortranslate

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "errortranslate",
	Doc: `repository methods must wrap generated errors through translateError

Generated methods return raw pgx errors (pgx.ErrNoRows, unique-constraint
violations, etc.). Repositories must convert them to apperrors sentinels
via translateError so the service layer can map them to domain errors.

Returning err directly from a generated call leaks pgx-specific errors past
the repository boundary, breaking error mapping in services and the chikit
handler error envelope.

Bad:

	row, err := r.GetProductByID(ctx, executorFromContext(ctx, r.db), id)
	if err != nil {
		return Product{}, err
	}

Good:

	row, err := r.GetProductByID(ctx, executorFromContext(ctx, r.db), id)
	if err != nil {
		return Product{}, translateError(err)
	}`,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	if !packageImportsGenerated(pass.Pkg) {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Nodes([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node, push bool) bool {
		if !push {
			return false
		}
		fn := n.(*ast.FuncDecl)
		if fn.Recv == nil || len(fn.Recv.List) == 0 || fn.Body == nil {
			return false
		}
		if !receiverIsRepository(fn.Recv) {
			return false
		}
		if !methodCallsGenerated(pass, fn.Body) {
			// Method doesn't touch the generated layer at all.
			return false
		}
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			ret, ok := node.(*ast.ReturnStmt)
			if !ok || len(ret.Results) == 0 {
				return true
			}
			last := ret.Results[len(ret.Results)-1]
			if !isBareErrIdent(last) {
				return true
			}
			pass.Reportf(last.Pos(),
				"wrap %s through translateError; bare returns leak pgx errors past the repository boundary",
				last.(*ast.Ident).Name)
			return true
		})
		return false
	})
	return nil, nil
}

func packageImportsGenerated(pkg *types.Package) bool {
	for _, imp := range pkg.Imports() {
		if strings.Contains(imp.Path(), "/internal/repository/generated") {
			return true
		}
	}
	return false
}

func receiverIsRepository(recv *ast.FieldList) bool {
	if len(recv.List) == 0 {
		return false
	}
	t := recv.List[0].Type
	if star, ok := t.(*ast.StarExpr); ok {
		t = star.X
	}
	ident, ok := t.(*ast.Ident)
	if !ok {
		return false
	}
	return strings.HasSuffix(ident.Name, "Repository")
}

func methodCallsGenerated(pass *analysis.Pass, body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		selection := pass.TypesInfo.Selections[sel]
		if selection == nil {
			return true
		}
		fn, ok := selection.Obj().(*types.Func)
		if !ok {
			return true
		}
		sig, ok := fn.Type().(*types.Signature)
		if !ok || sig.Recv() == nil {
			return true
		}
		recvType := sig.Recv().Type()
		if ptr, ok := recvType.(*types.Pointer); ok {
			recvType = ptr.Elem()
		}
		named, ok := recvType.(*types.Named)
		if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
			return true
		}
		if strings.Contains(named.Obj().Pkg().Path(), "/internal/repository/generated") {
			found = true
			return false
		}
		return true
	})
	return found
}

// isBareErrIdent matches identifiers that look like error variables
// (err, dbErr, queryErr, etc.) — not function calls, selector expressions,
// or composite literals. Wrapping err in translateError or any other call
// produces an *ast.CallExpr, which this function returns false for.
func isBareErrIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	name := ident.Name
	return name == "err" || strings.HasSuffix(name, "Err") || strings.HasSuffix(name, "Error")
}
