package base

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"management/infra"
)

var session *mgo.Session

//dbx 数据库实例
var database *mgo.Database

func MgoDatabase() *mgo.Database {
	Check(database)
	return database
}

//dbx数据库starter，并且设置为全局
type MongoStarter struct {
	infra.BaseStarter
}

func (s *MongoStarter) Setup(ctx infra.StarterContext) {
	// 创建链接
	session, mgoError := mgo.Dial("")
	if mgoError != nil {
		logrus.Error("mongodb链接失败！")
		panic(mgoError)
	}
	// 选择DB
	database = session.DB("web")

}

func (s *MongoStarter) Stop(ctx infra.StarterContext) {
	session.Close()
}
