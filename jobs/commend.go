package jobs

import (
	"encoding/csv"
	"github.com/go-redsync/redsync"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"management/core/cf_lib"
	"management/core/dao"
	"management/infra"
	"management/infra/base"
	"os"
	"strconv"
	"time"
)

type CommendJobStarter struct {
	infra.BaseStarter
	ticker *time.Ticker
	mutex  *redsync.Mutex
}

func (r *CommendJobStarter) Init(ctx infra.StarterContext) {
	r.ticker = time.NewTicker(time.Hour)
}

func (r *CommendJobStarter) Start(ctx infra.StarterContext) {
	go func() {
		for {
			_ = <-r.ticker.C
			db := base.MgoDatabase()
			userDao := dao.UserDao{db}
			videoDao := dao.VideoDao{db}
			createClickCsv()
			users := userDao.GetAll()
			log.Error("users", users)
			if users != nil {
				for _, user := range *users {
					cbVideoId := []bson.ObjectId{}
					if len(user.InterestsType) == 0 {
						//获取最热门的6个视频
						videos, err := videoDao.GetVideosByClickMore("", 18)
						log.Error("videos", videos)
						if err == nil {
							for _, video := range *videos {
								log.Error("video", video)
								log.Error("video", video.Id)

								cbVideoId = append(cbVideoId, video.Id)
							}
						} else {
							log.Error(err)
						}
					} else if len(user.InterestsType) == 1 {
						videos, err := videoDao.GetVideosByClickMore(user.InterestsType[0], 18)
						if err == nil {
							for _, video := range *videos {
								cbVideoId = append(cbVideoId, video.Id)
							}
						} else {
							log.Error(err)
						}
					} else if len(user.InterestsType) == 2 {
						//获取最热门的6个视频
						videos, err := videoDao.GetVideosByClickMore(user.InterestsType[0], 9)
						if err == nil {
							for _, video := range *videos {
								cbVideoId = append(cbVideoId, video.Id)
							}
						} else {
							log.Error(err)
						}
						videos, err = videoDao.GetVideosByClickMore(user.InterestsType[2], 9)
						if err == nil {
							for _, video := range *videos {
								cbVideoId = append(cbVideoId, video.Id)
							}
						} else {
							log.Error(err)
						}
					} else {
						//获取最热门的6个视频
						videos, err := videoDao.GetVideosByClickMore(user.InterestsType[0], 6)
						if err == nil {
							for _, video := range *videos {
								cbVideoId = append(cbVideoId, video.Id)
							}
						} else {
							log.Error(err)
						}
						videos, err = videoDao.GetVideosByClickMore(user.InterestsType[1], 6)
						if err == nil {
							for _, video := range *videos {
								cbVideoId = append(cbVideoId, video.Id)
							}
						} else {
							log.Error(err)
						}
						videos, err = videoDao.GetVideosByClickMore(user.InterestsType[2], 6)
						if err == nil {
							for _, video := range *videos {
								cbVideoId = append(cbVideoId, video.Id)
							}
						} else {
							log.Error(err)
						}
					}

					user.CBInterests = cbVideoId
					log.Error("d", user.CBInterests)
					//cf
					cf := cf_lib.GetItemCF()
					cf.DoCalculate()
					recommend := cf.Recommend(strconv.Itoa(user.UserInt))
					cfInterests := []bson.ObjectId{}
					for _, videoInt := range recommend {
						i, _ := strconv.Atoi(videoInt)
						video, _ := videoDao.GetVideosByInt(i)
						cfInterests = append(cfInterests, video.Id)
					}

					user.CFInterests = cfInterests
					userDao.Update(&user)
				}
				os.Remove("/Users/xiao_xiaoxiao/Desktop/web/mooc-web/core/data.csv")
			} else {
				log.Info("users==nil")
			}

		}
	}()

}

func (r *CommendJobStarter) Stop(ctx infra.StarterContext) {
	r.ticker.Stop()
}

func createClickCsv() {
	//创建文件
	f, err := os.Create("/Users/xiao_xiaoxiao/Desktop/web/mooc-web/core/data.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// 写入UTF-8 BOM
	f.WriteString("\xEF\xBB\xBF")
	//创建一个新的写入文件流
	w := csv.NewWriter(f)
	db := base.MgoDatabase()
	clickDao := dao.UserClickDao{DB: db}
	res := clickDao.GetAll()
	data := [][]string{}
	for _, click := range *res {
		content := []string{
			strconv.Itoa(click.UserId),
			strconv.Itoa(click.VideoId),
			strconv.Itoa(click.Click),
		}
		data = append(data, content)
	}
	//写入数据
	log.Error(data)
	w.WriteAll(data)
	w.Flush()
}
