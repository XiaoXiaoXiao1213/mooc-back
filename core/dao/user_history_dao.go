package dao

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"management/core/domain"
)

type UserHistoryDao struct {
	DB *mgo.Database
}

func (dao *UserHistoryDao) Insert(form *domain.UserVideoHistory) error {
	err := dao.DB.C("user_history").Insert(form)
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (dao *UserHistoryDao) GetByUserId(userId int) (*[]domain.UserVideoHistory, error) {
	var userHistory = new([]domain.UserVideoHistory)
	err := dao.DB.C("user_history").Find(bson.M{"user_id": userId}).Limit(50).Sort("-time").All(userHistory)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return userHistory, nil
}
