package dao

import (
	"errors"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"management/core/domain"
)

type UserClickDao struct {
	DB *mgo.Database
}

func (dao *UserClickDao) GetByUserAndVideo(videoClick domain.VideoClick) (*domain.VideoClick, error) {
	var video = new(domain.VideoClick)
	query := bson.M{
		"$and": []interface{}{
			bson.M{"user_id": videoClick.UserId},
			bson.M{"video_id": videoClick.VideoId},
		},
	}
	err := dao.DB.C("user_click").Find(query).One(video)
	if err != nil {
		log.Error(err)
		return nil, errors.New("not found")
	}
	return video, err
}

func (dao *UserClickDao) Insert(form *domain.VideoClick) error {
	err := dao.DB.C("user_click").Insert(form)
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (dao *UserClickDao) Update(video *domain.VideoClick) error {
	err := dao.DB.C("user_click").Update(bson.M{"_id": video.Id}, video)
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (dao *UserClickDao) GetAll() *[]domain.VideoClick {
	res := new([]domain.VideoClick)
	err := dao.DB.C("user_click").Find(nil).All(res)
	if err != nil {
		logrus.Error(err)
	}
	return res
}
