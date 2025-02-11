package wraperr

import (
	"github.com/cloverrose/rpcguard/pkg/filter"
)

// LogLevel is configuration of logging level.
// Available options are DEBUG, INFO, WARN, ERROR
var LogLevel = "INFO"

const (
	reportModeReturn   = "RETURN"
	reportModeFunction = "FUNCTION"
	reportModeBoth     = "BOTH"
)

// ReportMode is configuration of how to report violation.
// Available options are RETURN, FUNCTION, BOTH.
// - RETURN: Report violations at return instruction level. It provides detailed information, and you can find violation easily.
// - FUNCTION: Report violations at function level. It is convenient if you use nolint.
// - BOTH: Report violations both levels. This mode is mainly for unit test of wraperr itself.
var ReportMode = reportModeReturn

// IncludePackages is configuration which packages should be included.
// Multiple Packages can be specified by using commas (,).
// e.g. github.com/foo/bar/a/includedpkg/hello,github.com/foo/bar/common
//
// [Explain with example]
// There are two packages,
//   - package hello (path=github.com/foo/bar/a/includedpkg/hello) <- Let's say hello defines RPC method Hello.
//   - package errutil (path=github.com/foo/bar/common/errutil) <- Let's say errutil provides func that converts error with connect.NewError.
//
// hello's Hello calls errutil's ConvertErr.
// Regardless of this configuration, wraperr understands that Hello() is calling ConvertErr().
// If this configuration doesn't include "github.com/foo/bar/common/errutil", wraperr doesn't know whether ConvertErr()'s error is connect.NewError or not.
// Thus, wraperr reports, Hello() can return non connect.NewError errors.
// So, includes enough package path prefix is important for precise reports.
// On the other hand, extending the search area is computationally time-consuming.
// Narrowing down the search to only relevant Packages can improve efficiency.
var IncludePackages = ""

// ExcludePackages is configuration which packages should be excluded.
// This is useful to exclude vendor.
var ExcludePackages = ""

// ExcludeFiles is configuration which files should be excluded.
// This is useful to exclude test file, generated files.
// To set the same value with the default config, use this command line argument.
// -rpc_wraperr.ExcludeFiles='.+_test\.go,.+\.connect\.go'
var ExcludeFiles = `.+_test\.go,.+\.connect\.go`

// EnableErrGroupAnalyzer is configuration whether enable ErrGroupAnalyzer.
// Default is true and recommend to keep true.
// The reason for this configuration item is to share with users that
// wraperr can be incorrect in a false positive sense for propagation of errors that it does not support.
// errgroup is supported, but others such as hashicorp/go-multierror is not supported.
var EnableErrGroupAnalyzer = true

var (
	packageFilter *filter.Filter
	fileFilter    *filter.Filter
)
