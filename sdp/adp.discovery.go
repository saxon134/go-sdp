package sdp

import (
	"fmt"
	"github.com/saxon134/go-sdp/conf"
	"github.com/saxon134/go-sdp/db"
	"github.com/saxon134/go-utils/saCache"
	"github.com/saxon134/go-utils/saData/saHit"
	"math/rand"
	"time"
)

// Discovery 查找可用服务
func Discovery(app string) (host string, port int) {
	v := saCache.SMGet("discoveryApp:" + app)
	sdpAry := v.([]*Config)
	if sdpAry == nil || len(sdpAry) == 0 {
		sdpAry = make([]*Config, 0, 10)
		var key = fmt.Sprintf(RedisAppKey, saHit.Str(conf.Conf.Name != "", conf.Conf.Name, "sdp"), app)
		_ = db.Redis.GetObj(key, &sdpAry)
		if len(sdpAry) > 0 {
			saCache.SMSet("discoveryApp:"+app, sdpAry, time.Second*2)

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
	}
	return "", 0
}
