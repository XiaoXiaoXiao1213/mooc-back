package users

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type UserDao struct {
	runner *dbx.TxRunner
}

// 通过手机号和类型获取用户
func (dao *UserDao) GetOne(phone string, userType int) *User {
	form := &User{}
	ok, err := dao.runner.Get(form, "select * from user where phone=? and user_type=?", phone, userType)
	if err != nil || !ok {
		logrus.Error(err)
		return nil
	}
	return form
}

func (dao *UserDao) Insert(form *User) (int64, error) {
	rs, err := dao.runner.Insert(form)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

func (dao *UserDao) Update(user *User) (int64, error) {
	rs, err := dao.runner.Update(user)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}
