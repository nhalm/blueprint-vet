package nowriteheader_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/nhalm/blueprint-vet/internal/analysis/nowriteheader"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), nowriteheader.Analyzer, "myapp/internal/api")
}
