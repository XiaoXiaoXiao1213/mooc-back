package web

import (
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"management/core/orders"
	"management/core/users"
	"management/infra"
	"management/infra/base"
	"strconv"
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
	groupRouter.Get("/order/{page}/{page_size}", manamgeMeddle, u.getOrder)
	groupRouter.Get("/user/{page}/{page_size}", manamgeMeddle, u.getUser)

}

func (u *ManageApi) getOrder(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	page, err1 := strconv.Atoi(ctx.Params().Get("page"))
	pageSize, err2 := strconv.Atoi(ctx.Params().Get("page_size"))

	if err1 != nil || err2 != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err1, err2)
		return
	}
	//获取请求参数
	order := orders.Order{}
	err := ctx.ReadJSON(&order)
	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	order.Page = page
	order.PageSize = pageSize

	orders, err := u.orderService.GetOrdersByCond(order)
	if err != nil {
		r.Code = base.ResError
		r.Message = "查询失败"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	r.Data = map[string]interface{}{
		"token":  refreshToken(ctx),
		"orders": orders,
	}
	ctx.JSON(&r)

}

func (u *ManageApi) getUser(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	page, err1 := strconv.Atoi(ctx.Params().Get("page"))
	pageSize, err2 := strconv.Atoi(ctx.Params().Get("page_size"))
	if err1 != nil || err2 != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err1, err2)
		return
	}
	//获取请求参数
	user := users.User{}
	err := ctx.ReadJSON(&user)
	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	user.Page = page
	user.PageSize = pageSize

	users, err := u.userService.GetUserByCond(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = "查询错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	r.Data = map[string]interface{}{
		"token": refreshToken(ctx),
		"users": users,
	}
	ctx.JSON(&r)

}
