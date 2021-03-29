package web

import (
	"github.com/kataras/iris"
	"management/core/common"
	"management/core/users"
	"management/infra/base"
	"strconv"
)

func loginMeddle(ctx iris.Context) {
	res, user := parseTokenByRequest(ctx)
	if res.Code != base.ResCodeOk {
		ctx.JSON(&res)
		return
	}

	ctx.Request().Header.Set("phone", user.Phone)
	ctx.Request().Header.Set("user_type", strconv.Itoa(user.UserType))
	ctx.Request().Header.Set("user_id", strconv.FormatInt(user.Id, 10))
	ctx.Next()
}

func employeeMeddle(ctx iris.Context) {
	res, user := parseTokenByRequest(ctx)
	if res.Code != base.ResCodeOk {
		ctx.JSON(&res)
		return
	}

	if user.UserType != 2 {
		res.Code = base.ResError
		res.Message = "该用户不是员工"
		ctx.JSON(&res)
		return
	}

	ctx.Request().Header.Set("phone", user.Phone)
	ctx.Request().Header.Set("user_type", strconv.Itoa(user.UserType))
	ctx.Request().Header.Set("user_id", strconv.FormatInt(user.Id, 10))
	ctx.Next()
}

func manamgeMeddle(ctx iris.Context) {
	res, user := parseTokenByRequest(ctx)
	if res.Code != base.ResCodeOk {
		ctx.JSON(&res)
		return
	}

	if user.UserType != 2 {
		res.Code = base.ResError
		res.Message = "该用户不是员工"
		ctx.JSON(&res)
		return
	}

	user, err := users.GetUserService().GetUserByPhone(user.Phone, user.UserType)
	if err != nil || user == nil || user.SuperState != 2 {
		res.Code = base.ResError
		res.Message = "无权限"
		ctx.JSON(&res)
		return
	}

	ctx.Request().Header.Set("phone", user.Phone)
	ctx.Request().Header.Set("user_type", strconv.Itoa(user.UserType))
	ctx.Request().Header.Set("user_id", strconv.FormatInt(user.Id, 10))
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

func parseTokenByRequest(ctx iris.Context) (base.Res, *users.User) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	token := ctx.GetHeader("Authorization")
	if token == "" {
		r.Code = base.ResError
		r.Message = "用户未登录"
		return r, nil
	}

	user, err := common.ParseToken(token)
	if err != nil {
		r.Code = base.ResError
		r.Message = "token失效，请重新登陆"
		return r, nil
	}
	return r, user
}
