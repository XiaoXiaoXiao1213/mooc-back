package orders

import (
	"management/core/users"
	"management/infra/base"
)

var IOrderService OrderService

//用于对外暴露账户应用服务，唯一的暴露点
func GetOrderService() OrderService {
	base.Check(IOrderService)
	return IOrderService
}

type OrderService interface {
	Create(order Order, user users.User) (*Order, error)
	EditStage(stage OrderStage) error
	TakeOrder(phone string, userType int, orderId int64) error
	TakeEvaluation(evaluation Evaluation) error
	GetOrdersByUser(userId int64,userType int) (finishOrders, doingOrders OrderSlice, err error)
	GetOrdersById(orderId int64) (order *Order, err error)
	GetOrdersByCond(cond Order) (orders *[]Order, count int, err error)
}
