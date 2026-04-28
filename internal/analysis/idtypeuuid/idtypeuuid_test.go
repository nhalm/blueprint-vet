package idtypeuuid_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/nhalm/blueprint-vet/internal/analysis/idtypeuuid"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), idtypeuuid.Analyzer, "myapp/internal/models")
}
