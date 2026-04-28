package handlerjsonbool_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/nhalm/blueprint-vet/internal/analysis/handlerjsonbool"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), handlerjsonbool.Analyzer, "myapp/internal/api")
}
