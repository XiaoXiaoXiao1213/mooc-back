package houses

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var _ HouseService = new(houseService)
var once sync.Once

func init() {
	once.Do(func() {
		IHouseService = new(houseService)
	})
}

type houseService struct {
}

func (h houseService) SelectByHouseId(houseId string) (*House, error) {
	dao := HouseDao{}
	house := dao.GetHousesById(houseId)
	if house == nil {
		return nil, errors.New("该房子不存在")
	}
	return house, nil
}

func (h houseService) SelectByHouseholdId(household int) (*[]House, error) {
	dao := HouseDao{}
	houses := dao.GetUserHouses(household)
	if houses == nil {
		return nil, errors.New("该房子不存在")
	}
	return houses, nil

}
func (h houseService) Create(house House) error {
	house.UpdatedAt = time.Now()
	house.CreatedAt = time.Now()
	dao := HouseDao{}
	_, err := dao.runner.Insert(house)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (h houseService) Update(house House) error {
	house.UpdatedAt = time.Now()
	dao := HouseDao{}
	_, err := dao.runner.Update(house)
	if err != nil {
		log.Error(err)
	}
	return err

}
