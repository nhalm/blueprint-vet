package layerdirection_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/nhalm/blueprint-vet/analysis/layerdirection"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), layerdirection.Analyzer,
		"myapp/internal/models/good",
		"myapp/internal/models/bad",
		"myapp/internal/api/good",
		"myapp/internal/api/bad",
	)
}
