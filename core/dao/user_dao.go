package dao

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"management/core/domain"
)

type UserDao struct {
	DB *mgo.Database
}

// 通过手机号
func (dao *UserDao) GetUserByPhone(phone string) (*domain.User, error) {
	var user = new(domain.User)
	err := dao.DB.C("user").Find(bson.M{"phone": phone}).One(user)
	if err != nil {
		logrus.Error( err)
		return nil, err
	}
	return user, nil
}

// 通过用户id
func (dao *UserDao) GetUserByUserId(userId string) (*domain.User, error) {
	var user = new(domain.User)
	err := dao.DB.C("user").FindId(bson.ObjectIdHex(userId)).One(user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return user, nil
}

func (dao *UserDao) Insert(form *domain.User) error {
	all := dao.GetAll()
	form.UserInt = len(*all) + 1
	err := dao.DB.C("user").Insert(form)
	if err != nil {
		logrus.Error( err)
	}
	return err
}

func (dao *UserDao) Update(user *domain.User) error {
	err := dao.DB.C("user").Update(bson.M{"phone": user.Phone}, user)
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (dao *UserDao) GetAll() *[]domain.User {
	var users = new([]domain.User)
	err := dao.DB.C("user").Find(nil).All(users)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return users
}
