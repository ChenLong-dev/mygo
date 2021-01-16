/**
 .oooooo..o            .o.                  .o.             .oooooo..o
d8P'    `Y8           .888.                .888.           d8P'    `Y8
Y88bo.               .8"888.              .8"888.          Y88bo.
 `"Y8888o.          .8' `888.            .8' `888.          `"Y8888o.
     `"Y88b        .88ooo8888.          .88ooo8888.             `"Y88b
oo     .d8P       .8'     `888.        .8'     `888.       oo     .d8P
8""88888P'       o88o     o8888o      o88o     o8888o      8""88888P'
*/

package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/ChenLong-dev/gobase/mlog"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"myyy/src/availserverd/config"
	"time"
)

var (
	client      *clientv3.Client
	kv          clientv3.KV
	timeout     = 5 * time.Second
)

func NewClient() error {
	config := clientv3.Config{
		Endpoints:   config.Conf.Etcd.Hosts,
		DialTimeout: timeout,
		Username:    config.Conf.Etcd.UserName,
		Password:    config.Conf.Etcd.Password,
	}

	var err error
	if client, err = clientv3.New(config); err != nil {
		return err
	}

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	return nil
}

func NewTLSClient() error {
	mlog.Infof("etcd config info:%v", config.Conf.Etcd)
	cert, err := tls.LoadX509KeyPair(config.Conf.Etcd.EtcdCert, config.Conf.Etcd.EtcdCertKey)
	if err != nil {
		mlog.Error(err)
		return err
	}

	caData, e := ioutil.ReadFile(config.Conf.Etcd.EtcdCa)
	if e != nil {
		return e
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
		return err
	}

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	mlog.Infof("etcd NewTLSClient success")
	return nil
}

func Close() {
	client.Close()
}

func Put(key, val string) error {
	if _, err := kv.Put(context.TODO(), key, val); err != nil {
		return err
	}

	return nil
}

func Get(key string) ([]byte, error) {
	var getResp *clientv3.GetResponse
	var err error

	if getResp, err = kv.Get(context.TODO(), key); err != nil {
		return nil, err
	}

	for _, resp := range getResp.Kvs {
		mlog.Infof("getResp.Kvs key:%s, value:%s\n", string(resp.Key), string(resp.Value))
		if string(resp.Key) == key {
			return resp.Value, nil
		}
	}

	return nil, nil
}

func Delete(key string, opts ...clientv3.OpOption) error {
	if _, err := kv.Delete(context.TODO(), key, opts...); err != nil {
		return err
	}

	return nil
}

