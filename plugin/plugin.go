// Package plugin is the golangci-lint Module Plugin entry point for
// blueprint-vet. It registers the same set of analyzers exposed by the
// standalone blueprint-vet binary so they can be enabled from a consumer's
// .golangci.yml via golangci-lint's Module Plugin System.
//
// See https://golangci-lint.run/docs/plugins/module-plugins/ for how a
// custom golangci-lint binary is built that imports this package.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

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

func init() {
	register.Plugin("blueprint-vet", New)
}

// Settings is the per-plugin configuration block from .golangci.yml.
//
// It is intentionally empty for now: blueprint-vet's analyzers don't yet
// accept user-tunable knobs. Adding fields here later is backwards
// compatible — golangci-lint passes the raw `settings:` map through
// register.DecodeSettings, which will tolerate the new fields once they
// exist.
type Settings struct{}

// New constructs a LinterPlugin from the raw settings golangci-lint passes
// in. It satisfies register.NewPlugin.
func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}
	return &blueprintVetPlugin{settings: s}, nil
}

type blueprintVetPlugin struct {
	settings Settings
}

// BuildAnalyzers returns the full set of blueprint-vet analyzers. It mirrors
// the list registered by cmd/blueprint-vet/main.go so the plugin and the
// standalone binary stay in lockstep.
func (p *blueprintVetPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
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
	}, nil
}

// GetLoadMode reports the load mode required by the analyzers. All current
// blueprint-vet rules consult type information (e.g. resolving named types,
// method receivers, imported package paths), so we request typesinfo.
func (p *blueprintVetPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
