package orders

import "time"

type Order struct {
	Id            int64         `db:"id" json:"id,omitempty" form:"id"`
	HouseholdId   int64         `db:"household_id" json:"household_id,omitempty" form:"household_id"`
	HouseholdName string        `db:"household_name" json:"household_name,omitempty" form:"household_name"`
	EmployeeId    int64         `db:"employee_id" json:"employee_id,omitempty" form:"employee_id"`
	HouseId       int64         `db:"house_id" json:"house_id,omitempty" form:"house_id"`
	EmployeeName  string        `db:"employee_name" json:"employee_name,omitempty" form:"employee_name"`
	Type          string        `db:"type" json:"type,omitempty" form:"type"`
	Emergency     int           `db:"emergency" json:"emergency,omitempty" form:"emergency"`
	Stage         int           `db:"stage" json:"stage,omitempty" form:"stage"`
	EvaluationId  int64         `db:"evaluation_id" json:"evaluation_id,omitempty" form:"evaluation_id"`
	CreatedAt     time.Time     `db:"create_time" json:"update_time,omitempty" form:"update_time"`
	UpdatedAt     time.Time     `db:"update_time" json:"update_time,omitempty" form:"update_time"`
	Note          string        `db:"omitempty" json:"note,omitempty" form:"note"`
	OrderStage    *[]OrderStage `db:"omitempty" json:"order_stage,omitempty" form:"order_stage"`
	Page          int           `db:"omitempty" json:"page,omitempty" form:"page"`
	PageSize      int           `db:"omitempty" json:"page_size,omitempty" form:"page_size"`
	Level         int           `db:"level,omitempty" json:"level,omitempty" form:"level"`
}

type OrderStage struct {
	Id        int64     `db:"id" json:"id,omitempty"`
	OrderId   int64     `db:"order_id" json:"order_id,omitempty"`
	Note      string    `db:"note" json:"note,omitempty"`
	Stage     int       `db:"stage" json:"stage,omitempty"`
	CreatedAt time.Time `db:"create_time" json:"update_time,omitempty"`
	UpdatedAt time.Time `db:"update_time" json:"update_time,omitempty"`
}

type Evaluation struct {
	Id          int64     `db:"id" json:"id,omitempty"`
	OrderId     int64     `db:"order_id" json:"order_id,omitempty"`
	EmployeeId  int64     `db:"employee_id" json:"employee_id,omitempty"`
	Note        string    `db:"note" json:"note,omitempty"`
	Level       int       `db:"level" json:"level,omitempty"`
	TimelyScore int       `db:"timely_score" json:"timely_score,omitempty"`
	CreatedAt   time.Time `db:"create_time" json:"update_time,omitempty"`
	UpdatedAt   time.Time `db:"update_time" json:"update_time,omitempty"`
}

type OrderSlice []Order

func (o OrderSlice) Len() int           { return len(o) }
func (o OrderSlice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o OrderSlice) Less(i, j int) bool { return o[i].UpdatedAt.Before(o[j].UpdatedAt) }
