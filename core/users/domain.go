package users

import "time"

type User struct {
	Id         int64     `db:"id" json:"id,omitempty" form:"id"`
	Id_code    string    `db:"id_code" json:"id_code,omitempty" form:"id_code"`
	Name       string    `db:"name" json:"name,omitempty" form:"name"`
	Sex        int       `db:"sex" json:"sex,omitempty" form:"sex"`
	Birthday   time.Time `db:"birthday" json:"birthday,omitempty" form:"birthday"`
	Phone      string    `db:"phone" json:"phone,omitempty" form:"phone"`
	Wechat     string    `db:"wechat" json:"wechat,omitempty" form:"wechat"`
	Email      string    `db:"email" json:"email,omitempty" form:"email"`
	Score      int       `db:"score" json:"score,omitempty" form:"score"`
	State      int       `db:"state" json:"state,omitempty" form:"state"`
	SuperState int       `db:"super_state" json:"super_state,omitempty" form:"super_state"` //	 0	1:非管理员  2:管理员
	Skills     string    `db:"skills" json:"skills,omitempty" form:"skills"`           // 	“”(可修改)	技能 用数字映射
	Password   string    `db:"password" json:"password,omitempty" form:"password"`       // varchar	身份证后6位	密码
	Num        int       `db:"num" json:"num,omitempty" form:"num"`                 // 	默认为1	家庭人数
	UserType   int       `db:"user_type" json:"user_type,omitempty" form:"user_type"`     // 	用户类型	1为住户 2为员工
	CreatedAt  time.Time `db:"create_time" json:"update_time,omitempty" form:"update_time"`
	UpdatedAt  time.Time `db:"update_time" json:"update_time,omitempty" form:"update_time"`
	Page       int       `db:"omitempty" json:"page,omitempty" form:"page"`
	PageSize   int       `db:"omitempty" json:"page_size,omitempty" form:"page_size"`
}

type EmployeeScore struct {
	Id             int64   `db:"id" json:"id,omitempty"`
	EmployeeId     int64   `db:"employee_id" json:"employee_id,omitempty"`
	OrderScore     int     `db:"order_score" json:"order_score,omitempty"`
	TimelyScore    int     `db:"timely_score" json:"timely_score,omitempty"`
	DifficultScore int     `db:"difficult_score" json:"difficult_score,omitempty"`
	NonUrgentScore float64 `db:"non_urgent_score" json:"non_urgent_score,omitempty"`
	UrgentScore    float64 `db:"urgent_score" json:"urgent_score,omitempty"`
	State          int     `db:"state" json:"state,omitempty"`
	DoingOrder     int     `db:"doing_order" json:"doing_order,omitempty"`
	OrderCount     int     `db:"order_count" json:"order_count,omitempty"`
}
