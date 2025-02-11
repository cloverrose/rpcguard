package callvalidate

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func RegisterPlugin() {
	// https://golangci-lint.run/plugins/module-plugins/
	register.Plugin("rpc_callvalidate", newPlugin)
}

func newPlugin(conf any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[settings](conf)
	if err != nil {
		return nil, err
	}

	return &plugin{settings: &s}, nil
}

type settings struct {
	LogLevel     string
	ExcludeFiles string
}

type plugin struct {
	settings *settings
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	if p.settings.LogLevel != "" {
		LogLevel = p.settings.LogLevel
	}
	if p.settings.ExcludeFiles != "" {
		ExcludeFiles = p.settings.ExcludeFiles
	}
	return []*analysis.Analyzer{
		Analyzer,
	}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

var _ register.LinterPlugin = &plugin{}
