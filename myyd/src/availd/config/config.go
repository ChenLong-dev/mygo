/**
* @Author: cl
* @Date: 2021/1/14 19:37
 */
package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	AvailConf *AvailConf
}

type AvailConf struct {
	Region        string
	LoginAddr     string
	LoginPassword string
	Concurrent    int
	DetectType    string
	MdbgAddr      string
	PingInterval  int
}

var Conf Config

func init() {
	//读取本地配置获取可用性检测点配置信息
	if _, err := toml.DecodeFile("./conf/server.toml", &Conf); err != nil {
		panic(err)
	}
}
