/*
 * @Description:
 * @Author: Chen Long
 * @Date: 2020-12-16 12:01:13
 * @LastEditTime: 2020-12-16 12:02:27
 * @LastEditors: Chen Long
 * @Reference:
 */

package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"time"

	"config"
	"mlog"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

var (
	client  *clientv3.Client
	timeout = 5 * time.Second
	//DcSaltFlag = false
	//CloudEyeSeverFlag = false
	//SysLogLevelFlag = false
)

const (
	LoopholeInfo   = "loophole"
	TamperInfo     = "tamper"
	SaasConfig     = "saas_config"
	DcSalt         = "DC_SALT"
	CloudEyeServer = "CLOUD_EYE_SERVER"
	SysLogLevel    = "SYS_LOG_LEVEL"
)

func init() {
	mlog.Infof("etcd config info:%v", config.Conf.Etcd)
	var err error
	cert, err := tls.LoadX509KeyPair(config.Conf.Etcd.EtcdCert, config.Conf.Etcd.EtcdCertKey)
	if err != nil {
		mlog.Error(err)
		return
	}

	caData, e := ioutil.ReadFile(config.Conf.Etcd.EtcdCa)
	if e != nil {
		return
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}

	config := clientv3.Config{
		Endpoints:   config.Conf.Etcd.Hosts,
		DialTimeout: timeout,
		Username:    config.Conf.Etcd.UserName,
		Password:    config.Conf.Etcd.Password,
		TLS:         tlsConfig,
	}
	if client, err = clientv3.New(config); err != nil {
		panic(err)
	}
}

func Watch(scanType, key string) (map[string]string, error) {
	// 创建一个监听器
	var watchStartRevision int64
	watcher := clientv3.NewWatcher(client)
	watchRespChan := watcher.Watch(context.TODO(), key, clientv3.WithRev(watchStartRevision))

	// 用于读写etcd的键值对
	kv := clientv3.NewKV(client)

	// 读取入参key为前缀的所有key
	mapResp := make(map[string]string)
	if getResp, err := kv.Get(context.TODO(), key, clientv3.WithPrefix()); err != nil {
		return nil, err
	} else {
		// 获取成功
		mlog.Infof("respp kvs is %#v", getResp.Kvs)
		for _, resp := range getResp.Kvs {
			mapResp[string(resp.Key)] = string(resp.Value)
		}
		mlog.Infof("mapResp:%+v", mapResp)
	}

	// 处理kv变化事件
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				mlog.Infof("put data:%#v", string(event.Kv.Value))
				mapResp[string(event.Kv.Key)] = string(event.Kv.Value)
				switch string(event.Kv.Key) {
				case "/" + SaasConfig + "/" + DcSalt:
					getSaltInfo()
				case "/" + SaasConfig + "/" + CloudEyeServer:
					getSignTimeoutInfo()
				case "/" + SaasConfig + "/" + SysLogLevel:
					getSysLogLevelInfo(scanType)
				default:
					continue
				}
			case mvccpb.DELETE:
				mlog.Infof("delete key:%#v", string(event.Kv.Key))
				delete(mapResp, string(event.Kv.Key))
			}
		}
	}

	return mapResp, nil
}

func getSaltInfo() {
	var getSaltBytes []byte
	var err error
	if getSaltBytes, err = etcd.Get("/" + SaasConfig + "/" + DcSalt); err != nil {
		mlog.Error(err)
	}
	mlog.Debugf("saltBytes are %v\n", string(getSaltBytes))
	if err = json.Unmarshal(getSaltBytes, &config.Conf.SaltInfo); err != nil {
		mlog.Infof("salt json unmarshal fail, err:%v", err)
	}
	mlog.Debugf("SaltInfo is %#v\n", config.Conf.SaltInfo)
}

func getSignTimeoutInfo() {
	var getSignTimeoutBytes []byte
	var err error
	//获取etcd中sign_timeout
	if getSignTimeoutBytes, err = etcd.Get("/" + SaasConfig + "/" + CloudEyeServer); err != nil {
		mlog.Error(err)
	}
	mlog.Debugf("getSignTimeoutBytes are %v", string(getSignTimeoutBytes))

	//解析到SignTimeout结构
	if err = json.Unmarshal(getSignTimeoutBytes, &config.Conf.SignTimeout); err != nil {
		mlog.Infof("json unmarshal fail, err:%v", err)
	}
	mlog.Debugf("config.Conf.SignTimeout is %#v\n", config.Conf.SignTimeout)
}

func getSysLogLevelInfo(scanType string) {
	var getSysLogLevelBytes []byte
	var err error
	//获取etcd中tamper sys log
	if getSysLogLevelBytes, err = etcd.Get("/" + SaasConfig + "/" + SysLogLevel); err != nil {
		mlog.Error(err)
	}
	mlog.Debugf("getSysLogLevelBytes are %v", string(getSysLogLevelBytes))

	//解析到tamper sys log结构
	if err = json.Unmarshal(getSysLogLevelBytes, &config.Conf.SysLogLevel); err != nil {
		mlog.Infof("json unmarshal fail, err:%v", err)
	}
	mlog.Debugf("config.Conf.SysLogLevel is %#v\n", config.Conf.SysLogLevel)

	var SysLogLevelConfig string
	var level int
	switch scanType {
	case TamperInfo:
		SysLogLevelConfig = config.Conf.SysLogLevel.TamperLogLevel
	case LoopholeInfo:
		SysLogLevelConfig = config.Conf.SysLogLevel.LoopholeLogLevel
	default:
		SysLogLevelConfig = "info"
	}
	switch SysLogLevelConfig {
	case "all":
		level = 0
	case "trace":
		level = 1
	case "debug":
		level = 2
	case "info":
		level = 3
	case "warn":
		level = 4
	case "error":
		level = 5
	case "fatal":
		level = 6
	default:
		level = 3
	}
	mlog.Debugf("level value is %v", level)
	mlog.SetParams(config.LogConf.Log.Path, config.LogConf.Log.MaxSize, config.LogConf.Log.MaxBackups,
		config.LogConf.Log.MaxAge, level)
}
