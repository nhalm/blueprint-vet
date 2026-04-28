// Package repoexecutor implements an analyzer that flags repository
// methods calling generated repo methods with a direct field access
// (e.g. r.db) instead of routing through executorFromContext.
package repoexecutor

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "repoexecutor",
	Doc: `repository methods must route generated calls through executorFromContext

Generated repo methods take an Executor. Passing r.db works outside
transactions but bypasses the active *pgxkit.Tx when called inside BeginTx.
The blueprint's executorFromContext resolves the right executor from context.
Skipping it is the classic "transactional service silently doesn't transact"
bug.

Use:
    row, err := r.GetProductByAccountAndID(ctx, executorFromContext(ctx, r.db), id)`,
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
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if !isGeneratedMethod(pass, sel) {
				return true
			}
			if len(call.Args) < 2 {
				return true
			}
			if _, ok := call.Args[1].(*ast.SelectorExpr); ok {
				pass.Reportf(call.Args[1].Pos(),
					"pass executorFromContext(ctx, r.db) instead of a direct field; %s.%s bypasses the active transaction",
					exprString(sel.X), sel.Sel.Name)
			}
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

func isGeneratedMethod(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	selection := pass.TypesInfo.Selections[sel]
	if selection == nil {
		return false
	}
	fn, ok := selection.Obj().(*types.Func)
	if !ok {
		return false
	}
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Recv() == nil {
		return false
	}
	recvType := sig.Recv().Type()
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}
	named, ok := recvType.(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return false
	}
	return strings.Contains(named.Obj().Pkg().Path(), "/internal/repository/generated")
}

func exprString(e ast.Expr) string {
	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}
	return "<recv>"
}
