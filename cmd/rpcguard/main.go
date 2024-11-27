package main

import (
	"log"
	"os"

	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/cloverrose/rpcguard"
)

func main() {
	log.SetOutput(os.Stdout)
	unitchecker.Main(rpcguard.Analyzer)
}
