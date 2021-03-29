package orders

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"management/core/common"
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

func (o orderService) GetOrdersByCond(cond Order) (orders *[]Order, err error) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := OrderDao{runner: runner}
		if cond.Page == 0 {
			cond.Page = 1
		}
		if cond.PageSize < 1 || cond.PageSize > 10 {
			cond.PageSize = 10
		}
		log.Error(cond)
		orders = dao.GetByCond(cond)
		if orders == nil {
			return nil
		}
		stageDao := OrderStageDao{runner: runner}
		for _, order := range *orders {
			orderStage := stageDao.GetByOrderId(order.Id)
			order.OrderStage = orderStage
		}
		return nil
	})
	return
}

func (o orderService) GetOrdersByUser(userId int64, userType int) (finishOrders, doingOrders OrderSlice, err error) {
	var orders *[]Order
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := OrderDao{runner: runner}
		if userType == 1 {
			orders = dao.GetByUserId(userId)

		} else {
			orders = dao.GetByEmployeeId(userId)

		}
		if orders == nil {
			return nil
		}
		stageDao := OrderStageDao{runner: runner}
		for _, order := range *orders {
			orderStage := stageDao.GetByOrderId(order.Id)
			order.OrderStage = orderStage
			if order.Stage <= 7 || order.Stage >= 5 {
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

// 创建订单
func (o orderService) Create(order Order, user users.User) (*Order, error) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		// 1. 查找用户
		userDao := users.UserDao{runner}
		user := userDao.GetOne(user.Phone, user.UserType)
		if user == nil {
			err := errors.New("查找用户失败")
			return err
		}

		// 2. 创建订单
		dao := OrderDao{runner: runner}
		order.HouseholdId = user.Id
		order.HouseholdName = user.Name
		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()
		order.Stage = 1
		orderId, err := dao.Insert(&order)
		if err != nil {
			log.Error(err)
			err := errors.New("创建订单失败")
			return err
		}

		// 3.创建订单阶段
		orderStage := &OrderStage{
			Stage:     order.Stage,
			OrderId:   orderId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Note:      order.Note,
		}
		stageDao := OrderStageDao{runner}
		_, err = stageDao.Insert(orderStage)
		if err != nil {
			log.Error(err)
			err := errors.New("创建订单失败")
			return err
		}
		order.OrderStage = &[]OrderStage{*orderStage}
		// 4.分配/抢单
		scoreDao := users.EmployeeScoreDao{Runner: runner}
		if order.Emergency == 1 { // 非紧急，直接分配
			employee := scoreDao.GetNonUrgentEmployee()
			if employee == nil {
				return errors.New("分配员工失败")
			}

			order.EmployeeId = employee.EmployeeId
			_, err = dao.Update(&order)
			if err != nil {
				log.Error(err)
				err := errors.New("分配员工失败")
				return err
			}

			employee.DoingOrder += 1
			employee.DifficultScore += order.Level
			employee.NonUrgentScore, employee.UrgentScore = common.AllocationAlgorithm(*employee)
			_, err = scoreDao.Update(employee)
			if err != nil {
				log.Error(err)
				err := errors.New("分配员工失败")
				return err
			}

			// 邮件通知
			user = userDao.GetOneById(employee.EmployeeId)
			if user == nil {
				err := errors.New("查找员工失败")
				return err
			}
			common.SendMail(user.Email, "内容", "内容")
		} else {
			employees := scoreDao.GetUrgentEmployee()
			if employees == nil {
				log.Error(err)
				err := errors.New("分配员工失败")
				return err
			}
			for _, employee := range *employees {
				user = userDao.GetOneById(employee.EmployeeId)
				if user == nil {
					err := errors.New("查找员工失败")
					return err
				}
				// TODO 员工和订单要关联起来
				common.SendMail(user.Email, "内容", "内容")

			}
		}
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
		if order.Stage < 5 && (stage.Stage >= 5) {
			scoreDao := users.EmployeeScoreDao{Runner: runner}
			employee := scoreDao.GetByEmployeeId(order.EmployeeId)
			if employee == nil {
				return errors.New("找不到该订单对应的员工")
			}

			employee.DoingOrder -= 1
			_, err := scoreDao.Update(employee)
			if err != nil {
				log.Error(err)
				err := errors.New("更新员工订单数失败")
				return err
			}
		}
		order.Stage = stage.Stage
		order.UpdatedAt = time.Now()
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
