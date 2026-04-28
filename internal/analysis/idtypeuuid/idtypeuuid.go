// Package idtypeuuid implements an analyzer that flags ID-typed struct
// fields in internal/models that are not uuid.UUID.
package idtypeuuid

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "idtypeuuid",
	Doc: `model ID fields must be uuid.UUID

The blueprint's ID strategy is "internal uuid.UUID, wire string." Domain
models that hold IDs as strings short-circuit the boundary, push UUID parsing
into queries, and produce worse error messages on bad input.

Fields under /internal/models whose name is ID or ends in ID must be
uuid.UUID or *uuid.UUID.`,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	if !strings.Contains(pass.Pkg.Path(), "/internal/models") {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.StructType)(nil)}, func(n ast.Node) {
		st := n.(*ast.StructType)
		if st.Fields == nil {
			return
		}
		for _, field := range st.Fields.List {
			for _, name := range field.Names {
				if !isIDName(name.Name) {
					continue
				}
				if isUUIDType(pass, field.Type) {
					continue
				}
				pass.Reportf(field.Pos(),
					"field %s must be uuid.UUID (or *uuid.UUID); domain models keep IDs as UUIDs, not strings",
					name.Name)
			}
		}
	})
	return nil, nil
}

func isIDName(name string) bool {
	return name == "ID" || strings.HasSuffix(name, "ID")
}

func isUUIDType(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return false
	}
	return named.Obj().Pkg().Path() == "github.com/google/uuid" && named.Obj().Name() == "UUID"
}
