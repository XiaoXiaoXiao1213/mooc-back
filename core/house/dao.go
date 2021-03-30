package houses

import (
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"strconv"
)

type HouseDao struct {
	Runner *dbx.TxRunner
}

// 通过用户获取房子
func (dao *HouseDao) GetUserHouses(userId int64) *[]House {
	form := &[]House{}
	err := dao.Runner.Find(form, "select * from house where household_id=?", userId)
	if err != nil {
		logrus.Error(err)
	}
	return form
}

// 通过用户获取房子
func (dao *HouseDao) GetHousesById(houseId string) *House {
	form := &House{}
	ok, err := dao.Runner.Get(form, "select * from house where house_id=?", houseId)
	if err != nil || !ok {
		return nil
	}
	return form
}

func (dao *HouseDao) Insert(house House) (int64, error) {
	rs, err := dao.Runner.Insert(house)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

func (dao *HouseDao) Update(house House) (int64, error) {
	rs, err := dao.Runner.Update(house)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

func (dao *HouseDao) DeleteByHouseId(houseId string) (int64, error) {
	rs, err := dao.Runner.Exec("delete from house where house_id=? limit 1",houseId)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

func (dao *HouseDao) GetByCond(house House) (*[]House, int) {
	form := []House{}
	sql := "select * from `house` where 1"
	if house.HouseId != "" {
		sql += " and houseId=" + house.HouseId
	}
	if house.HouseholdId != 0 {
		sql += " and household_id=" + strconv.FormatInt(house.HouseholdId, 10)
	}

	err := dao.Runner.Find(&form, sql)
	if err != nil {
		log.Error(err)
		return nil, 0
	}
	count := len(form)
	err = dao.Runner.Find(&form, sql+" limit ?,?", house.Page-1, house.PageSize)
	if err != nil {
		log.Error(err)
		return nil, 0
	}
	return &form, count
}
