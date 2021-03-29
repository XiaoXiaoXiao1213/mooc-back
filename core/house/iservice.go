package houses

import (
	"management/infra/base"
)

var IHouseService HouseService

//用于对外暴露账户应用服务，唯一的暴露点
func GetHouseService() HouseService {
	base.Check(IHouseService)
	return IHouseService
}

type HouseService interface {
	Create(house *House) error
	Update(house *House) error
	SelectByHouseId(houseId string) (*House, error)
	SelectByHouseholdId(household int) (*[]House, error)
	GetHousesByCond(cond House) (house *[]House, count int, err error)
	DeleteHousesByHouseId(houseId string) (err error)

}
