package callvalidate_test

import (
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/gostaticanalysis/testutil"

	"github.com/cloverrose/rpcguard/passes/callvalidate"
)

func Test(t *testing.T) {
	t.Parallel()
	testdata := analysistest.TestData()
	testdata = testutil.WithModules(t, testdata, nil)
	callvalidate.LogConfig.Level = "INFO"
	callvalidate.ValidateMethods = "github.com/bufbuild/protovalidate-go:Validate,a:customValidate"
	pkgs := "a"
	analysistest.Run(t, testdata, callvalidate.Analyzer, strings.Split(pkgs, ",")...)
}
