package sdp

import (
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saUrl"
	"github.com/saxon134/go-utils/saHttp"
	"github.com/saxon134/go-utils/saLog"
	"time"
)

// Register 注册服务
func Register(address string, host string, port int, secret string) {
	var params = map[string]string{"host": host, "port": saData.Itos(port)}
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

	for {
		go ping(address, host, port, secret)
		time.Sleep(time.Second * 2)
	}
}

func ping(address string, host string, port int, secret string) {
	defer func() {
		if err := recover(); err != nil {
			saLog.Err("Sdp ping panic:", err)
		}
	}()

	var params = map[string]string{"host": host, "port": saData.Itos(port)}
	if secret != "" {
		var timestamp = saData.I64tos(time.Now().Unix())
		params["timestamp"] = timestamp
		params["sign"] = saData.Md5(secret+timestamp, true)
	}
	_, err := saHttp.Get(saUrl.ConnectUri(address, "ping"), params)
	if err != nil {
		saLog.Err("Sdp ping error:", err)
		return
	}
}
