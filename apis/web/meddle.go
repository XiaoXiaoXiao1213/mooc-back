package web

import (
	"github.com/kataras/iris"
	"management/core/common"
	"management/core/users"
	"management/infra/base"
	"strconv"
)

func loginMeddle(ctx iris.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		r := base.Res{
			Code:    base.ResError,
			Message: "用户未登录",
		}
		ctx.JSON(&r)
	}
	user, err := common.ParseToken(token)
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
	token := ctx.GetHeader("Authorization")
	if token == "" {
		r := base.Res{
			Code:    base.ResError,
			Message: "用户未登录",
		}
		ctx.JSON(&r)
	}
	user, err := common.ParseToken(token)
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


func manamgeMeddle(ctx iris.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		r := base.Res{
			Code:    base.ResError,
			Message: "用户未登录",
		}
		ctx.JSON(&r)
	}
	user, err := common.ParseToken(token)
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
	res, err := users.GetUserService().GetUserByPhone(user.Phone, user.UserType)
	if err != nil || res == nil || res.SuperState != 2 {
		r := base.Res{
			Code:    base.ResError,
			Message: "无权限",
		}
		ctx.JSON(&r)
		return
	}
	ctx.Request().Header.Set("phone", user.Phone)
	ctx.Request().Header.Set("user_type", strconv.Itoa(user.UserType))
	ctx.Request().Header.Set("user_id", strconv.FormatInt(user.Id,10))
	ctx.Next()
}

func refreshToken(ctx iris.Context) string {
	phone := ctx.GetHeader("phone")
	userId, _ := strconv.ParseInt(ctx.GetHeader("user_id"), 10, 64)
	userType, _ := strconv.Atoi(ctx.GetHeader("user_type"))
	user := users.User{
		Phone:    phone,
		Id:       userId,
		UserType: userType,
	}
	token, _ := common.GenerateToken(user)
	return token
}