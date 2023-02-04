package sdp

import (
	"fmt"
	"github.com/saxon134/go-sdp/conf"
	"github.com/saxon134/go-sdp/db"
	"github.com/saxon134/go-sdp/db/models"
	"github.com/saxon134/go-sdp/io"
	"github.com/saxon134/go-utils/saData/saHit"
	"time"
)

const RedisAppKey = "%s:apps:%s"

type Config struct {
	Host   string `json:"h"`
	Port   int    `json:"p"`
	Weight int    `json:"w,omitempty"` //权重
	Time   int64  `json:"t"`           //上次ping的时间
}

var Chan chan io.SdpRequest

func init() {
	Chan = make(chan io.SdpRequest, 10)
	go appSdp()
}

// 处理的时候必须保证要能拿到Redis数据
// 所以使用channel保证一致性，如果是部署多个实例，因为并发可能性极低，暂不考虑
func appSdp() {
	for {
		if in, ok := <-Chan; ok {
			var key = saHit.Str(conf.Conf.Name != "", conf.Conf.Name, "sdp")
			key = fmt.Sprintf(RedisAppKey, key, in.App)

			var sdpAry = make([]*Config, 0, 10)
			_ = db.Redis.GetObj(key, &sdpAry)

			var existed *Config
			for _, c := range sdpAry {
				if c.Host == in.Host && c.Port == in.Port {
					existed = c
					break
				}
			}

			var now = time.Now().UnixMilli()

			//不存在是注册
			if existed == nil {
				var m = &Config{
					Host:   in.Host,
					Port:   in.Port,
					Weight: 100,
					Time:   now,
				}
				if in.Cpu > 0.9 || in.Memo > 0.9 {
					m.Weight = 1
				} else if in.Cpu > 0.8 || in.Memo > 0.8 {
					m.Weight = 10
				}
				sdpAry = append(sdpAry, m)
			} else
			//存在是ping
			{
				//超过2秒可能服务就慢了
				var delay = now - existed.Time
				if delay > 1000 {
					existed.Weight = 1
				} else if delay > 200 {
					existed.Weight = 10
				} else {
					existed.Weight = 100
				}

				if in.Cpu > 0.9 || in.Memo > 0.9 {
					existed.Weight = 1
				} else if in.Cpu > 0.8 || in.Memo > 0.8 {
					existed.Weight = 10
				}
			}
			_ = db.Redis.SetObj(key, sdpAry, time.Hour*48)

			//将服务器信息保存到资源
			{
				if db.MySql == nil {
					continue
				}

				var obj = new(models.TblResource)
				db.MySql.Table(models.TBNResource).
					Where("address = ?", in.Host).
					Where("`type` = 1").
					First(obj)
				obj.Type = 1
				obj.Status = saHit.Int(obj.Status != 0, obj.Status, 1)
				obj.Address = in.Host
				obj.Name = in.Host
				obj.Cpu = in.Cpu
				obj.Memo = in.Memo
				db.MySql.Save(obj)
			}
		}
	}
}