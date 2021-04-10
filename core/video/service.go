package videos

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"management/core/dao"
	"management/core/domain"
	"management/infra/base"
	"strconv"
	"sync"
	"time"
)

var _ VideoService = new(videoService)
var once sync.Once

func init() {
	once.Do(func() {
		IVideoService = new(videoService)
	})
}

type videoService struct {
}

func (u *videoService) GetRecommendVideo() (*[]domain.Video, error) {
	db := base.MgoDatabase()
	videoDao := dao.VideoDao{db}
	video, err := videoDao.GetVideosByClickMore("", 8)
	if err != nil {
		log.Error(err)
		err := errors.New("获取失败")
		return nil, err
	}
	return video, nil
}

func (u *videoService) GetHotVideo(video domain.Video) (*[]domain.Video, int, error) {
	db := base.MgoDatabase()
	videoDao := dao.VideoDao{db}
	videoRes, total, err := videoDao.GetVideosByTypeAndHot(video)
	if err != nil {
		log.Error(err)
		err := errors.New("获取失败")
		return nil, 0, err
	}
	return videoRes, total, nil
}

func (u *videoService) Edit(video domain.Video) error {
	db := base.MgoDatabase()
	videoDao := dao.VideoDao{db}
	err := videoDao.Update(&video)
	if err != nil {
		log.Error(err)
		err := errors.New("更新失败")
		return err
	}
	return nil
}

func (u *videoService) GetHistoryVideo(user domain.User) (*[]domain.Video, int, error) {
	userId := user.UserInt
	db := base.MgoDatabase()
	historyDao := dao.UserHistoryDao{db}
	userHistory, err := historyDao.GetByUserId(userId)
	if err != nil {
		return nil, 0, err
	}
	videoDao := dao.VideoDao{db}
	res := new([]domain.Video)
	resMap := make(map[string]bool)
	for _, history := range *userHistory {
		video, _ := videoDao.GetVideosByInt(history.VideoId)
		_, ok := resMap[strconv.Itoa(video.VideoInt)]
		if !ok {
			*res = append(*res, *video)
			resMap[strconv.Itoa(video.VideoInt)] = true
		}
	}
	return res, len(*res), nil

}

func (u *videoService) ClickVideo(videoClick domain.VideoClick) error {
	db := base.MgoDatabase()
	clickDao := dao.UserClickDao{db}
	historyDao := dao.UserHistoryDao{db}

	if videoClick.UserId != 0 {
		res, err := clickDao.GetByUserAndVideo(videoClick)
		if err == nil && res != nil {
			res.Click += 1
			err := clickDao.Update(res)
			if err != nil {
				log.Error(err)
				return err
			}
		} else {
			videoClick.Click = 1
			err := clickDao.Insert(&videoClick)
			if err != nil {
				log.Error(err)
				return err
			}
		}
		history := domain.UserVideoHistory{
			VideoId: videoClick.VideoId,
			UserId:  videoClick.UserId,
			Time:    time.Now().String(),
		}
		_ = historyDao.Insert(&history)

	}

	videoDao := dao.VideoDao{db}
	video, _ := videoDao.GetVideosByInt(videoClick.VideoId)
	video.Click += 1
	_ = videoDao.Update(video)
	return nil
}

func (u *videoService) Create(video domain.Video) error {
	if video.VideoUrl == "" || video.Name == "" || video.VideoType == "" || video.VideoImg == "" {
		return errors.New("缺少参数")
	}
	db := base.MgoDatabase()
	videoDao := dao.VideoDao{db}
	video.Click = 1
	err := videoDao.Insert(&video)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (u *videoService) GetVideoByType(video domain.Video) (*[]domain.Video, int, error) {
	db := base.MgoDatabase()
	videoDao := dao.VideoDao{db}
	videos, total, err := videoDao.GetVideosByType(video)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}
	return videos, total, nil
}
func (u *videoService) GetVideos(video domain.Video) (*[]domain.Video, int, error) {
	db := base.MgoDatabase()
	videoDao := dao.VideoDao{db}
	videos, total, err := videoDao.GetVideos(video)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}
	return videos, total, nil
}

func (u *videoService) DeleteVideoById(videoId string) error {
	db := base.MgoDatabase()
	videoDao := dao.VideoDao{db}
	err := videoDao.DeleteById(videoId)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (u *videoService) GetVideoById(videoId string) (*domain.Video, error) {
	db := base.MgoDatabase()
	videoDao := dao.VideoDao{db}
	video, err := videoDao.GetVideosById(videoId)
	if err != nil {
		log.Error(err)
	}
	return video, err
}

func (u *videoService) GetVideoByObjectId(videoId bson.ObjectId) (*domain.Video, error) {
	db := base.MgoDatabase()
	videoDao := dao.VideoDao{db}
	video, err := videoDao.GetVideosByObjectId(videoId)
	if err != nil {
		log.Error(err)
	}
	return video, err
}
