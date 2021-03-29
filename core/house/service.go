package houses

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"management/infra/base"
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
	var house *House
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := HouseDao{runner: runner}
		house = dao.GetHousesById(houseId)
		if house == nil {
			return errors.New("该房子不存在")
		}
		return nil
	})
	return house, err
}

func (h houseService) SelectByHouseholdId(household int) (houses *[]House, err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := HouseDao{runner: runner}
		houses = dao.GetUserHouses(household)
		return nil
	})
	return

}

func (h houseService) Create(house *House) (err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		house.UpdatedAt = time.Now()
		house.CreatedAt = time.Now()
		dao := HouseDao{runner: runner}
		_, err := dao.runner.Insert(house)
		if err != nil {
			log.Error(err)
		}
		return err
	})
	return
}

func (h houseService) Update(house *House) error {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		house.UpdatedAt = time.Now()
		dao := HouseDao{runner: runner}
		_, err := dao.runner.Update(house)
		if err != nil {
			log.Error(err)
		}
		return err
	})
	return err

}
