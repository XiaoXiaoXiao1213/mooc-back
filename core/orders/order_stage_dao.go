package orders

import (
	"github.com/tietang/dbx"
)


type OrderStageDao struct {
	runner *dbx.TxRunner
}


func (dao *OrderStageDao) GetByOrderId(orderId int64) *[]OrderStage {
	form := &[]OrderStage{}
	ok, err := dao.runner.Get(form, "select * from order_stage where order_id=?", orderId)
	if err != nil || !ok {
		return nil
	}
	return form
}




func (dao *OrderStageDao) Insert(form *OrderStage) (int64, error) {
	rs, err := dao.runner.Insert(form)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

