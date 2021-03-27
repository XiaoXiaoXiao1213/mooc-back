package orders

import "time"

type Order struct {
	Id            int64         `db:"id" json:"id,omitempty"`
	HouseholdId   int64         `db:"household_id" json:"household_id,omitempty"`
	HouseholdName string        `db:"household_name" json:"household_name,omitempty"`
	EmployeeId    int64         `db:"employee_id" json:"employee_id,omitempty"`
	HouseId       int64         `db:"house_id" json:"house_id,omitempty"`
	EmployeeName  string        `db:"employee_name" json:"employee_name,omitempty"`
	Type          int           `db:"type" json:"type,omitempty"`
	Emergency     int           `db:"emergency" json:"emergency,omitempty"`
	Stage         int           `db:"stage" json:"stage,omitempty"`
	EvaluationId  int64         `db:"evaluation_id" json:"evaluation_id,omitempty"`
	CreatedAt     time.Time     `db:"create_time" json:"update_time,omitempty"`
	UpdatedAt     time.Time     `db:"update_time" json:"update_time,omitempty"`
	OrderStage    *[]OrderStage `db:"omitempty" json:"order_stage,omitempty"`
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
	Id         int64     `db:"id" json:"id,omitempty"`
	OrderId    int64     `db:"order_id" json:"order_id,omitempty"`
	EmployeeId int64     `db:"employee_id" json:"employee_id,omitempty"`
	Note       string    `db:"note" json:"note,omitempty"`
	Level      int64     `db:"level" json:"level,omitempty"`
	CreatedAt  time.Time `db:"create_time" json:"update_time,omitempty"`
	UpdatedAt  time.Time `db:"update_time" json:"update_time,omitempty"`
}

type OrderSlice []Order

func (o OrderSlice) Len() int           { return len(o) }
func (o OrderSlice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o OrderSlice) Less(i, j int) bool { return o[i].UpdatedAt.Before(o[j].UpdatedAt) }
