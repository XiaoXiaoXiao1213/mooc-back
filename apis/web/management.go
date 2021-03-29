package web

import (
	"github.com/kataras/iris"
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

// 条件查询订单
func (u *ManageApi) getOrder(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	//获取请求参数
	orderCond := orders.Order{}
	err := ctx.ReadForm(&orderCond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		return
	}

	orders, err := u.orderService.GetOrdersByCond(orderCond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "查询失败"
		ctx.JSON(&r)
		return
	}

	ctx.ResponseWriter().Header().Set("token",refreshToken(ctx))
	r.Data = map[string]interface{}{
		"orders": orders,
	}
	ctx.JSON(&r)
}

// 条件查询住户/员工
func (u *ManageApi) getUser(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	//获取请求参数
	userCond := users.User{}
	err := ctx.ReadForm(&userCond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		return
	}

	users, err := u.userService.GetUserByCond(userCond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "查询错误"
		ctx.JSON(&r)
		return
	}

	ctx.ResponseWriter().Header().Set("token",refreshToken(ctx))
	r.Data = map[string]interface{}{
		"users": users,
	}
	ctx.JSON(&r)
}
