package callvalidate

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/cloverrose/rpcguard/pkg/logger"
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
	Log             logger.Config
	ExcludeFiles    string
	ValidateMethods string
}

type plugin struct {
	settings *settings
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	if p.settings.Log.Level != "" {
		LogConfig.Level = p.settings.Log.Level
	}
	if p.settings.Log.File != "" {
		LogConfig.File = p.settings.Log.File
	}
	if p.settings.Log.Format != "" {
		LogConfig.Format = p.settings.Log.Format
	}
	if p.settings.ExcludeFiles != "" {
		ExcludeFiles = p.settings.ExcludeFiles
	}
	if p.settings.ValidateMethods != "" {
		ValidateMethods = p.settings.ValidateMethods
	}
	return []*analysis.Analyzer{
		Analyzer,
	}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

var _ register.LinterPlugin = &plugin{}
