package users

import (
	"github.com/prometheus/common/log"
	"github.com/tietang/dbx"
	"strconv"
)

type EmployeeScoreDao struct {
	Runner *dbx.TxRunner
}

func (dao *EmployeeScoreDao) GetByCond(user User) *[]User {
	form := []User{}
	sql := "select * from `user` where 1"
	if user.Id != 0 {
		sql = sql + " and id=" + strconv.FormatInt(user.Id, 10)
	}
	if user.SuperState != 0 {
		sql = sql + " and super_state=" + strconv.Itoa(user.SuperState)
	}
	if user.Sex != 0 {
		sql = sql + " and sex=" + strconv.Itoa(user.Sex)
	}
	if user.Id_code != "" {
		sql = sql + " and id_code=\"%" + user.Id_code + "%\""
	}
	if user.Name != "" {
		sql = sql + " and name=\"%" + user.Name + "%\""
	}
	err := dao.Runner.Find(&form, sql+" limit ?,?", user.Page-1, user.PageSize)
	if err != nil {
		log.Error(err)
		return nil
	}
	return &form
}

// 通过手机号和类型获取用户
func (dao *EmployeeScoreDao) GetUrgentEmployee(orderType string) *[]EmployeeScore {
	form := []EmployeeScore{}
	sql := "select * from `employee_score` where state=2 and doing_order<2 and skills like \"%" + orderType + "%\" order by urgent_score desc limit 3"
	log.Error(sql)
	err := dao.Runner.Find(&form, sql)
	if err != nil || len(form) == 0 {
		log.Error(err)
		sql = "select * from `employee_score` where state=2 and skills like \"%" + orderType + "%\" order by doing_order desc limit 3"
		err = dao.Runner.Find(&form, sql)
		if err != nil {
			log.Error(err)
			return nil
		}
	}
	return &form
}
func (dao *EmployeeScoreDao) GetNonUrgentEmployee(orderType string) *EmployeeScore {
	form := &[]EmployeeScore{}
	sql := "select * from employee_score  where state=2 and doing_order<5 and skills like \"%" + orderType + "%\" order by non_urgent_score desc limit 1"
	err := dao.Runner.Find(form, sql)
	log.Error(form)
	if err != nil || len(*form) == 0 {
		sql = "select * from `employee_score` where state=2 and skills like \"%" + orderType + "%\" order by non_urgent_score desc limit 1"
		err = dao.Runner.Find(&form, sql)
		if err != nil || len(*form) == 0 {
			log.Error(err)
			return nil
		}
	}
	return &(*form)[0]
}

func (dao *EmployeeScoreDao) GetByEmployeeId(employeeId int64) *EmployeeScore {
	form := EmployeeScore{}
	sql := "select * from `employee_score` where employee_id=?"
	ok, err := dao.Runner.Get(&form, sql, employeeId)
	if err != nil || !ok {
		log.Error(err)
		return nil
	}
	return &form
}

func (dao *EmployeeScoreDao) Insert(form *EmployeeScore) (int64, error) {
	rs, err := dao.Runner.Insert(form)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

func (dao *EmployeeScoreDao) Update(user *EmployeeScore) (int64, error) {
	rs, err := dao.Runner.Update(user)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}
