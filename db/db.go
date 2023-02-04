package db

import (
	"github.com/saxon134/go-sdp/conf"
	"github.com/saxon134/go-utils/saLog"
	"github.com/saxon134/go-utils/saOrm"
	"github.com/saxon134/go-utils/saRedis"
)

var Redis *saRedis.Redis
var MySql *saOrm.DB

func Init() {
	//初始化MySql数据库
	{
		var config = conf.Conf.MySql
		if config.User != "" && config.Db != "" {
			dsn := config.User + ":" + config.Pass + "@tcp(" + config.Host + ")/ads" + "?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local"
			MySql = saOrm.Open(dsn, saOrm.Conf{})
			saLog.Log("MySql init ok")
		}
	}

	//初始化Redis
	{
		var err error
		var config = conf.Conf.Redis
		Redis, err = saRedis.Init(config.Host, config.Pass, 0)
		if err != nil {
			panic("Redis初始化出错" + err.Error())
		}

		_, err = Redis.Pool.Dial()
		if err != nil {
			panic("Redis初始化出错" + err.Error())
		}

		saLog.Log("Redis init ok")
	}
}
