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
	var sdpAry []*Config
	v, expired := saCache.MGet("discoveryApp:" + app)
	if v != nil {
		sdpAry, _ = v.([]*Config)
	}

	if expired || sdpAry == nil || len(sdpAry) == 0 {
		sdpAry = make([]*Config, 0, 10)
		saCache.MSetWithFunc("discoveryApp:"+app, time.Second*20, func() (interface{}, error) {
			var params = map[string]string{"app": app}
			if secret != "" {
				var timestamp = saData.I64tos(time.Now().Unix())
				params["timestamp"] = timestamp
				params["sign"] = saData.Md5(secret+timestamp, true)
			}

			res, err := saHttp.Get(saUrl.ConnectUri(address, "discovery"), params)
			if err != nil {
				saLog.Err("discovery error:", err)
				return nil, err
			}

			_ = saData.StrToModel(res, &sdpAry)
			return sdpAry, err
		})
	}

	if sdpAry != nil && len(sdpAry) > 0 {
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
