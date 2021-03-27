package houses

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type HouseDao struct {
	runner *dbx.TxRunner
}

// 通过用户获取房子
func (dao *HouseDao) GetUserHouses(userId int) *[]House {
	form := &[]House{}
	err := dao.runner.Find(form, "select * from house where household_id=?", userId)
	if err != nil {
		logrus.Error(err)
	}
	return form
}

// 通过用户获取房子
func (dao *HouseDao) GetHousesById(houseId string) *House {
	form := &House{}
	ok, err := dao.runner.Get(form, "select * from house where house_id=?", houseId)
	if err != nil || !ok {
		return nil
	}
	return form
}

func (dao *HouseDao) Insert(house House) (int64, error) {
	rs, err := dao.runner.Insert(house)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

func (dao *HouseDao) Update(house House) (int64, error) {
	rs, err := dao.runner.Update(house)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}
