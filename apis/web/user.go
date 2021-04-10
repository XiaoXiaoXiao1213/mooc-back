package web

import (
	"github.com/kataras/iris"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"management/core/common"
	"management/core/domain"
	"management/core/users"
	videos "management/core/video"
	"management/infra"
	"management/infra/base"
	"math/rand"
	"strconv"
)

func init() {
	infra.RegisterApi(new(UserApi))
}

type UserApi struct {
	service      users.UserService
	videoService videos.VideoService
}

func (u *UserApi) Init() {
	u.service = users.GetUserService()
	u.videoService = videos.GetVideoService()
	groupRouter := base.Iris().Party("/api/1.0/user")
	groupRouter.Use(Cors)
	groupRouter.Post("/register", u.register)
	groupRouter.Post("/login", u.login)
	groupRouter.Put("/logout", loginMeddle, u.logout)
	groupRouter.Post("/password/reset", loginMeddle, u.reset)
	groupRouter.Get("/message", loginMeddle, u.message)
	groupRouter.Post("/message/edit", loginMeddle, u.editMessage)
	groupRouter.Get("/commend/video", loginMeddle, u.getCommendVideo)
	groupRouter.Get("/commend/video/default", u.getCommendVideoDefault)
	groupRouter.Get("/history/video", loginMeddle, u.getHistoryVideo)

}

func (u *UserApi) getHistoryVideo(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	userId := ctx.GetHeader("user_id")
	userRes, err := u.service.GetUserById(userId)
	videos, total, err := u.videoService.GetHistoryVideo(*userRes)
	if err != nil {
		log.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = map[string]interface{}{
		"videos": videos,
		"total":  total,
		"token":  refreshToken(ctx),
	}
	ctx.JSON(&r)
}

// 用户注册
func (u *UserApi) register(ctx iris.Context) {
	r := base.Res{
		Code:    base.ResCodeOk,
		Message: "注册成功",
	}

	//获取请求参数
	user := domain.User{}
	err := ctx.ReadJSON(&user)

	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		return
	}

	userRes, err := u.service.Create(user)
	if err != nil {
		log.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = map[string]interface{}{
		"user": userRes,
	}
	ctx.JSON(&r)
}

// 登陆
func (u *UserApi) login(ctx iris.Context) {
	r := base.Res{
		Code:    base.ResCodeOk,
		Message: "登陆成功",
	}

	//获取请求参数
	user := &domain.User{}
	err := ctx.ReadJSON(user)
	if err != nil {
		logrus.Error("userApi login ctx.ReadJSON ：", err)
		r.Code = base.ResError
		r.Message = "参数错误"
		ctx.JSON(&r)
		return
	}
	user, err = u.service.Login(user.Phone, user.Password)

	if err != nil {
		logrus.Error("userApi login u.service.Login：", err)
		r.Code = base.ResError
		r.Message = "没有该用户或密码错误"
		ctx.JSON(&r)
		return
	}
	token, _ := common.GenerateToken(*user)
	r.Data = map[string]interface{}{
		"user":  user,
		"token": token,
	}
	ctx.JSON(&r)
}

// 退出登陆
func (u *UserApi) logout(ctx iris.Context) {
	r := base.Res{
		Code:    base.ResCodeOk,
		Message: "退出登陆成功",
	}
	ctx.JSON(&r)
}

// 修改密码
func (u *UserApi) reset(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
		Message: "修改成功",
	}
	user:= domain.User{}
	userP := domain.PasswordUser{}
	ctx.ReadJSON(&userP)
	log.Error("oldPassword",userP)
	if userP.OldPassword == "" || userP.NewPassword == "" {
		r.Code = base.ResError
		r.Message = "缺少参数"
		ctx.JSON(&r)
		return
	}
	if userP.ConfitPassword != userP.NewPassword {
		r.Code = base.ResError
		r.Message = "两次密码不一样"
		ctx.JSON(&r)
		return
	}
	userId := ctx.GetHeader("user_id")
	user.Id = bson.ObjectId(userId)
	user.Password = userP.OldPassword
	err := u.service.ResetPassword(user, userP.NewPassword)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	ctx.JSON(&r)
}

// 获取用户信息
func (u *UserApi) message(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	userId := ctx.GetHeader("user_id")
	user, err := u.service.GetUserById(userId)
	if err != nil || user == nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = map[string]interface{}{
		"user":  user,
		"token": refreshToken(ctx),
	}
	ctx.JSON(&r)
}

// 修改信息
func (u *UserApi) editMessage(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	userId := ctx.GetHeader("user_id")
	userRes, err := u.service.GetUserById(userId)
	log.Error(userRes)
	if err != nil || userRes == nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	user := domain.User{}
	ctx.ReadJSON(&user)

	err = u.service.Update(user)

	if err != nil {
		logrus.Error("u.service.Update(user):", err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}

	r.Data = map[string]interface{}{
		"user":  user,
		"token": refreshToken(ctx),
	}
	ctx.JSON(&r)
}

// 获取推荐的视频信息
func (u *UserApi) getCommendVideo(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	userId := ctx.GetHeader("user_id")
	userRes, err := u.service.GetUserById(userId)
	if err != nil || userRes == nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	cbVideoId := getSeed(userRes.CBInterests, 8)
	cfVideoId := getSeed(userRes.CFInterests, 8)
	res := append(cbVideoId, cfVideoId...)
	res = getSeed(res, 8)
	videos := []domain.Video{}
	for _, id := range res {
		video, err := u.videoService.GetVideoByObjectId(id)
		if err != nil {
			logrus.Error(" u.videoService.GetVideoById(id)", err)
			r.Code = base.ResError
			r.Message = err.Error()
			ctx.JSON(&r)
			return
		}
		videos = append(videos, *video)

	}
	r.Data = map[string]interface{}{
		"videos": videos,
		"token":  refreshToken(ctx),
	}
	ctx.JSON(&r)
}

func (u *UserApi) getCommendVideoDefault(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	video, err := u.videoService.GetRecommendVideo()
	if err != nil {
		log.Error(video)
	}
	r.Data = map[string]interface{}{
		"videos": video,
		"token":  refreshToken(ctx),
	}
	ctx.JSON(&r)
}

func getSeed(array []bson.ObjectId, count int) []bson.ObjectId {
	if len(array) <= count {
		return array
	}
	mapRes := make(map[string]bool)
	res := []bson.ObjectId{}
	for len(res) < count {
		i := rand.Intn(len(array))
		istr := strconv.Itoa(i)
		m := mapRes[istr]
		if m {
			continue
		}
		mapRes[istr] = true
		res = append(res, array[i])
	}
	return res
}
