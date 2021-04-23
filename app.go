package resk_individual

import (
	"gzl-tommy/resk-individual/apis/gorpc"
	_ "gzl-tommy/resk-individual/apis/web"
	_ "gzl-tommy/resk-individual/core/accounts"
	_ "gzl-tommy/resk-individual/core/envelopes"
	"gzl-tommy/resk-individual/infra"
	"gzl-tommy/resk-individual/infra/base"
	"gzl-tommy/resk-individual/jobs"
	_ "gzl-tommy/resk-individual/public/ui"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.GoRPCStarter{})
	infra.Register(&gorpc.GoRpcApiStarter{})
	infra.Register(&jobs.RefundExpiredJobStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.EurekaStarter{})
	infra.Register(&base.HookStarter{})
}
