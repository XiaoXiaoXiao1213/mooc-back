package dao

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"management/core/domain"
)

type VideoDao struct {
	DB *mgo.Database
}

func (dao *VideoDao) GetVideosByTypeAndHot(video domain.Video) (*[]domain.Video, int, error) {
	var videos = new([]domain.Video)
	var err error
	var total int
	if video.Content != "" && video.VideoType != "" {
		query := bson.M{
			"video_type": video.VideoType,
			"$or": []interface{}{
				bson.M{"name": bson.M{"$regex": video.Content, "$options": "i"}},
				bson.M{"video_type": bson.M{"$regex": video.Content, "$options": "i"}},
			},
		}
		err = dao.DB.C("video").Find(query).Sort("-click").Skip((video.Page - 1) * video.PageSize).Limit(video.PageSize).All(videos)
		total, err = dao.DB.C("video").Find(query).Count()
	} else if video.Content == "" && video.VideoType != "" {
		err = dao.DB.C("video").Find(bson.M{"video_type": video.VideoType}).Sort("-click").Skip((video.Page - 1) * video.PageSize).Limit(video.PageSize).All(videos)
		total, err = dao.DB.C("video").Find(bson.M{"video_type": video.VideoType}).Count()
	} else if video.Content != "" && video.VideoType == "" {
		query := bson.M{
			"$or": []interface{}{
				bson.M{"name": bson.M{"$regex": video.Content, "$options": "i"}},
				bson.M{"video_type": bson.M{"$regex": video.Content, "$options": "i"}},
			},
		}
		err = dao.DB.C("video").Find(query).Sort("-click").Skip((video.Page - 1) * video.PageSize).Limit(video.PageSize).All(videos)
		total, err = dao.DB.C("video").Find(query).Count()
	} else {
		err = dao.DB.C("video").Find(nil).Sort("-click").Skip((video.Page - 1) * video.PageSize).Limit(video.PageSize).All(videos)
		total, err = dao.DB.C("video").Find(nil).Count()
	}
	if err != nil {
		logrus.Error(err)
		return nil, 0, err
	}
	return videos, total, nil
}

func (dao *VideoDao) GetVideosByType(video domain.Video) (*[]domain.Video, int, error) {
	var videos = new([]domain.Video)
	var err error
	var total int
	if video.Content != "" {
		query := bson.M{
			"video_type": video.VideoType,
			"$or": []interface{}{
				bson.M{"name": bson.M{"$regex": video.Content, "$options": "i"}},
				bson.M{"video_type": bson.M{"$regex": video.Content, "$options": "i"}},
			},
		}
		err = dao.DB.C("video").Find(query).Skip((video.Page - 1) * video.PageSize).Limit(video.PageSize).All(videos)
		total, err = dao.DB.C("video").Find(query).Count()
	} else {
		err = dao.DB.C("video").Find(bson.M{"video_type": video.VideoType}).Skip((video.Page - 1) * video.PageSize).Limit(video.PageSize).All(videos)
		total, err = dao.DB.C("video").Find(bson.M{"video_type": video.VideoType}).Count()
	}
	if err != nil {
		logrus.Error(err)
		return nil, 0, err
	}
	return videos, total, nil
}
func (dao *VideoDao) GetVideos(video domain.Video) (*[]domain.Video, int, error) {
	var videos = new([]domain.Video)
	var err error
	var total int
	if video.Content != "" {
		query := bson.M{
			"$or": []interface{}{
				bson.M{"name": bson.M{"$regex": video.Content, "$options": "i"}},
				bson.M{"video_type": bson.M{"$regex": video.Content, "$options": "i"}},
			},
		}
		err = dao.DB.C("video").Find(query).Skip((video.Page - 1) * video.PageSize).Limit(video.PageSize).All(videos)
		total, err = dao.DB.C("video").Find(query).Count()
	} else {
		err = dao.DB.C("video").Find(nil).Skip((video.Page - 1) * video.PageSize).Limit(video.PageSize).All(videos)
		total, err = dao.DB.C("video").Find(nil).Count()

	}
	if err != nil {
		logrus.Error(err)
		return nil, 0, err
	}
	return videos, total, nil
}

func (dao *VideoDao) GetVideosById(videoId string) (*domain.Video, error) {
	var video = new(domain.Video)
	err := dao.DB.C("video").FindId(bson.ObjectIdHex(videoId)).One(video)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return video, nil
}
func (dao *VideoDao) GetVideosByObjectId(videoId bson.ObjectId) (*domain.Video, error) {
	var video = new(domain.Video)
	err := dao.DB.C("video").FindId(videoId).One(video)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return video, nil
}

func (dao *VideoDao) GetVideosByInt(videoInt int) (*domain.Video, error) {
	var video = new(domain.Video)
	err := dao.DB.C("video").Find(bson.M{"video_int": videoInt}).One(video)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return video, nil
}

func (dao *VideoDao) Insert(form *domain.Video) error {
	all := dao.GetAll()
	form.VideoInt = len(*all) + 1
	err := dao.DB.C("video").Insert(form)
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (dao *VideoDao) GetAll() *[]domain.Video {
	var form = new([]domain.Video)
	err := dao.DB.C("video").Find(nil).All(form)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return form
}

func (dao *VideoDao) Update(video *domain.Video) error {
	err := dao.DB.C("video").Update(bson.M{"_id": video.Id}, video)
	if err != nil {
		logrus.Error(err)
	}
	return err
}
func (dao *VideoDao) DeleteById(videoId string) error {
	err := dao.DB.C("video").RemoveId(bson.ObjectIdHex(videoId))
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (dao *VideoDao) GetVideosByClickMore(tag string, count int) (*[]domain.Video, error) {
	var videos = new([]domain.Video)
	var err error
	if tag != "" {
		query := bson.M{"video_type": bson.M{"$regex": tag, "$options": "i"}}
		err = dao.DB.C("video").Find(query).Limit(count).Sort("-click").All(videos)
	} else {
		err = dao.DB.C("video").Find(nil).Limit(count).Sort("-click").All(videos)
	}
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return videos, nil
}
