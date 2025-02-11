package rpcguard

import (
	"github.com/cloverrose/rpcguard/passes/callvalidate"
)

func init() {
	callvalidate.RegisterPlugin()
}
