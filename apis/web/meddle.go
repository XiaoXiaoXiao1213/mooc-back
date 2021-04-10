package web

import (
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
	"management/core/common"
	"management/core/domain"
	"management/infra/base"
)

//
func Cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
	ctx.Next()
}

func loginMeddle(ctx iris.Context) {
	res, user := parseTokenByRequest(ctx)
	if res.Code != base.ResCodeOk {
		ctx.JSON(&res)
		return
	}
	ctx.Request().Header.Set("phone", (*user)["phone"])
	ctx.Request().Header.Set("user_id", (*user)["user_id"])
	ctx.Next()
}

func refreshToken(ctx iris.Context) string {
	phone := ctx.GetHeader("phone")
	userId := ctx.GetHeader("user_id")
	user := domain.User{
		Phone: phone,
		Id:    bson.ObjectId(userId),
	}
	token, _ := common.GenerateToken(user)
	return token
}

func parseTokenByRequest(ctx iris.Context) (base.Res, *map[string]string) {
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
