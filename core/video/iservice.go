package videos

import (
	"gopkg.in/mgo.v2/bson"
	"management/core/domain"
	"management/infra/base"
)

var IVideoService VideoService

//用于对外暴露账户应用服务，唯一的暴露点
func GetVideoService() VideoService {
	base.Check(IVideoService)
	return IVideoService
}

type VideoService interface {
	Create(video domain.Video) error
	Edit(video domain.Video) error
	GetVideoByType(video domain.Video) (*[]domain.Video, int, error)
	DeleteVideoById(videoId string) error
	GetVideoById(videoId string) (*domain.Video, error)
	GetVideoByObjectId(videoId bson.ObjectId) (*domain.Video, error)
	GetRecommendVideo() (*[]domain.Video, error)

	GetVideos(video domain.Video) (*[]domain.Video, int, error)
	ClickVideo(video domain.VideoClick) error
	GetHotVideo(video domain.Video) (*[]domain.Video, int, error)
	GetHistoryVideo(user domain.User) (*[]domain.Video, int, error)
}
