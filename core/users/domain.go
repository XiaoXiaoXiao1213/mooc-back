package users

import "time"

type User struct {
	Id         int64     `db:"id" json:"id,omitempty"`
	Id_code    string    `db:"id_code" json:"id_code,omitempty"`
	Name       string    `db:"name" json:"name,omitempty"`
	Sex        int       `db:"sex" json:"sex,omitempty"`
	Birthday   time.Time `db:"birthday" json:"birthday,omitempty"`
	Phone      string    `db:"phone" json:"phone,omitempty"`
	Wechat     string    `db:"wechat" json:"wechat,omitempty"`
	Email      string    `db:"email" json:"email,omitempty"`
	Score      int64     `db:"score" json:"score,omitempty"`
	State      int       `db:"state" json:"state,omitempty"`
	SuperState int       `db:"super_state" json:"super_state,omitempty"` //	 0	1:非管理员  2:管理员
	Skills     string    `db:"skills" json:"skills,omitempty"`           // 	“”(可修改)	技能 用数字映射
	Password   string    `db:"password" json:"password,omitempty"`       // varchar	身份证后6位	密码
	Num        int       `db:"num" json:"num,omitempty"`                 // 	默认为1	家庭人数
	UserType   int       `db:"user_type" json:"user_type,omitempty"`     // 	用户类型	1为住户 2为员工
	CreatedAt  time.Time `db:"create_time" json:"update_time,omitempty"`
	UpdatedAt  time.Time `db:"update_time" json:"update_time,omitempty"`
	Page       int       `db:"omitempty" json:"page,omitempty"`
	PageSize   int       `db:"omitempty" json:"page_size,omitempty"`
}
