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
		dao := HouseDao{Runner: runner}
		house = dao.GetHousesById(houseId)
		if house == nil {
			return errors.New("该房子不存在")
		}
		return nil
	})
	return house, err
}

func (h houseService) SelectByHouseholdId(household int64) (houses *[]House, err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := HouseDao{Runner: runner}
		houses = dao.GetUserHouses(household)
		return nil
	})
	return

}

func (h houseService) Create(house *House) (err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		house.UpdatedAt = time.Now()
		house.CreatedAt = time.Now()
		dao := HouseDao{Runner: runner}
		_, err := dao.Runner.Insert(house)
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
		dao := HouseDao{Runner: runner}
		_, err := dao.Runner.Update(house)
		if err != nil {
			log.Error(err)
		}
		return err
	})
	return err

}

func (h houseService) DeleteHousesByHouseId(houseId string) (err error) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := HouseDao{Runner: runner}
		_, err := dao.DeleteByHouseId(houseId)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	})
	return
}

func (h houseService) GetHousesByCond(cond House) (orders *[]House, count int, err error) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := HouseDao{Runner: runner}
		if cond.Page == 0 {
			cond.Page = 1
		}
		if cond.PageSize < 1 || cond.PageSize > 10 {
			cond.PageSize = 10
		}
		log.Error(cond)
		orders, count = dao.GetByCond(cond)
		return nil
	})
	return
}
