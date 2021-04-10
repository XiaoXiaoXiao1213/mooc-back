package web

import (
	_ "management/apis/web"
	_ "management/core/users"
	"management/infra"
	"management/infra/base"
	"management/jobs"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.MongoStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.HookStarter{})
	infra.Register(&jobs.CommendJobStarter{})
}
