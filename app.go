package resk

import (
	_ "management/apis/web"
	_ "management/core/users"
	"management/infra"
	"management/infra/base"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.HookStarter{})
}
