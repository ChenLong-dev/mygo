/**
* @Author: cl
* @Date: 2021/1/16 11:04
 */
package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Mongo     *Mongo
	Mq        *Mq
	AvailConf *AvailConf
}

type Mongo struct {
	Host     string
	UserName string
	Password string
	DbName   string
}

type Mq struct {
	TaskName     string
	QueueName    string
	RoutingKey   string
	ExchangeName string
	ExchangeType string
	BrokerUrl    string
}

type AvailConf struct {
	UpdateTaskIntvl  uint64
	CheckIntvl       uint64
	CheckTimeout     uint64
	KeepLastPointNum int
	MaxFlyCheckCount int64
	EntryAddr        string
	EntryPassword    string
	MdbgAddr         string
	DetectType       string
	PointBound       int
	RegionBound      int
	DetectLimit      int
	AlarmAlias       []string
	GrpcAddr         string
	AclGrpcAddr      string
	AlarmTemplate    string
	AlarmIntvl       int64
	TargetUrl        []string
}

var Conf Config

func init() {
	//读取本地配置获取服务IP、端口、监听url和ETCD配置信息
	if _, err := toml.DecodeFile("./conf/server.toml", &Conf); err != nil {
		panic(err)
	}
}