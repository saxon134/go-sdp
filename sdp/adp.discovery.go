package sdp

import (
	"github.com/saxon134/go-utils/saCache"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saUrl"
	"github.com/saxon134/go-utils/saHttp"
	"github.com/saxon134/go-utils/saLog"
	"math/rand"
	"time"
)

// Discovery 查找可用服务
func Discovery(address string, app string, secret string) (host string, port int) {
	sdpAry, _ := saCache.SMGet("discoveryApp:" + app).([]*Config)
	if sdpAry == nil || len(sdpAry) == 0 {
		sdpAry = make([]*Config, 0, 10)

		var params = map[string]string{"app": app}
		if secret != "" {
			var timestamp = saData.I64tos(time.Now().Unix())
			params["timestamp"] = timestamp
			params["sign"] = saData.Md5(secret+timestamp, true)
		}

		res, err := saHttp.Get(saUrl.ConnectUri(address, "discovery"), params)
		if err != nil {
			saLog.Err("discovery error:", err)
			return
		}

		_ = saData.StrToModel(res, &sdpAry)
		if len(sdpAry) > 0 {
			saCache.SMSet("discoveryApp:"+app, sdpAry, time.Second*20)
		}
	}

	if len(sdpAry) > 0 {
		var weight = 0
		for _, v := range sdpAry {
			weight += v.Weight
		}

		var r = rand.Intn(weight)
		weight = 0
		for _, v := range sdpAry {
			if weight >= r {
				return v.Host, v.Port
			}
			weight += v.Weight
		}
	}
	return "", 0
}
