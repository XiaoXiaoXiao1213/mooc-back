package houses

import "time"

type House struct {
	Id          int64     `db:"id" json:"id,omitempty"`
	HouseholdId int64     `db:"household_id" json:"household_id,omitempty"`
	HouseId     string    `db:"house_id" json:"house_id,omitempty"`
	CreatedAt   time.Time `db:"create_time" json:"create_time,omitempty"`
	UpdatedAt   time.Time `db:"update_time" json:"update_time,omitempty"`
	Page         int           `db:"omitempty" json:"page,omitempty"`
	PageSize     int           `db:"omitempty" json:"page_size,omitempty"`
}
