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
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	err = h.service.Create(house)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		logrus.Error(err)
	}
	ctx.JSON(&r)
}

// 登记房子
func (h *HouseApi) register(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}

	//获取请求参数
	house := &houses.House{}
	err := ctx.ReadJSON(house)
	if err != nil || house.HouseId == "" || house.HouseholdId == 0 {
		r.Code = base.ResError
		r.Message = "字段或字段值格式错误"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	hou, err := h.service.SelectByHouseId(house.HouseId)

	if err != nil || hou == nil || hou.HouseholdId != 0 {
		r.Code = base.ResError
		r.Message = "该房子不存在或已经被注册"
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	hou.HouseholdId = house.HouseholdId
	err = h.service.Update(*hou)
	if err != nil {
		r.Code = base.ResError
		r.Message = err.Error()
		logrus.Error(err)
	}
	ctx.JSON(&r)

}
