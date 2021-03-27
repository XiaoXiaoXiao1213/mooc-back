package web

import (
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"management/core"
	"management/core/orders"
	"management/core/users"
	"management/infra"
	"management/infra/base"
	"strconv"
	"time"
)

func init() {
	infra.RegisterApi(new(UserApi))
}

type UserApi struct {
	service users.UserService
}

func (a *UserApi) Init() {
	a.service = users.GetUserService()
	groupRouter := base.Iris().Party("/api/1.0/user")
	groupRouter.Post("/register", a.register)
	groupRouter.Post("/login", loginMeddle, a.login)
	groupRouter.Put("/reset", a.reset)
	groupRouter.Get("/message/{phone}/{user_type}", loginMeddle, a.message)
	groupRouter.Post("/personal", loginMeddle, a.personal)
	groupRouter.Put("/click", loginMeddle, a.click)
	groupRouter.Get("/order", loginMeddle, a.order)
}

// 用户注册
func (a *UserApi) register(ctx iris.Context) {
	//获取请求参数
	user := users.User{}
	err := ctx.ReadJSON(&user)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	// TODO 参数检测
	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	//创建用户
	err = a.service.Create(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		logrus.Error(err)
		ctx.JSON(&r)
		return
	}

	ctx.JSON(&r)

}

func (a *UserApi) login(ctx iris.Context) {
	//获取请求参数，
	user := users.User{}
	err := ctx.ReadJSON(&user)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	// TODO 参数检测
	err = a.service.Login(user.Phone, user.Password, user.UserType)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	token, err := core.GenerateToken(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	r.Data = map[string]string{
		"token": token,
	}
	ctx.JSON(&r)

}

func (u *UserApi) order(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	//获取请求参数
	phone := ctx.GetHeader("phone")
	userId, _ := strconv.ParseInt(ctx.GetHeader("user_id"), 10, 64)
	userType, _ := strconv.Atoi(ctx.GetHeader("user_type"))
	user := users.User{
		Phone:    phone,
		Id:       userId,
		UserType: userType,
	}
	service := orders.GetOrderService()
	finishOrder, doingOrders, err := service.GetOrdersByUser(userId)

	token, err := core.GenerateToken(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	r.Data = map[string]interface{}{
		"token":           token,
		"processing_list": doingOrders,
		"complete_list":   finishOrder,
	}
	ctx.JSON(&r)

}

func (a *UserApi) reset(ctx iris.Context) {
	user := users.User{}
	err := ctx.ReadJSON(&user)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	err = a.service.ResetPassword(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	ctx.JSON(&r)

}

func (a *UserApi) message(ctx iris.Context) {
	phone := ctx.Params().Get("phone")
	userType, err := strconv.Atoi(ctx.Params().Get("user_type"))
	if err != nil {
		r := base.Res{
			Code:    base.ResError,
			Message: "地址格式错误",
		}
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	r := base.Res{
		Code: base.ResCodeOk,
	}
	user, err := a.service.GetUserByPhone(phone, userType)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	token, err := core.GenerateToken(*user)
	r.Data = map[string]interface{}{
		"token": token,
		"user":  user,
	}
	ctx.JSON(&r)
}

func (a *UserApi) personal(ctx iris.Context) {
	//获取请求参数
	phone := ctx.GetHeader("phone")
	userType := ctx.GetHeader("user_type")
	atoi, _ := strconv.Atoi(userType)
	user := users.User{}
	err := ctx.ReadJSON(&user)
	user.Phone = phone
	user.UserType = atoi
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	err = a.service.Edit(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		logrus.Error(err)
	}

	token, err := core.GenerateToken(user)
	r.Data = map[string]interface{}{
		"token": token,
		"user":  user,
	}
	ctx.JSON(&r)
}

func (a *UserApi) click(ctx iris.Context) {
	phone := ctx.GetHeader("phone")
	userType := ctx.GetHeader("user_type")
	atoi, _ := strconv.Atoi(userType)
	if atoi != 2 {
		r := base.Res{
			Code:    base.ResError,
			Message: "该用户不是员工",
		}
		ctx.JSON(&r)
	}
	//获取请求参数
	user := users.User{Phone: phone, UserType: atoi}
	r := base.Res{
		Code: base.ResCodeOk,
	}

	//创建用户
	u, err := a.service.GetUserByPhone(phone, atoi)
	if u.State == 2 {
		u.State = 1
	} else {
		u.State = 2
	}
	u.UpdatedAt = time.Now()
	err = a.service.Edit(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = "打卡失败"
		logrus.Error(err)
	}

	token, err := core.GenerateToken(user)
	r.Data = map[string]interface{}{
		"token": token,
	}
	ctx.JSON(&r)
}
