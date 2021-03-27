package web

import (
	"github.com/kataras/iris"
	"management/core"
	"management/infra/base"
	"strconv"
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
	ctx.Request().Header.Set("user_type", strconv.Itoa(user.UserType))
	ctx.Request().Header.Set("user_id", strconv.FormatInt(user.Id,10))
	ctx.Next()

}
func employeeMeddle(ctx iris.Context) {
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
	if user.UserType != 2 {
		r := base.Res{
			Code:    base.ResError,
			Message: "该用户不是员工",
		}
		ctx.JSON(&r)
		return
	}
	ctx.Request().Header.Set("phone", user.Phone)
	ctx.Request().Header.Set("user_type", strconv.Itoa(user.UserType))
	ctx.Request().Header.Set("user_id", strconv.FormatInt(user.Id,10))
	ctx.Next()
}
