package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/saxon134/go-utils/saHttp"
	"github.com/saxon134/go-utils/saLog"
	"time"
)

type FeiShu struct {
	webhookUrl string
	secret     string
}

func New(webhookUrl string, secret string) *FeiShu {
	return &FeiShu{webhookUrl: webhookUrl, secret: secret}
}

func (m *FeiShu) SendTxt(txt string) {
	var args = map[string]interface{}{
		"msg_type": "text",
		"content":  map[string]string{"text": txt},
	}
	if m.secret != "" {
		var timestamp = time.Now().Unix()
		var sign = sign(timestamp, m.secret)
		args["timestamp"] = timestamp
		args["sign"] = sign
	}
	_, _, err := saHttp.PostJson(m.webhookUrl, args)
	if err != nil {
		saLog.Err("飞书通知发送失败：", err.Error())
	}
}

func sign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
