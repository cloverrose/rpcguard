package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/cloverrose/rpcguard/passes/callvalidate"
)

func main() {
	unitchecker.Main(callvalidate.Analyzer)
}
