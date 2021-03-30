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
	groupRouter.Post("/create", loginMeddle, a.createOrder)
	groupRouter.Get("/{id}", loginMeddle, a.getOrderById)
	groupRouter.Post("/{id}", loginMeddle, a.editOrder)
	groupRouter.Post("/evaluation/{order_id}", loginMeddle, a.evaluationOrder)
	groupRouter.Put("/{order_id}", employeeMeddle, a.takeOrder)

}

// 创建订单
func (a *OrderApi) createOrder(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	order := orders.Order{}
	err := ctx.ReadJSON(&order)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		return
	}

	//获取请求参数
	phone := ctx.GetHeader("phone")
	userType, _ := strconv.Atoi(ctx.GetHeader("user_type"))
	user := users.User{Phone: phone, UserType: userType}
	res, err := a.service.Create(order, user)

	ctx.ResponseWriter().Header().Set("token", refreshToken(ctx))
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}

	r.Data = map[string]interface{}{
		"order": res,
	}
	ctx.JSON(&r)
}

// TODO 抢单
func (a *OrderApi) takeOrder(ctx iris.Context) {
	//r := base.Res{
	//	Code: base.ResCodeOk,
	//}
	//orderId, err := strconv.ParseInt(ctx.Params().Get("order_id"), 10, 64)
	//if err != nil {
	//	r.Code = base.ResError
	//	r.Message = "字段或字段值格式错误"
	//	ctx.JSON(&r)
	//	logrus.Error(err)
	//	return
	//}
	//
	//order := orders.Order{}
	//err := ctx.ReadJSON(&order)
	//if err != nil {
	//	r.Code = base.ResError
	//	r.Message = "字段或字段值格式错误"
	//	ctx.JSON(&r)
	//	logrus.Error(err)
	//	return
	//}
	//
	////获取请求参数
	//phone := ctx.GetHeader("phone")
	//userType, _ := strconv.Atoi(ctx.GetHeader("user_type"))
	//user := users.User{Phone: phone, UserType: userType}
	//res, err := a.service.Create(order, user)
	//if err != nil {
	//	r.Code = base.ResError
	//	r.Message = err.Error()
	//	logrus.Error(err)
	//}
	//ctx.ResponseWriter().Header().Set("token",refreshToken(ctx))
	//
	//r.Data = map[string]interface{}{
	//	"order": res,
	//}
	//ctx.JSON(&r)
}

// 评价订单
func (a *OrderApi) evaluationOrder(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	// 获取订单
	orderId, err := strconv.ParseInt(ctx.Params().Get("order_id"), 10, 64)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		return
	}

	order, err := a.service.GetOrdersById(orderId)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "找不到该订单"
		ctx.JSON(&r)
		return
	}

	if order.EvaluationId != 0 {
		r.Code = base.ResError
		r.Message = "该订单已评价"
		ctx.JSON(&r)
		return
	}

	//获取请求参数
	evalutaion := orders.Evaluation{}
	err = ctx.ReadJSON(&evalutaion)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		return
	}

	evalutaion.EmployeeId = order.EmployeeId
	evalutaion.OrderId = orderId
	err = a.service.TakeEvaluation(evalutaion)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
	}
	ctx.ResponseWriter().Header().Set("token", refreshToken(ctx))
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
		return
	}

	order, err := a.service.GetOrdersById(id)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "找不到订单"
		ctx.JSON(&r)
		return
	}
	ctx.ResponseWriter().Header().Set("token", refreshToken(ctx))
	r.Data = map[string]interface{}{
		"order": order,
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
		return
	}

	stage := orders.OrderStage{}
	err = ctx.ReadJSON(&stage)
	if err != nil || stage.Stage == 0 || stage.OrderId == 0 {
		r.Code = base.ResError
		r.Message = "参数错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	err = a.service.EditStage(stage)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	ctx.ResponseWriter().Header().Set("token", refreshToken(ctx))
	ctx.JSON(&r)
}
