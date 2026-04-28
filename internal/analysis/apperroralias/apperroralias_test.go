package apperroralias_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/nhalm/blueprint-vet/internal/analysis/apperroralias"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), apperroralias.Analyzer,
		"myapp/internal/service/good",
		"myapp/internal/service/badbare",
		"myapp/internal/service/badalias",
	)
}
