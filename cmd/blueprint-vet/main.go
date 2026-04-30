package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/nhalm/blueprint-vet/analysis/apperroralias"
	"github.com/nhalm/blueprint-vet/analysis/errortranslate"
	"github.com/nhalm/blueprint-vet/analysis/handlerjsonbool"
	"github.com/nhalm/blueprint-vet/analysis/idtypeuuid"
	"github.com/nhalm/blueprint-vet/analysis/layerdirection"
	"github.com/nhalm/blueprint-vet/analysis/mockgendirective"
	"github.com/nhalm/blueprint-vet/analysis/nofmtprint"
	"github.com/nhalm/blueprint-vet/analysis/nojsonencode"
	"github.com/nhalm/blueprint-vet/analysis/nowriteheader"
	"github.com/nhalm/blueprint-vet/analysis/repoexecutor"
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
