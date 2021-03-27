package users

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type UserDao struct {
	Runner *dbx.TxRunner
}

// 通过手机号和类型获取用户
func (dao *UserDao) GetOne(phone string, userType int) *User {
	form := &User{}

	ok, err := dao.Runner.Get(form, "select * from user where phone=? and user_type=?", phone, userType)
	if err != nil || !ok {
		logrus.Error(err)
		return nil
	}

	return form
}

// 通过手机号和类型获取用户
func (dao *UserDao) GetOneById(userId int64) *User {
	form := &User{}
	logrus.Error("userId",userId)

	ok, err := dao.Runner.Get(form, "select * from user where id=?", userId)
	if err != nil || !ok {
		logrus.Error(err)
		return nil
	}

	return form
}

func (dao *UserDao) Insert(form *User) (int64, error) {
	rs, err := dao.Runner.Insert(form)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

func (dao *UserDao) Update(user *User) (int64, error) {
	rs, err := dao.Runner.Update(user)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}
