package nofmtprint_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/nhalm/blueprint-vet/analysis/nofmtprint"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), nofmtprint.Analyzer,
		"myapp/internal/api",
		"myapp/internal/config",
		"myapp/cmd/server",
	)
}
