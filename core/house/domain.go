package houses

import "time"

type House struct {
	Id          int64     `db:"id" json:"id,omitempty"  form:"id"`
	HouseholdId int64     `db:"household_id" json:"household_id,omitempty" form:"household_id"`
	HouseId     string    `db:"house_id" json:"house_id,omitempty" form:"house_id"`
	CreatedAt   time.Time `db:"create_time" json:"create_time,omitempty" form:"create_time"`
	UpdatedAt   time.Time `db:"update_time" json:"update_time,omitempty" form:"update_time"`
	Page         int           `db:"omitempty" json:"page,omitempty" form:"page"`
	PageSize     int           `db:"omitempty" json:"page_size,omitempty" form:"page_size"`
}
