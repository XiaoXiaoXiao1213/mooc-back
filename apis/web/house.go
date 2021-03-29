package web

import (
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	houses "management/core/house"
	"management/infra"
	"management/infra/base"
)

func init() {
	infra.RegisterApi(new(HouseApi))
}

type HouseApi struct {
	service houses.HouseService
}

func (h *HouseApi) Init() {
	h.service = houses.GetHouseService()
	groupRouter := base.Iris().Party("/api/1.0/house")
	groupRouter.Post("/create", employeeMeddle, h.create)
	groupRouter.Post("/register", employeeMeddle, h.register)
}

// 创建房子
func (h *HouseApi) create(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	//获取请求参数
	house := houses.House{}
	err := ctx.ReadJSON(&house)
	if err != nil || house.HouseId == "" {
		if err != nil {
			logrus.Error(err)
		}
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		return
	}

	err = h.service.Create(&house)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
	}
	ctx.JSON(&r)
}

// 登记房子
func (h *HouseApi) register(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	//获取请求参数
	houseCond := houses.House{}
	err := ctx.ReadJSON(&houseCond)
	if err != nil || houseCond.HouseId == "" || houseCond.HouseholdId == 0 {
		if err != nil {
			logrus.Error(err)
		}
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		return
	}

	house, err := h.service.SelectByHouseId(houseCond.HouseId)
	if err != nil || house == nil || house.HouseholdId != 0 {
		if err != nil {
			logrus.Error(err)
		}
		r.Code = base.ResError
		r.Message = "该房子不存在或已经被注册"
		ctx.JSON(&r)
		return
	}

	house.HouseholdId = houseCond.HouseholdId
	err = h.service.Update(house)
	if err != nil {
		logrus.Error(err)
		r.Code = base.ResError
		r.Message = err.Error()
	}
	ctx.JSON(&r)
}
