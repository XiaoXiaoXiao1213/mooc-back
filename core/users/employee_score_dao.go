package users

import (
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"strconv"
)

type EmployeeScoreDao struct {
	Runner *dbx.TxRunner
}

// 通过手机号和类型获取用户
func (dao *EmployeeScoreDao) Get(phone string, userType int) *User {
	form := &User{}

	ok, err := dao.Runner.Get(form, "select * from user where phone=? and user_type=?", phone, userType)
	if err != nil || !ok {
		logrus.Error(err)
		return nil
	}

	return form
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
func (dao *EmployeeScoreDao) GetUrgentEmployee() *[]EmployeeScore {
	form := []EmployeeScore{}
	sql := "select * from `employee_score` where state=2 and doing_order<2 order by urgent_score desc limit 3"
	err := dao.Runner.Find(&form, sql)
	if err != nil || len(form) == 0 {
		log.Error(err)
		sql = "select * from `employee_score` where state=2 order by doing_order desc limit 3"
		err = dao.Runner.Find(&form, sql)
		if err != nil {
			log.Error(err)
			return nil
		}
	}
	return &form
}
func (dao *EmployeeScoreDao) GetNonUrgentEmployee() *EmployeeScore {
	form := EmployeeScore{}
	sql := "select * from `employee_score` where state=2 and doing_order<5 order by non_urgent_score desc limit 1"
	ok, err := dao.Runner.Get(&form, sql)
	if err != nil || !ok {
		log.Error(err)
		sql = "select * from `employee_score` where state=2 order by non_urgent_score desc limit 1"
		ok, err = dao.Runner.Get(&form, sql)
		if err != nil || !ok {
			log.Error(err)
			return nil
		}
	}
	return &form
}

func (dao *EmployeeScoreDao) GetByEmployeeId(employeeId int64) *EmployeeScore {
	form := EmployeeScore{}
	sql := "select * from `employee_score` where employee_id=?"
	ok, err := dao.Runner.Get(&form, sql,employeeId)
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
