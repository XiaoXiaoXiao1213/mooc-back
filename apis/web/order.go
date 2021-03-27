package web

import (
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"management/core"
	"management/core/orders"
	"management/core/users"
	"management/infra"
	"management/infra/base"
	"strconv"
)

func init() {
	infra.RegisterApi(new(OrderApi))
}

type OrderApi struct {
	service orders.OrderService
}

func (a *OrderApi) Init() {
	a.service = orders.GetOrderService()
	groupRouter := base.Iris().Party("/api/1.0/order")
	groupRouter.Post("/create", loginMeddle, employeeMeddle, a.create)
	groupRouter.Get("/{id}", loginMeddle, a.orderId)
	groupRouter.Post("/{id}", loginMeddle, a.editOrder)
}

// 创建订单
func (a *OrderApi) create(ctx iris.Context) {
	//获取请求参数
	phone := ctx.GetHeader("phone")
	user := users.User{Phone: phone}
	order := orders.Order{}
	err := ctx.ReadJSON(&order)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	// TODO 参数检测
	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	//创建用户
	err = a.service.Create(order, user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		logrus.Error(err)
		ctx.JSON(&r)
		return
	}

	ctx.JSON(&r)

}

func (a *OrderApi) orderId(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	//获取请求参数
	orderId := ctx.Params().Get("id")
	id, err := strconv.ParseInt(orderId, 10, 64)
	if id <= 0 || err != nil {
		r.Code = base.ResError
		r.Message = "参数错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	order, err := a.service.GetOrdersById(id)
	if err != nil {
		r.Code = base.ResError
		r.Message = "找不到订单"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	phone := ctx.GetHeader("phone")
	userId, _ := strconv.ParseInt(ctx.GetHeader("user_id"), 10, 64)
	userType, _ := strconv.Atoi(ctx.GetHeader("user_type"))
	user := users.User{
		Phone:    phone,
		Id:       userId,
		UserType: userType,
	}
	token, _ := core.GenerateToken(user)
	r.Data = map[string]interface{}{
		"order": order,
		"token": token,
	}
	ctx.JSON(&r)

}

// 创建订单
func (a *OrderApi) editOrder(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	//获取请求参数
	orderId := ctx.Params().Get("id")
	id, err := strconv.ParseInt(orderId, 10, 64)
	if id <= 0 || err != nil {
		r.Code = base.ResError
		r.Message = "参数错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	order, err := a.service.GetOrdersById(id)
	if err != nil {
		r.Code = base.ResError
		r.Message = "找不到订单"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	stage := orders.OrderStage{}
	err = ctx.ReadJSON(&stage)
	if err != nil {
		r.Code = base.ResError
		r.Message = "参数错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	phone := ctx.GetHeader("phone")
	userId, _ := strconv.ParseInt(ctx.GetHeader("user_id"), 10, 64)
	userType, _ := strconv.Atoi(ctx.GetHeader("user_type"))
	user := users.User{
		Phone:    phone,
		Id:       userId,
		UserType: userType,
	}
	token, _ := core.GenerateToken(user)
	r.Data = map[string]interface{}{
		"order": order,
		"token": token,
	}
	ctx.JSON(&r)

}
