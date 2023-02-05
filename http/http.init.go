package http

import (
	"fmt"
	"github.com/saxon134/go-sdp/conf"
	"github.com/saxon134/go-sdp/db"
	"github.com/saxon134/go-sdp/db/models"
	"github.com/saxon134/go-sdp/sdp"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saUrl"
	"github.com/saxon134/go-utils/saLog"
	"net/http"
	"net/url"
)

func Init() {
	if conf.Conf.Port == "" {
		panic("http port can not be empty")
	}

	http.HandleFunc(saUrl.ConnPath(conf.Conf.Path, "register"), RegisterHandler)
	http.HandleFunc(saUrl.ConnPath(conf.Conf.Path, "ping"), PingHandler)
	http.HandleFunc(saUrl.ConnPath(conf.Conf.Path, "discovery"), DiscoveryHandler)

	http.HandleFunc(saUrl.ConnPath(conf.Conf.Path, "resource"), GetResourceHandler)

	saLog.Log("Http listening on " + conf.Conf.Port)
	err := http.ListenAndServe(":"+conf.Conf.Port, nil)
	if err != nil {
		panic("http err:" + err.Error())
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var query = r.URL.Query()
	if checkSign(query) == false {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("sign error"))
		return
	}

	var in = sdp.Request{}
	in.App = query.Get("app")
	in.Host = query.Get("host")
	in.Port = saData.Int(query.Get("port"))
	in.Cpu, _ = saData.ToFloat32(query.Get("cpu"))
	in.Memo, _ = saData.ToFloat32(query.Get("memo"))
	if in.App == "" || in.Host == "" || in.Port <= 0 {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("missing params"))
		return
	}

	//注册app服务
	sdp.Chan <- in
	w.WriteHeader(200)
	_, _ = w.Write([]byte(saData.String(map[string]interface{}{"code": 0})))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	RegisterHandler(w, r)
}

func DiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	var query = r.URL.Query()
	if checkSign(query) == false {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("sign error"))
		return
	}

	var app = query.Get("app")
	if app == "" {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("leak app"))
		return
	}

	var sdpAry = make([]*sdp.Config, 0, 10)
	var key = fmt.Sprintf(sdp.RedisAppKey, app)
	_ = db.Redis.GetObj(key, &sdpAry)
	if len(sdpAry) > 0 {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(saData.String(sdpAry)))
	}
}

func GetResourceHandler(w http.ResponseWriter, r *http.Request) {
	if db.MySql == nil {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(saData.String(map[string]interface{}{
			"code":   0,
			"result": []struct{}{},
		})))
		return
	}

	var query = r.URL.Query()
	if checkSign(query) == false {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("sign error"))
		return
	}

	var ary = make([]*models.TblResource, 0, 10)
	var conn = db.MySql.Table(models.TBNResource).Where("status > -1")
	t, _ := saData.ToInt(query.Get("type"))
	if t > 0 {
		conn.Where("`type` = ?", t)
	}
	conn.Find(&ary)
	w.WriteHeader(200)
	_, _ = w.Write([]byte(saData.String(map[string]interface{}{
		"code":   0,
		"result": ary,
	})))
}

func checkSign(query url.Values) bool {
	if conf.Conf.Secret == "" {
		return true
	}

	var sign = query.Get("sign")
	var timestamp = query.Get("timestamp")
	if sign == "" || timestamp == "" {
		return false
	}

	sign2 := saData.Md5(sign+timestamp, true)
	return sign == sign2
}
