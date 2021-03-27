package orders

import (
	"github.com/tietang/dbx"
)

type EvaluationDao struct {
	runner *dbx.TxRunner
}


func (dao *EvaluationDao) GetOneByOrderId(orderId int64) *Evaluation {
	form := &Evaluation{}
	ok, err := dao.runner.Get(form, "select * from evaluation where order_id=?", orderId)
	if err != nil || !ok {
		return nil
	}
	return form
}

func (dao *EvaluationDao) Insert(evaluation Evaluation) (int64, error) {
	rs, err := dao.runner.Insert(evaluation)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

