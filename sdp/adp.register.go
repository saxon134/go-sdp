package sdp

import (
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saUrl"
	"github.com/saxon134/go-utils/saHttp"
	"github.com/saxon134/go-utils/saLog"
	"time"
)

var pingParams map[string]string

// Register 注册服务
func Register(address string, app string, host string, port int, secret string) {
	if host == "" || port <= 0 {
		saLog.Err("RPC register error: leak params")
		return
	}

	var params = map[string]string{"app": app, "host": host, "port": saData.Itos(port)}
	if secret != "" {
		var timestamp = saData.I64tos(time.Now().Unix())
		params["timestamp"] = timestamp
		params["sign"] = saData.Md5(secret+timestamp, true)
	}
	_, err := saHttp.Get(saUrl.ConnectUri(address, "register"), params)
	if err != nil {
		saLog.Err(saError.StackError(err))
		return
	}

	pingParams = params
	pingParams["address"] = address
	pingParams["secret"] = secret
	for {
		time.Sleep(time.Second * 2)
		go ping()
	}
}

func ping() {
	defer func() {
		if err := recover(); err != nil {
			saLog.Err("Sdp ping panic:", err)
		}
	}()

	var params = map[string]string{"app": pingParams["app"], "host": pingParams["host"], "port": pingParams["port"]}
	var secret = pingParams["secret"]
	if secret != "" {
		var timestamp = saData.I64tos(time.Now().Unix())
		params["timestamp"] = timestamp
		params["sign"] = saData.Md5(secret+timestamp, true)
	}

	_, err := saHttp.Get(saUrl.ConnectUri(pingParams["address"], "ping"), params)
	if err != nil {
		saLog.Err("Sdp ping error:", err)
		return
	}
}
