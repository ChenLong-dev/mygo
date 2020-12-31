/**
* @Author: cl
* @Date: 2020/12/31 16:56
 */
package config

var Cfg ConfigExt

//属性设计为指针类型
type ConfigExt struct {
	AppMode   string
	Mdbg      *MdbgCfg
	PopServer *PopServer
	Etcd      *Etcd
	Mongo     *Mongo
	Log       *Log
}

type MdbgCfg struct {
	Host   string
	Port   int
	Enable int
}

type PopServer struct {
	Host        string
	Port        int
	SrvLockKey  string
	EtcdCluster string
	ProjectName string
	CronExpres  string
	LimitNum    int
}

type Etcd struct {
	Hosts       []string
	UserName    string
	Salt        string
	EtcdCert    string
	EtcdCertKey string
	EtcdCa      string
}

type Mongo struct {
	Host     string
	UserName string
	Password string
	DbName   string
	AuthDB   string
	MaxConns uint64
	Port     int
	Key      string
}

type Log struct {
	Path       string
	Level      int
	MaxSize    int
	MaxBackups int
	MaxAge     int
}
