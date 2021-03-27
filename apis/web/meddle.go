package web

import (
	"github.com/kataras/iris"
	"management/core"
	"management/infra/base"
)

func loginMeddle(ctx iris.Context) {
	token := ctx.GetHeader("authorization")
	if token == "" {
		r := base.Res{
			Code:    base.ResError,
			Message: "用户未登录",
		}
		ctx.JSON(&r)
	}
	user, err := core.ParseToken(token)
	if err != nil {
		r := base.Res{
			Code:    base.ResError,
			Message: "token失效，请重新登陆",
		}
		ctx.JSON(&r)
	}
	ctx.Request().Header.Set("phone", user.Phone)
	ctx.Request().Header.Set("user_type", string(user.UserType))
	ctx.Request().Header.Set("user_id", string(user.Id))
	ctx.Next()

}
func employeeMeddle(ctx iris.Context) {
	userType := ctx.GetHeader("user_type")
	if userType != "2" {
		r := base.Res{
			Code:    base.ResError,
			Message: "该用户不是员工",
		}
		ctx.JSON(&r)
	}
	ctx.Next()

}
