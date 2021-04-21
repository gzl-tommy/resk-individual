package gorpc

import (
	"gzl-tommy/resk-individual/infra"
	"gzl-tommy/resk-individual/infra/base"
)

type GoRpcApiStarter struct {
	infra.BaseStarter
}

func (g *GoRpcApiStarter) Init(ctx infra.StarterContext) {
	base.RpcRegister(new(EnvelopeRpc))
}