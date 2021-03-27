package orders

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type OrderDao struct {
	runner *dbx.TxRunner
}




func (dao *OrderDao) GetByUserId(userId int64) *[]Order {
	form := &[]Order{}
	ok, err := dao.runner.Get(form, "select * from order where household_id=?", userId)
	if err != nil || !ok {
		return nil
	}
	return form
}
func (dao *OrderDao) GetOneByOrderId(orderId int64) *Order {
	form := &Order{}
	ok, err := dao.runner.Get(form, "select * from `order` where id=?", orderId)
	log.Error("e",form)

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
