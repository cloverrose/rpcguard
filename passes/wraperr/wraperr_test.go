package wraperr_test

import (
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/gostaticanalysis/testutil"

	"github.com/cloverrose/rpcguard/passes/wraperr"
)

func Test(t *testing.T) {
	t.Parallel()
	testdata := analysistest.TestData()
	testdata = testutil.WithModules(t, testdata, nil)
	wraperr.LogLevel = "INFO"
	wraperr.ReportMode = "BOTH"
	wraperr.EnableErrGroupAnalyzer = true
	pkgs := "a/a01core,a/a02phi,a/a03interface,a/a04closure,a/a05global,a/a06parameter,a/a07generics,a/a08import/a,a/a08import/includedpkg,a/a09cyclic,a/a10defer,eg/eg01core,eg/eg02generics,eg/eg03interface"
	wraperr.IncludePackages = "^(a/a01core|a/a02phi|a/a03interface|a/a04closure|a/a05global|a/a06parameter|a/a07generics|a/a08import/a|a/a08import/includedpkg|a/a09cyclic|a/a10defer|eg/eg01core|eg/eg02generics|eg/eg03interface)$"
	wraperr.ExcludePackages = "(.+/)?vendor$"
	analysistest.Run(t, testdata, wraperr.Analyzer, strings.Split(pkgs, ",")...)
}
