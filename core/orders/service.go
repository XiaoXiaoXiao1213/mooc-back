package orders

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"management/core/common"
	"management/core/users"
	"sort"
	"strconv"
	"strings"

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

func (o orderService) DeleteOrdersById(orderId int64) (err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := OrderDao{runner: runner}
		// 获取订单
		order := dao.GetOneByOrderId(orderId)
		if order == nil {
			return errors.New("订单查找失败")
		}
		// 如果订单还在进行中，那么删除用户正在做的订单数
		if order.Stage < 5 && order.EmployeeId != 0 {
			scoreDao := users.EmployeeScoreDao{Runner: runner}
			employee := scoreDao.GetByEmployeeId(order.EmployeeId)
			if employee == nil {
				return errors.New("用户查找失败")

			}
			employee.DoingOrder -= 1
			scoreDao.Update(employee)
		}
		_, err := dao.DeleteByOrderId(orderId)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	})
	return err
}

func (o orderService) GetOrdersByCond(cond Order) (orders *[]Order, count int, err error) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := OrderDao{runner: runner}
		if cond.Page == 0 {
			cond.Page = 1
		}
		if cond.PageSize < 1 || cond.PageSize > 10 {
			cond.PageSize = 10
		}
		log.Error(cond)
		orders, count = dao.GetByCond(cond)
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
		order.Id = orderId
		if err != nil {
			log.Error(err)
			err := errors.New("创建订单失败")
			return err
		}

		types := strings.Split(order.Type, "-")
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
			employee := scoreDao.GetNonUrgentEmployee(types[0])
			if employee == nil {
				return errors.New("分配员工失败")
			}

			order.EmployeeId = employee.EmployeeId
			order.UpdatedAt = time.Now()
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
			employees := scoreDao.GetUrgentEmployee(types[0])
			if employees == nil {
				err := errors.New("分配员工失败")
				log.Error(err)
				return err
			}
			var employeeIds []string
			for _, employee := range *employees {
				user = userDao.GetOneById(employee.EmployeeId)
				if user == nil {
					err := errors.New("查找员工失败")
					log.Error(err)
					return err
				}
				orderIds, err := common.Get(strconv.FormatInt(employee.EmployeeId, 10))
				if err != nil {
					orderIds = strconv.FormatInt(orderId, 10)
				} else {
					orderIds += "|" + strconv.FormatInt(orderId, 10)
				}
				err = common.Set(strconv.FormatInt(employee.EmployeeId, 10), orderIds)
				if err != nil {
					log.Error(err)
					return err
				}
				employeeIds = append(employeeIds, strconv.FormatInt(employee.EmployeeId, 10))
				common.SendMail(user.Email, "内容", "内容")
			}
			employeeIdString := strings.Join(employeeIds, "|")
			err := common.Set(strconv.FormatInt(order.Id, 10), employeeIdString)
			if err != nil {
				log.Error(err)
				return err
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
		if order.Stage >= 5 && stage.Stage < 5 {
			err := errors.New("该订单已完成，不能修改为进行中")
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
	err := base.Tx(func(runner *dbx.TxRunner) error {
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
			err := errors.New("更新订单失败")
			return err
		}

		scoreDao := users.EmployeeScoreDao{Runner: runner}
		employeeScore := scoreDao.GetByEmployeeId(user.Id)
		if employeeScore == nil {
			err := errors.New("查找员工分数表失败")
			return err
		}
		employeeScore.DoingOrder += 1
		update, err = scoreDao.Update(employeeScore)
		if update < 1 || err != nil {
			err := errors.New("更新员工分数表失败")
			return err
		}

		employeeIds, err := common.Get(strconv.FormatInt(orderId, 10))
		if err != nil {
			err := errors.New("查找紧急订单人员失败")
			return err
		}
		common.Delete(strconv.FormatInt(orderId, 10))
		employeeIdArr := strings.Split(employeeIds, "|")
		for _, employeeId := range employeeIdArr {
			var newOrders []string
			orders, _ := common.Get(employeeId)
			ordersArr := strings.Split(orders, "|")
			for _, employeeOrder := range ordersArr {
				if employeeOrder == strconv.FormatInt(orderId, 10) {
					continue
				}
				newOrders = append(newOrders, employeeOrder)
			}
			err = common.Set(employeeId, strings.Join(newOrders, "|"))
			if err != nil {
				return err
			}

		}
		return nil
	})
	return err
}

func (o orderService) TakeEvaluation(evaluation Evaluation) error {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		evaluationDao := EvaluationDao{runner: runner}
		orderDao := OrderDao{runner: runner}
		// 查看评价过了没
		eval := evaluationDao.GetOneByOrderId(evaluation.OrderId)
		if eval != nil {
			err := errors.New("此订单已评价过")
			return err
		}
		// 查看订单完成了没
		order := orderDao.GetOneByOrderId(evaluation.OrderId)
		if order.Stage < 5 {
			err := errors.New("订单未完成，无法评价")
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

		order.Stage = 7
		order.UpdatedAt = time.Now()
		order.EvaluationId = evalId
		update, err := orderDao.Update(order)
		if update < 1 || err != nil {
			err := errors.New("订单评价更新失败")
			return err
		}

		// 更新用户分数
		userDao := users.UserDao{Runner: runner}
		scoreDao := users.EmployeeScoreDao{Runner: runner}
		employee := userDao.GetOneById(evaluation.EmployeeId)
		employeeScore := scoreDao.GetByEmployeeId(evaluation.EmployeeId)

		if employee == nil || employeeScore == nil {
			err := errors.New("找不到员工，无法评价")
			return err
		}
		employee.Score += evaluation.Level
		_, err = userDao.Update(employee)
		if err != nil {
			log.Error(err)
			return err
		}
		employeeScore.OrderScore += evaluation.Level
		employeeScore.OrderCount += 1
		employeeScore.UrgentScore, employeeScore.UrgentScore = common.AllocationAlgorithm(*employeeScore)
		_, err = scoreDao.Update(employeeScore)
		return err
	})
	return err
}
