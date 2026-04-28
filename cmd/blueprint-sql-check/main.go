package main

import (
	"fmt"
	"os"

	"github.com/nhalm/blueprint-vet/internal/sqlcheck"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: blueprint-sql-check <dir-of-.sql-files>...")
		os.Exit(2)
	}
	failed := false
	for _, dir := range os.Args[1:] {
		findings := sqlcheck.Run(dir)
		for _, f := range findings {
			fmt.Println(f)
		}
		if len(findings) > 0 {
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}
