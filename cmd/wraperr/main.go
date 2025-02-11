package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/cloverrose/rpcguard/passes/wraperr"
)

func main() {
	unitchecker.Main(wraperr.Analyzer)
}
