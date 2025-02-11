package rpcguard

import (
	"github.com/cloverrose/rpcguard/passes/callvalidate"
	"github.com/cloverrose/rpcguard/passes/wraperr"
)

func init() {
	callvalidate.RegisterPlugin()
	wraperr.RegisterPlugin()
}
