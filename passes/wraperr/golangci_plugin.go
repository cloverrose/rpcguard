package wraperr

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/cloverrose/rpcguard/pkg/logger"
)

func RegisterPlugin() {
	// https://golangci-lint.run/plugins/module-plugins/
	register.Plugin("rpc_wraperr", newPlugin)
}

func newPlugin(conf any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[settings](conf)
	if err != nil {
		return nil, err
	}

	return &plugin{settings: &s}, nil
}

type settings struct {
	Log                    logger.Config
	ReportMode             string
	IncludePackages        string
	ExcludePackages        string
	ExcludeFiles           string
	EnableErrGroupAnalyzer bool
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
	if p.settings.ReportMode != "" {
		ReportMode = p.settings.ReportMode
	}
	if p.settings.IncludePackages != "" {
		IncludePackages = p.settings.IncludePackages
	}
	if p.settings.ExcludePackages != "" {
		ExcludePackages = p.settings.ExcludePackages
	}
	if p.settings.ExcludeFiles != "" {
		ExcludeFiles = p.settings.ExcludeFiles
	}
	return []*analysis.Analyzer{
		Analyzer,
	}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

var _ register.LinterPlugin = &plugin{}
