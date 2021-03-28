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

func (dao *OrderDao) GetByCond(order Order) *[]Order {
	form := []Order{}
	sql := "select * from `order` where 1"
	if order.Id != 0 {
		sql = sql + " and id=" + strconv.FormatInt(order.Id, 10)
	}
	if order.HouseholdId != 0 {
		sql = sql + " and household_id=" + strconv.FormatInt(order.HouseholdId, 10)
	}
	if order.HouseId != 0 {
		sql = sql + " and house_id=" + strconv.FormatInt(order.HouseId, 10)
	}
	if order.EmployeeId != 0 {
		sql = sql + " and employee_id=" + strconv.FormatInt(order.EmployeeId, 10)
	}
	if order.Type != "" {
		sql = sql + " and type=\"" + order.Type + "\""
	}
	if order.Emergency != 0 {
		sql = sql + " and emergency=" + strconv.Itoa(order.Emergency)

	}
	if order.Stage != 0 {
		sql = sql + " and stage=" + strconv.Itoa(order.Stage)

	}
	log.Error(sql)
	err := dao.runner.Find(&form, sql+" limit ?,?", order.Page-1, order.PageSize)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Error(len(form))

	return &form
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
