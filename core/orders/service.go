package orders

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"management/core/users"
	"sort"

	//"management/core/users"
	"management/infra/base"
	"sync"
	"time"
)

var _ OrderService = new(orderService)
var once sync.Once

func init() {
	once.Do(func() {
		IOrderService = new(orderService)
	})
}

type orderService struct {
}

func (o orderService) GetOrdersByUser(userId int64) (finishOrders, doingOrders OrderSlice, err error) {
	var orders *[]Order
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := OrderDao{runner: runner}
		orders = dao.GetByUserId(userId)
		if orders==nil{
			return nil
		}
		stageDao := OrderStageDao{runner: runner}
		for _, order := range *orders {
			orderStage := stageDao.GetByOrderId(order.Id)
			order.OrderStage = orderStage
			if order.Stage == 7 {
				finishOrders = append(finishOrders, order)
			} else {
				doingOrders = append(doingOrders, order)
			}
		}
		sort.Sort(sort.Reverse(finishOrders))
		sort.Sort(sort.Reverse(doingOrders))
		return nil
	})
	return
}

func (o orderService) Create(order Order, user users.User) (*Order, error) {

	u, err := users.GetUserService().GetUserByPhone(user.Phone, user.UserType)
	if err != nil {
		err := errors.New("查找用户失败")
		return nil, err
	}

	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := OrderDao{runner: runner}
		order.HouseholdId = u.Id
		order.HouseholdName = u.Name
		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()
		order.Stage = 1
		res, err := dao.runner.Insert(order)
		order.Id, _ = res.LastInsertId()
		if err != nil {
			err := errors.New("查找用户失败")
			log.Error(err)
			return err
		}

		orderId, _ := res.LastInsertId()
		orderStage := &OrderStage{
			Stage:     order.Stage,
			OrderId:   orderId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Note:      order.Note,
		}

		stageDao := OrderStageDao{runner: runner}
		_, err = stageDao.Insert(orderStage)
		order.OrderStage = &[]OrderStage{*orderStage}
		return err
	})
	return &order, err
}

func (o orderService) EditStage(stage OrderStage) error {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := OrderDao{runner: runner}
		order := dao.GetOneByOrderId(stage.OrderId)
		if order == nil {
			err := errors.New("找不到订单")
			return err
		}

		order.Stage = stage.Stage
		update, err := dao.Update(order)
		if update < 1 || err != nil {
			err := errors.New("订单进度更新失败")
			return err
		}

		stageDao := OrderStageDao{runner: runner}
		stage.UpdatedAt = time.Now()
		stage.CreatedAt = time.Now()
		_, err = stageDao.Insert(&stage)
		return err
	})
	return err
}

func (o orderService) GetOrdersById(orderId int64) (order *Order, err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := OrderDao{runner: runner}
		order = dao.GetOneByOrderId(orderId)
		if order == nil {
			err := errors.New("找不到订单")
			return err
		}
		return nil
	})
	return order, err
}
func (o orderService) TakeOrder(phone string, userType int, orderId int64) error {
	if userType != 2 {
		err := errors.New("您不是员工，不能接单")
		return err
	}
	userDao := users.UserDao{}
	user := userDao.GetOne(phone, 2)
	if user != nil {
		err := errors.New("找不到该员工，不能接单")
		return err
	}
	orderDao := OrderDao{}
	order := orderDao.GetOneByOrderId(orderId)
	order.EmployeeId = user.Id
	order.EmployeeName = user.Name
	order.UpdatedAt = time.Now()
	order.Stage = 2
	update, err := orderDao.Update(order)
	if update < 1 || err != nil {
		err := errors.New("接单失败")
		return err
	}
	return nil
}

func (o orderService) TakeEvaluation(evaluation Evaluation) error {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		// 评价
		evaluationDao := EvaluationDao{runner: runner}
		orderDao := OrderDao{runner: runner}
		// 查看评价过了没
		eval := evaluationDao.GetOneByOrderId(evaluation.OrderId)
		if eval != nil {
			err := errors.New("此订单已评价过")
			return err
		}
		evaluation.CreatedAt = time.Now()
		evaluation.UpdatedAt = time.Now()
		evalId, err := evaluationDao.Insert(evaluation)
		if err != nil {
			log.Error(err)
			err := errors.New("评价失败")
			return err
		}
		// 更新订单
		order := orderDao.GetOneByOrderId(evaluation.OrderId)
		order.UpdatedAt = time.Now()
		order.EvaluationId = evalId
		update, err := orderDao.Update(order)
		if update < 1 || err != nil {
			err := errors.New("订单评价更新失败")
			return err
		}
		// 更新用户分数
		userDao := users.UserDao{Runner: runner}
		employee := userDao.GetOneById(evaluation.EmployeeId)
		if employee == nil {
			err := errors.New("订单还没被接，无法评价")
			return err
		}
		employee.Score = employee.Score + evaluation.Level
		_, err = userDao.Update(employee)
		return err
	})
	return err
}

func (o orderService) GetUserByPhone(phone string, userType int) (*Order, error) {
	panic("implement me")
}
