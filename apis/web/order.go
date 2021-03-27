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
	infra.RegisterApi(new(OrderApi))
}

type OrderApi struct {
	service orders.OrderService
}

func (a *OrderApi) Init() {
	a.service = orders.GetOrderService()
	groupRouter := base.Iris().Party("/api/1.0/order")
	groupRouter.Post("/create", loginMeddle, employeeMeddle, a.create)
	groupRouter.Get("/{id}", loginMeddle, a.getOrderById)
	groupRouter.Post("/{id}", loginMeddle, a.editOrder)
	groupRouter.Post("/evaluation/{order_id}", loginMeddle, a.evaluationOrder)
}

// 创建订单
func (a *OrderApi) create(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	order := orders.Order{}
	err := ctx.ReadJSON(&order)
	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	//获取请求参数
	phone := ctx.GetHeader("phone")
	user := users.User{Phone: phone}
	err = a.service.Create(order, user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		logrus.Error(err)
	}
	ctx.JSON(&r)
}

func (a *OrderApi) evaluationOrder(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	//获取请求参数
	evalutaion := orders.Evaluation{}
	err := ctx.ReadJSON(&evalutaion)
	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	err = a.service.TakeEvaluation(evalutaion)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		logrus.Error(err)
	}
	ctx.JSON(&r)
}

func (a *OrderApi) getOrderById(ctx iris.Context) {
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
	r.Data = map[string]interface{}{
		"order": order,
		"token": refreshToken(ctx),
	}
	ctx.JSON(&r)
}

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

	stage := orders.OrderStage{}
	err = ctx.ReadJSON(&stage)
	if err != nil {
		r.Code = base.ResError
		r.Message = "参数错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	err = a.service.EditStage(stage)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	r.Data = map[string]interface{}{
		"token": refreshToken(ctx),
	}
	ctx.JSON(&r)

}
