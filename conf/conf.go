package conf

import (
	"gopkg.in/yaml.v2"
	"os"
)

var Conf *ModelConf

type ModelConf struct {
	Name   string
	Secret string //如果配置了秘钥，所有接口都需要加密
	Port   string //http端口
	Path   string //对外提供注册的接口路径，不要包含域名&端口

	Redis struct {
		Host string
		Pass string
	}

	MySql struct {
		Host string
		Pass string
		User string
		Db   string
	}

	Feishu struct {
		Webhookurl string
		Secret     string
	}
}

func Init() *ModelConf {
	if Conf == nil {
		Conf = new(ModelConf)

		f_n := "./config.yaml"
		yamlData, err := os.ReadFile(f_n)
		if err != nil {
			panic("配置文件路径有误")
		}

		err = yaml.Unmarshal(yamlData, Conf)
		if err != nil {
			panic("配置文件信息有误")
		}
	}
	return Conf
}
