package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/nhalm/blueprint-vet/internal/analysis/apperroralias"
	"github.com/nhalm/blueprint-vet/internal/analysis/errortranslate"
	"github.com/nhalm/blueprint-vet/internal/analysis/handlerjsonbool"
	"github.com/nhalm/blueprint-vet/internal/analysis/idtypeuuid"
	"github.com/nhalm/blueprint-vet/internal/analysis/layerdirection"
	"github.com/nhalm/blueprint-vet/internal/analysis/mockgendirective"
	"github.com/nhalm/blueprint-vet/internal/analysis/nofmtprint"
	"github.com/nhalm/blueprint-vet/internal/analysis/nojsonencode"
	"github.com/nhalm/blueprint-vet/internal/analysis/nowriteheader"
	"github.com/nhalm/blueprint-vet/internal/analysis/repoexecutor"
)

func main() {
	multichecker.Main(
		apperroralias.Analyzer,
		errortranslate.Analyzer,
		handlerjsonbool.Analyzer,
		idtypeuuid.Analyzer,
		layerdirection.Analyzer,
		mockgendirective.Analyzer,
		nofmtprint.Analyzer,
		nojsonencode.Analyzer,
		nowriteheader.Analyzer,
		repoexecutor.Analyzer,
	)
}
