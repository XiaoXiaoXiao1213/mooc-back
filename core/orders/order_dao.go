package orders

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"strconv"
)

type OrderDao struct {
	runner *dbx.TxRunner
}

func (dao *OrderDao) GetByUserId(userId int64) *[]Order {
	form := []Order{}
	err := dao.runner.Find(&form, "select * from `order` where household_id=?", userId)
	if err != nil {
		log.Error(err)
		return nil
	}
	return &form
}

func (dao *OrderDao) GetByEmployeeId(employeeId int64) *[]Order {
	form := []Order{}
	err := dao.runner.Find(&form, "select * from `order` where employee_id=?", employeeId)
	if err != nil {
		log.Error(err)
		return nil
	}
	return &form
}

func (dao *OrderDao) GetByCond(order Order) (*[]Order, int) {
	form := []Order{}
	sql := "select * from `order` where 1"
	if order.Id != 0 {
		sql += " and id=" + strconv.FormatInt(order.Id, 10)
	}
	if order.HouseholdId != 0 {
		sql += " and household_id=" + strconv.FormatInt(order.HouseholdId, 10)
	}
	if order.HouseId != 0 {
		sql += " and house_id=" + strconv.FormatInt(order.HouseId, 10)
	}
	if order.EmployeeId != 0 {
		sql += " and employee_id=" + strconv.FormatInt(order.EmployeeId, 10)
	}
	if order.Type != "" {
		sql += " and type=\"" + order.Type + "\""
	}
	if order.Emergency != 0 {
		sql += " and emergency=" + strconv.Itoa(order.Emergency)

	}
	if order.Stage != 0 {
		sql += " and stage=" + strconv.Itoa(order.Stage)
	}
	err := dao.runner.Find(&form, sql)
	if err != nil {
		log.Error(err)
		return nil, 0
	}
	count := len(form)
	err = dao.runner.Find(&form, sql+" limit ?,?", order.Page-1, order.PageSize)
	if err != nil {
		log.Error(err)
		return nil, 0
	}
	return &form, count
}
func (dao *OrderDao) GetOneByOrderId(orderId int64) *Order {
	form := &Order{}
	ok, err := dao.runner.Get(form, "select * from `order` where id=?", orderId)
	log.Error("e", form)

	if err != nil || !ok {
		return nil
	}
	return form
}

func (dao *OrderDao) Insert(form *Order) (int64, error) {
	rs, err := dao.runner.Insert(form)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

func (dao *OrderDao) Update(order *Order) (int64, error) {
	rs, err := dao.runner.Update(order)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}


func (dao *OrderDao) DeleteByOrderId(orderId int64) (int64, error) {
	rs, err := dao.runner.Exec("delete from `order` where id=? limit 1",orderId)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}