// Package sqlcheck applies blueprint conformance rules to .sql files.
//
// Two rules ship today: softdelete (R-9) and paginatedorderby (R-10). See
// proposals/blueprint-vet.md for rule statements and rationale.
package sqlcheck

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Block is a single named query: a `-- name: <Name> :<type>` annotation and
// the SQL text following it up to the next annotation or EOF.
type Block struct {
	Name string
	Type string
	SQL  string
	File string
	Line int // 1-indexed line of the annotation
}

// Finding is a single rule violation.
type Finding struct {
	File    string
	Line    int
	Rule    string
	Message string
}

func (f Finding) String() string {
	if f.Line == 0 {
		return fmt.Sprintf("%s: %s: %s", f.File, f.Rule, f.Message)
	}
	return fmt.Sprintf("%s:%d: %s: %s", f.File, f.Line, f.Rule, f.Message)
}

var blockHeaderRe = regexp.MustCompile(`^\s*--\s*name:\s*(\w+)\s*:(\w+)\s*$`)

// Parse reads a .sql file and returns its query blocks.
func Parse(file string) ([]Block, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var blocks []Block
	var current *Block
	var sqlLines []string
	flush := func() {
		if current != nil {
			current.SQL = strings.Join(sqlLines, "\n")
			blocks = append(blocks, *current)
		}
		current = nil
		sqlLines = nil
	}
	for i, line := range strings.Split(string(data), "\n") {
		if m := blockHeaderRe.FindStringSubmatch(line); m != nil {
			flush()
			current = &Block{Name: m[1], Type: m[2], File: file, Line: i + 1}
			continue
		}
		if current != nil {
			sqlLines = append(sqlLines, line)
		}
	}
	flush()
	return blocks, nil
}

// Run walks dir for .sql files and returns sorted findings across all rules.
func Run(dir string) []Finding {
	var findings []Finding
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			findings = append(findings, Finding{File: path, Rule: "walk", Message: err.Error()})
			return nil
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".sql") {
			return nil
		}
		blocks, err := Parse(path)
		if err != nil {
			findings = append(findings, Finding{File: path, Rule: "parse", Message: err.Error()})
			return nil
		}
		for _, b := range blocks {
			findings = append(findings, SoftDelete(b)...)
			findings = append(findings, PaginatedOrderBy(b)...)
		}
		return nil
	})
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].File != findings[j].File {
			return findings[i].File < findings[j].File
		}
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
		}
		return findings[i].Rule < findings[j].Rule
	})
	return findings
}

var softDeleteOptOutRe = regexp.MustCompile(`(?i)(IncludingDeleted|Audit|Trash|AllVersions)$`)

// SoftDelete (R-9) flags read queries (`:one`, `:many`, `:paginated`) whose SQL
// lacks a `deleted_at` token, unless the query name opts out via convention.
func SoftDelete(b Block) []Finding {
	if b.Type != "one" && b.Type != "many" && b.Type != "paginated" {
		return nil
	}
	upper := strings.ToUpper(b.SQL)
	if !strings.Contains(upper, "SELECT") {
		return nil
	}
	if strings.Contains(upper, "DELETED_AT") {
		return nil
	}
	if softDeleteOptOutRe.MatchString(b.Name) {
		return nil
	}
	return []Finding{{
		File: b.File, Line: b.Line, Rule: "softdelete",
		Message: fmt.Sprintf("%s: read query missing `deleted_at IS NULL` filter; if including soft-deleted rows is intentional, name it *IncludingDeleted, *Audit, *Trash, or *AllVersions", b.Name),
	}}
}

var orderByRe = regexp.MustCompile(`(?i)\bORDER\s+BY\b`)

// PaginatedOrderBy (R-10) flags `:paginated` queries that lack an ORDER BY.
// skimatik requires a stable sort to construct the cursor.
func PaginatedOrderBy(b Block) []Finding {
	if b.Type != "paginated" {
		return nil
	}
	if orderByRe.MatchString(b.SQL) {
		return nil
	}
	return []Finding{{
		File: b.File, Line: b.Line, Rule: "paginatedorderby",
		Message: fmt.Sprintf("%s: :paginated query must contain ORDER BY (skimatik builds the cursor from it)", b.Name),
	}}
}
