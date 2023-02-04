package sdp

import (
	"fmt"
	"github.com/saxon134/go-sdp/conf"
	"github.com/saxon134/go-sdp/db"
	"github.com/saxon134/go-sdp/http"
	"github.com/saxon134/go-utils/saCache"
	"github.com/saxon134/go-utils/saData/saHit"
	"math/rand"
	"time"
)

func Discovery(app string) (host string, port int) {
	v := saCache.SMGet("discoveryApp:" + app)
	sdpAry := v.([]*http.SdpConfig)
	if sdpAry == nil || len(sdpAry) == 0 {
		sdpAry = make([]*http.SdpConfig, 0, 10)
		var key = fmt.Sprintf(http.SdpAppKey, saHit.Str(conf.Conf.Name != "", conf.Conf.Name, "sdp"), app)
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
