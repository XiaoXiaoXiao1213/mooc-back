package web

import (
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	houses "management/core/house"
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
	houseService houses.HouseService
}

func (u *ManageApi) Init() {
	u.orderService = orders.GetOrderService()
	u.userService = users.GetUserService()
	u.houseService = houses.GetHouseService()
	groupRouter := base.Iris().Party("/api/1.0/management")
	groupRouter.Get("/order", manamgeMeddle, u.getOrder)
	groupRouter.Get("/user", manamgeMeddle, u.getUser)
	groupRouter.Get("/house", manamgeMeddle, u.getHouse)
	groupRouter.Delete("/house/delete/{house_id}", manamgeMeddle, u.deleteHouse)
	groupRouter.Delete("/order/delete/{order_id}", manamgeMeddle, u.deleteOrder)
	groupRouter.Delete("/user/delete/{user_id}", manamgeMeddle, u.deleteUser)

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

	orders, count, err := u.orderService.GetOrdersByCond(orderCond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "查询失败"
		ctx.JSON(&r)
		return
	}

	ctx.ResponseWriter().Header().Set("token", refreshToken(ctx))
	r.Data = map[string]interface{}{
		"orders": orders,
		"total":  count,
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

	users, total, err := u.userService.GetUserByCond(userCond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "查询错误"
		ctx.JSON(&r)
		return
	}

	ctx.ResponseWriter().Header().Set("token", refreshToken(ctx))
	r.Data = map[string]interface{}{
		"users": users,
		"total": total,
	}
	ctx.JSON(&r)
}

// 条件查询房子
func (u *ManageApi) getHouse(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	//获取请求参数
	houseCond := houses.House{}
	err := ctx.ReadForm(&houseCond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		return
	}

	houses, total, err := u.houseService.GetHousesByCond(houseCond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "查询错误"
		ctx.JSON(&r)
		return
	}

	ctx.ResponseWriter().Header().Set("token", refreshToken(ctx))
	r.Data = map[string]interface{}{
		"houses": houses,
		"total":  total,
	}
	ctx.JSON(&r)
}

func (u *ManageApi) deleteHouse(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	houseId := ctx.Params().Get("house_id")
	err := u.houseService.DeleteHousesByHouseId(houseId)
	if err != nil {
		r.Code = base.ResError
		r.Message = "删除住宅失败"
	}
	ctx.ResponseWriter().Header().Set("token", refreshToken(ctx))
	ctx.JSON(&r)
}

func (u *ManageApi) deleteOrder(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	ctx.ResponseWriter().Header().Set("token", refreshToken(ctx))
	orderId, err := strconv.ParseInt(ctx.Params().Get("order_id"), 10, 64)
	if err != nil {
		r.Code = base.ResError
		r.Message = "订单id格式错误"
		ctx.JSON(&r)
	}
	err = u.orderService.DeleteOrdersById(orderId)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
	}
	ctx.JSON(&r)
}

func (u *ManageApi) deleteUser(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	ctx.ResponseWriter().Header().Set("token", refreshToken(ctx))
	userId, err := strconv.ParseInt(ctx.Params().Get("user_id"), 10, 64)
	if err != nil {
		r.Code = base.ResError
		r.Message = "用户id格式错误"
		ctx.JSON(&r)
	}
	err = u.userService.DeleteUserById(userId)
	if err != nil {
		r.Code = base.ResError
		r.Message = "删除用户失败"
	}
	ctx.JSON(&r)
}
