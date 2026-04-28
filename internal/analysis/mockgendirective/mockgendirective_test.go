package mockgendirective_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/nhalm/blueprint-vet/internal/analysis/mockgendirective"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), mockgendirective.Analyzer, "example")
}
