package io

type SdpRequest struct {
	App  string  `json:"app" form:"app"`
	Host string  `json:"host" form:"host"`
	Port int     `json:"port" form:"port"`
	Cpu  float32 `json:"cpu" form:"cpu"`
	Memo float32 `json:"memo" form:"memo"`
}
