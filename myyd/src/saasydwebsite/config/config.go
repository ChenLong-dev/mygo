/**
* @Author: cl
* @Date: 2021/1/16 15:35
 */
package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Mongo     *Mongo
	AvailConf *AvailConf
}

type Mongo struct {
	Host     string
	UserName string
	Password string
	DbName   string
}

type AvailConf struct {
	ListenAddr string
}

var Conf Config

func init() {
	//读取本地配置获取可用性检测点配置信息
	if _, err := toml.DecodeFile("./conf/server.toml", &Conf); err != nil {
		panic(err)
	}
}

