package domain

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id            bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty" form:"id"`
	UserInt       int             `bson:"user_int,omitempty" json:"user_int,omitempty" form:"user_int"`
	Name          string          `bson:"name,omitempty" json:"name,omitempty" form:"name"`
	Password      string          `bson:"password,omitempty" json:"password,omitempty" form:"password"`
	Phone         string          `bson:"phone,omitempty" json:"phone,omitempty" form:"phone"`
	Image         string          `bson:"image,omitempty" json:"image,omitempty" form:"image"`
	Intro         string          `bson:"intro,omitempty" json:"intro,omitempty" form:"intro"`
	InterestsType []string        `bson:"interests_type,omitempty" json:"interests_type,omitempty" form:"interests_type"`
	CFInterests   []bson.ObjectId `bson:"c_f_interests,omitempty" json:"c_f_interests,omitempty" form:"c_f_interests"`
	CBInterests   []bson.ObjectId `bson:"c_b_interests,omitempty" json:"c_b_interests,omitempty" form:"c_b_interests"`
	Page          int             `bson:"page,omitempty" json:"page,omitempty" form:"page"`
	PageSize      int             `bson:"page_size,omitempty" json:"page_size,omitempty" form:"page_size"`
}
type PasswordUser struct {
	OldPassword    string `json:"old_password" form:"old_password"`
	NewPassword    string `json:"new_password" form:"new_password"`
	ConfitPassword string `json:"confit_password" form:"confit_password"`
}
