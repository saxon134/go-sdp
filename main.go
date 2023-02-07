package main

import (
	"github.com/saxon134/go-sdp/conf"
	"github.com/saxon134/go-sdp/db"
	"github.com/saxon134/go-sdp/http"
	"github.com/saxon134/go-sdp/sdp"
	"github.com/saxon134/go-utils/saLog"
)

func main() {
	//初始化
	conf.Init()

	//初始化日志
	saLog.Init(saLog.WarnLevel, saLog.ZapType)

	//初始化数据库
	db.Init()

	//初始化http服务
	go http.Init()

	//初始化SDP
	sdp.Init()

	//防止应用退出
	<-make(chan bool)
}
