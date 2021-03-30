package users

import (
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"strconv"
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
func (dao *UserDao) GetByCond(user User) (*[]User, int) {
	form := []User{}
	sql := "select * from `user` where 1"
	if user.Id != 0 {
		sql += " and id=" + strconv.FormatInt(user.Id, 10)
	}
	if user.SuperState != 0 {
		sql += " and super_state=" + strconv.Itoa(user.SuperState)
	}
	if user.Sex != 0 {
		sql += " and sex=" + strconv.Itoa(user.Sex)
	}
	if user.Id_code != "" {
		sql += " and id_code like\"%" + user.Id_code + "%\""
	}
	if user.Name != "" {
		sql += " and name like\"%" + user.Name + "%\""
	}
	if user.UserType != 0 {
		sql += " and user_type=" + strconv.Itoa(user.UserType)
	}
	err := dao.Runner.Find(&form, sql)
	if err != nil {
		log.Error(err)
		return nil, 0
	}
	count := len(form)
	form = []User{}
	err = dao.Runner.Find(&form, sql+" limit ?,?", user.Page-1, user.PageSize)
	if err != nil {
		log.Error(err)
		return nil, 0
	}
	return &form, count
}

// 通过手机号和类型获取用户
func (dao *UserDao) GetOneById(userId int64) *User {
	form := &User{}
	logrus.Error("userId", userId)

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

func (dao *UserDao) DeleteByUserId(userId int64) (int64, error) {
	rs, err := dao.Runner.Exec("delete from user where id=? limit 1", userId)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}
