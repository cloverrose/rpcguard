package rpcguard

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/plugin-module-register/register"
)

func init() {
	RegisterPlugin()
}

func RegisterPlugin() {
	// https://golangci-lint.run/plugins/module-plugins/
	register.Plugin("rpcguard", newPlugin)
}

func newPlugin(conf any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[settings](conf)
	if err != nil {
		return nil, err
	}

	return &plugin{settings: &s}, nil
}

type settings struct {
	Verbose string
}

type plugin struct {
	settings *settings
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	if p.settings.Verbose == "true" {
		Verbose = true
	}
	return []*analysis.Analyzer{
		Analyzer,
	}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

var _ register.LinterPlugin = &plugin{}
