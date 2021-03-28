package web

import (
	"github.com/kataras/iris"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"management/core/orders"
	"management/core/users"
	"management/infra"
	"management/infra/base"
)

func init() {
	infra.RegisterApi(new(ManageApi))
}

type ManageApi struct {
	orderService orders.OrderService
	userService  users.UserService
}

func (u *ManageApi) Init() {
	u.orderService = orders.GetOrderService()
	u.userService = users.GetUserService()
	groupRouter := base.Iris().Party("/api/1.0/management")
	groupRouter.Get("/order", manamgeMeddle, u.getOrder)
	groupRouter.Get("/user", manamgeMeddle, u.getUser)

}

func (u *ManageApi) getOrder(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	//获取请求参数
	order := orders.Order{}
	err := ctx.ReadForm(&order)
	log.Error(order)
	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	orders, err := u.orderService.GetOrdersByCond(order)
	if err != nil {
		r.Code = base.ResError
		r.Message = "查询失败"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	ctx.ResponseWriter().Header().Set("token",refreshToken(ctx))
	r.Data = map[string]interface{}{
		"orders": orders,
	}
	ctx.JSON(&r)

}

func (u *ManageApi) getUser(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	//获取请求参数
	user := users.User{}
	err := ctx.ReadForm(&user)
	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	users, err := u.userService.GetUserByCond(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = "查询错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	ctx.ResponseWriter().Header().Set("token",refreshToken(ctx))

	r.Data = map[string]interface{}{
		"users": users,
	}
	ctx.JSON(&r)

}
