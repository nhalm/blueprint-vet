package sqlcheck_test

import (
	"strings"
	"testing"

	"github.com/nhalm/blueprint-vet/internal/sqlcheck"
)

func TestRun(t *testing.T) {
	findings := sqlcheck.Run("testdata/queries")

	want := []struct {
		rule    string
		query   string
		lineGTE int
	}{
		{"softdelete", "ListProducts", 1},
		{"paginatedorderby", "ListProductsPage", 1},
		{"softdelete", "ListProductsPage", 1},
	}

	if len(findings) != len(want) {
		t.Fatalf("got %d findings, want %d:\n%s", len(findings), len(want), formatFindings(findings))
	}
	for i, w := range want {
		f := findings[i]
		if f.Rule != w.rule {
			t.Errorf("finding %d: rule=%q want %q", i, f.Rule, w.rule)
		}
		if !strings.Contains(f.Message, w.query) {
			t.Errorf("finding %d: message %q does not name query %q", i, f.Message, w.query)
		}
		if f.Line < w.lineGTE {
			t.Errorf("finding %d: line=%d, want >= %d", i, f.Line, w.lineGTE)
		}
	}
}

func formatFindings(fs []sqlcheck.Finding) string {
	var b strings.Builder
	for _, f := range fs {
		b.WriteString(f.String())
		b.WriteByte('\n')
	}
	return b.String()
}

func TestParseAttributes(t *testing.T) {
	blocks, err := sqlcheck.Parse("testdata/queries/products.sql")
	if err != nil {
		t.Fatal(err)
	}
	want := []struct {
		name string
		typ  string
	}{
		{"GetProductByID", "one"},
		{"ListProducts", "many"},
		{"ListProductsPage", "paginated"},
		{"ListProductsByName", "paginated"},
		{"ListProductsIncludingDeleted", "many"},
		{"GetProductAudit", "many"},
		{"SoftDeleteProduct", "exec"},
		{"InsertProduct", "one"},
	}
	if len(blocks) != len(want) {
		t.Fatalf("got %d blocks, want %d", len(blocks), len(want))
	}
	for i, w := range want {
		if blocks[i].Name != w.name || blocks[i].Type != w.typ {
			t.Errorf("block %d: %s :%s, want %s :%s", i, blocks[i].Name, blocks[i].Type, w.name, w.typ)
		}
	}
}
