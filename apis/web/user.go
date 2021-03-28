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

func (u *UserApi) Init() {
	u.service = users.GetUserService()
	groupRouter := base.Iris().Party("/api/1.0/user")
	groupRouter.Post("/register", u.register)
	groupRouter.Post("/login", u.login)
	groupRouter.Put("/reset", u.reset)
	groupRouter.Get("/message", loginMeddle, u.message)
	groupRouter.Post("/personal", loginMeddle, u.personal)
	groupRouter.Put("/click", employeeMeddle, u.click)
	groupRouter.Get("/order", loginMeddle, u.order)
}

// 用户注册
func (u *UserApi) register(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	//获取请求参数
	user := users.User{}
	err := ctx.ReadJSON(&user)
	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	if !checkUser(user) {
		r.Code = base.ResError
		r.Message = "缺少参数"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	_, err = u.service.GetUserByPhone(user.Phone, user.UserType)
	if err == nil {
		r.Code = base.ResError
		r.Message = "用户已注册"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	//创建用户
	err = u.service.Create(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		logrus.Error(err)
		ctx.JSON(&r)
		return
	}

	ctx.JSON(&r)

}

func (u *UserApi) login(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	//获取请求参数
	user := &users.User{}
	err := ctx.ReadJSON(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	user, err = u.service.Login(user.Phone, user.Password, user.UserType)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	token, _ := core.GenerateToken(*user)
	r.Data = map[string]interface{}{
		"token":            token,
		"default_password": user.Password == user.Id_code[len(user.Id_code)-6:],
	}
	ctx.JSON(&r)

}

func (u *UserApi) order(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	userId, _ := strconv.ParseInt(ctx.GetHeader("user_id"), 10, 64)
	service := orders.GetOrderService()
	finishOrder, doingOrders, err := service.GetOrdersByUser(userId)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	r.Data = map[string]interface{}{
		"token":           refreshToken(ctx),
		"processing_list": doingOrders,
		"complete_list":   finishOrder,
	}
	ctx.JSON(&r)

}

func (u *UserApi) reset(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	user := users.User{}
	err := ctx.ReadJSON(&user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	err = u.service.ResetPassword(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	ctx.JSON(&r)
}

func (u *UserApi) message(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	phone := ctx.GetHeader("phone")
	userType, err := strconv.Atoi(ctx.GetHeader("user_type"))
	user, err := u.service.GetUserByPhone(phone, userType)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	r.Data = map[string]interface{}{
		"token": refreshToken(ctx),
		"user":  user,
	}
	ctx.JSON(&r)
}

func (u *UserApi) personal(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	//获取请求参数
	user := users.User{}
	err := ctx.ReadJSON(&user)

	if err != nil {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	user.Phone = ctx.GetHeader("phone")
	user.UserType, _ = strconv.Atoi(ctx.GetHeader("user_type"))
	err = u.service.Edit(user)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		logrus.Error(err)
	}

	r.Data = map[string]interface{}{
		"token": refreshToken(ctx),
		"user":  user,
	}
	ctx.JSON(&r)
}

func (u *UserApi) click(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	phone := ctx.GetHeader("phone")
	userType, _ := strconv.Atoi(ctx.GetHeader("user_type"))
	if userType != 2 {
		r.Code = base.ResError
		r.Message = "该用户不是员工"
		ctx.JSON(&r)
	}

	oldUser, err := u.service.GetUserByPhone(phone, userType)
	if oldUser.State == 2 {
		oldUser.State = 1
	} else {
		oldUser.State = 2
	}
	oldUser.UpdatedAt = time.Now()
	err = u.service.Edit(*oldUser)
	if err != nil {
		r.Code = base.ResError
		r.Message = "打卡失败"
		logrus.Error(err)
	}

	r.Data = map[string]interface{}{
		"token": refreshToken(ctx),
	}
	ctx.JSON(&r)
}

func checkUser(user users.User) bool {
	if user.Name == "" || user.Id_code == "" || user.Phone == "" || user.Email == "" ||
		user.UserType == 0 {
		return false
	}
	return true
}

func refreshToken(ctx iris.Context) string {
	phone := ctx.GetHeader("phone")
	userId, _ := strconv.ParseInt(ctx.GetHeader("user_id"), 10, 64)
	userType, _ := strconv.Atoi(ctx.GetHeader("user_type"))
	user := users.User{
		Phone:    phone,
		Id:       userId,
		UserType: userType,
	}
	token, _ := core.GenerateToken(user)
	return token
}
