package web

import (
	"github.com/kataras/iris"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"management/core/common"
	"management/core/domain"
	"management/core/users"
	videos "management/core/video"
	"management/infra"
	"management/infra/base"
	"strconv"
)

func init() {
	infra.RegisterApi(new(VideoApi))
}

type VideoApi struct {
	service videos.VideoService
}

func (v *VideoApi) Init() {
	v.service = videos.GetVideoService()

	groupRouter := base.Iris().Party("/api/1.0/video")
	groupRouter.Use(Cors)
	// common
	common := groupRouter.Party("/")
	{
		common.Options("*", func(ctx iris.Context) {
			ctx.Next()
		})
	}
	groupRouter.Post("/create", v.createVideo)
	groupRouter.Get("/id/{id}", v.getVideoById)
	groupRouter.Get("/type/{type}",v.getVideoByType)
	groupRouter.Get("/all", v.getVideo)
	groupRouter.Put("/delete/{id}", loginMeddle, v.deleteVideo)
	groupRouter.Put("/click/{id}", loginMeddle, v.clickVideo)
	groupRouter.Get("/hot", v.getHotVideo)

}

// 用户行为记录 点击课程
func (v *VideoApi) clickVideo(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	//获取请求参数
	token := ctx.GetHeader("Authorization")
	userId := ""
	if token != "" {
		user, _ := common.ParseToken(token)
		userId = (*user)["user_id"]
	}
	videoId := ctx.Params().Get("id")
	videoRes, _ := v.service.GetVideoById(videoId)
	userRes, _ := users.GetUserService().GetUserById(userId)

	click := domain.VideoClick{
		UserId:  userRes.UserInt,
		VideoId: videoRes.VideoInt,
	}
	err := v.service.ClickVideo(click)
	if err != nil {
		log.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	ctx.JSON(&r)
}

//删除课程
func (v *VideoApi) deleteVideo(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	//获取请求参数
	id := ctx.Params().Get("id")
	err := v.service.DeleteVideoById(id)
	if err != nil {
		log.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	ctx.JSON(&r)
}

// 创建课程
func (v *VideoApi) createVideo(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	//获取请求参数
	video := domain.Video{}
	err := ctx.ReadForm(&video)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		return
	}
	url, err := common.Upload(ctx, "video_img")
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "图片上传错误"
		ctx.JSON(&r)
		return
	}
	video.VideoImg = url
	if video.Id != "" {
		err = v.service.Edit(video)
	} else {
		err = v.service.Create(video)
	}
	if err != nil {
		log.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	ctx.JSON(&r)
}

// 获取课程详情
func (v *VideoApi) getVideoById(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	videoId := ctx.Params().Get("id")
	videoRes, err := v.service.GetVideoById(videoId)
	log.Error(videoRes)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = map[string]interface{}{
		"video": videoRes,
	}
	ctx.JSON(&r)
}

//根据类型获取课程
func (v *VideoApi) getVideoByType(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	videoType := ctx.Params().Get("type")
	page, _ := strconv.Atoi(ctx.FormValueDefault("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.FormValueDefault("page_size", "15"))
	content := ctx.FormValueDefault("content", "")
	cond := domain.Video{
		VideoType: videoType,
		PageSize:  pageSize,
		Page:      page,
		Content: content,
	}
	videoRes, total, err := v.service.GetVideoByType(cond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = map[string]interface{}{
		"videos": videoRes,
		"total":  total,
	}
	ctx.JSON(&r)
}

func (v *VideoApi) getHotVideo(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	videoType := ctx.FormValueDefault("type", "")
	page, _ := strconv.Atoi(ctx.FormValueDefault("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.FormValueDefault("page_size", "15"))
	content := ctx.FormValueDefault("content", "")
	cond := domain.Video{
		VideoType: videoType,
		PageSize:  pageSize,
		Page:      page,
		Content: content,
	}
	videoRes, total, err := v.service.GetHotVideo(cond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = map[string]interface{}{
		"videos": videoRes,
		"total":  total,
	}
	ctx.JSON(&r)
}

// 模糊查询课程
func (v *VideoApi) getVideo(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	page, _ := strconv.Atoi(ctx.FormValueDefault("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.FormValueDefault("page_size", "15"))
	content := ctx.FormValueDefault("content", "")
	cond := domain.Video{
		PageSize: pageSize,
		Page:     page,
		Content:  content,
	}

	videoRes, total, err := v.service.GetVideos(cond)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = map[string]interface{}{
		"videos": videoRes,
		"total":  total,
	}
	ctx.JSON(&r)
}
