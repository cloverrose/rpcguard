package callvalidate

import (
	"fmt"
	"strings"

	"github.com/cloverrose/rpcguard/pkg/filter"
	"github.com/cloverrose/rpcguard/pkg/logger"
)

// log related configuration.
var LogConfig = logger.Config{
	Level:  "INFO",
	File:   "",
	Format: "json",
}

// ExcludeFiles is configuration which files should be excluded.
// This is useful to exclude test file, generated files.
// To set the same value with the default config, use this command line argument.
// -rpc_callvalidate.ExcludeFiles='.+_test\.go,.+\.connect\.go'
var ExcludeFiles = `.+_test\.go,.+\.connect\.go`

// ValidateMethods is configuration which methods should be called.
// Default is "buf.build/go/protovalidate:Validate,github.com/bufbuild/protovalidate-go:Validate"
// Package and Method join with `:`
// You can specify multiple methods by using `,` separated value.
var ValidateMethods = "buf.build/go/protovalidate:Validate,github.com/bufbuild/protovalidate-go:Validate"

var fileFilter *filter.Filter

type Method struct {
	packagePath string
	name        string
}

func parseMethods(input string) ([]Method, error) {
	values := strings.Split(input, ",")
	methods := make([]Method, len(values))
	for i, value := range values {
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid method format: %s", value)
		}
		methods[i] = Method{parts[0], parts[1]}
	}
	return methods, nil
}
