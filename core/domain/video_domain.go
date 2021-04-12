package domain

import (
	"gopkg.in/mgo.v2/bson"
)

type Video struct {
	Id        bson.ObjectId `bson:"_id,omitempty" form:"id" json:"id"`
	VideoInt  int           `bson:"video_int,omitempty" form:"video_int" json:"video_int"`
	Name      string        `bson:"name,omitempty" form:"name" json:"name"`
	Content   string        `bson:"content,omitempty" form:"content" json:"content"`
	VideoType string        `bson:"video_type,omitempty" form:"video_type" json:"video_type"`
	VideoImg  string        `bson:"video_img,omitempty"  form:"video_img" json:"video_img"`
	VideoUrl  string        `bson:"video_url,omitempty" form:"video_url" json:"video_url"`
	Click     int           `bson:"click,omitempty" form:"click" json:"click"`
	Page      int           `bson:"page,omitempty" form:"page" json:"page"`
	PageSize  int           `bson:"page_size,omitempty" form:"page_size" json:"page_size"`
}

//
type VideoClick struct {
	Id      bson.ObjectId `bson:"_id,omitempty"  form:"id" json:"id"`
	VideoId int           `bson:"video_id,omitempty"  form:"video_id" json:"video_id"`
	UserId  int           `bson:"user_id,omitempty"  form:"user_id" json:"user_id"`
	Click   int           `bson:"click,omitempty"  form:"click" json:"click"`
}

type UserVideoHistory struct {
	Id      bson.ObjectId `bson:"_id,omitempty"  form:"id" json:"id"`
	VideoId int           `bson:"video_id,omitempty"  form:"video_id" json:"video_id"`
	UserId  int           `bson:"user_id,omitempty"  form:"user_id" json:"user_id"`
	Time    string        `bson:"time,omitempty"  form:"time" json:"time"`
}
