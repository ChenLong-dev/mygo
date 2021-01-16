/**
* @Author: cl
* @Date: 2021/1/16 16:41
 */
package register
import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/ChenLong-dev/gobase/mlog"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"myyy/src/availserverd/config"
	"time"
)

type EtcdService struct {
	Cluster string // 集群名称
	Name    string // 服务名称
	Info    string // 节点信息
	stop    chan error
	leaseid clientv3.LeaseID
	client  *clientv3.Client
}

type EtcdRegister struct {
	Cluster     string
	MapRegister map[string]*EtcdService
}

// 注册
func (reg *EtcdRegister) Register(service, key, info string) error {
	var s *EtcdService
	var err error
	name := service + "/" + key
	if s, err = NewTlsService(reg.Cluster, name, info); err != nil {
		mlog.Errorf("Register service:%s error:%v", service, err.Error())
		return err
	}

	if reg.MapRegister == nil {
		reg.MapRegister = make(map[string]*EtcdService)
	}

	if _, ok := reg.MapRegister[name]; ok {
		mlog.Errorf("Service:%s Have Registered!", name)
		return nil
	}

	reg.MapRegister[name] = s
	// w维持心跳
	s.Start()

	return nil
}

// 更新
func (reg *EtcdRegister) UpdateInfo(service, key, info string) {
	name := service + "/" + key
	if s, ok := reg.MapRegister[name]; ok {
		if s != nil {
			s.SetValue(info)
		}
	}
}

// 注册ETCD服务
func NewService(cluster, name, info string) (*EtcdService, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Conf.Etcd.Hosts,
		DialTimeout: 5 * time.Second,
		Username:    config.Conf.Etcd.UserName,
		Password:    config.Conf.Etcd.Password,
	})

	if err != nil {
		mlog.Error(err)
		return nil, err
	}
	// 返回服务对象
	return &EtcdService{
		Cluster: cluster,
		Name:    name,
		Info:    info,
		stop:    make(chan error),
		client:  cli,
	}, err
}

// 注册ETCD服务
func NewTlsService(cluster, name, info string) (*EtcdService, error) {
	cert, err := tls.LoadX509KeyPair(config.Conf.Etcd.EtcdCert, config.Conf.Etcd.EtcdCertKey)
	if err != nil {
		mlog.Error(err)
		return nil, err
	}
	caData, e := ioutil.ReadFile(config.Conf.Etcd.EtcdCa)
	if e != nil {
		return nil, e
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	cli, tlsErr := clientv3.New(clientv3.Config{
		Endpoints:   config.Conf.Etcd.Hosts,
		DialTimeout: 5 * time.Second,
		Username:    config.Conf.Etcd.UserName,
		Password:    config.Conf.Etcd.Password,
		TLS:         tlsConfig,
	})

	if tlsErr != nil {
		mlog.Error(tlsErr)
		return nil, tlsErr
	}

	// 返回服务对象
	return &EtcdService{
		Cluster: cluster,
		Name:    name,
		Info:    info,
		stop:    make(chan error),
		client:  cli,
	}, tlsErr
}// 启动
func (s *EtcdService) Start() error {
	// 获取心跳的通道
	ch, err := s.keepLive()
	if err != nil {
		mlog.Error(err)
		return err
	}
	go func() {
		// 死循环
		for {
			select {
			case <-s.stop:
				s.revoke()
				return
			case <-s.client.Ctx().Done():
				mlog.Info("server closed")
				return
			case ka, ok := <-ch:
				if !ok {
					mlog.Info("keep live channel closed")
					s.revoke()
					return
				} else {
					mlog.Infof("recv reply frem service:[%s], ttl:[%d]", s.Name, ka.TTL)
				}
			}
		}
	}()
	return nil
}

// 保持心跳
func (s *EtcdService) keepLive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	key := "/" + s.Cluster + "/" + s.Name
	value, _ := json.Marshal(s.Info)

	// minimum lease TTL is 60-second
	resp, err := s.client.Grant(context.TODO(), 60)
	if err != nil {
		mlog.Error(err)
		return nil, err
	}

	mlog.Infof("Register, Key:%s, Value:%s\n", key, string(value))
	_, err = s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if nil != err {
		mlog.Error(err)
		return nil, err
	}

	s.leaseid = resp.ID
	return s.client.KeepAlive(context.TODO(), resp.ID)
}

// 设置节点信息
func (s *EtcdService) SetValue(info string) error {
	s.Info = info
	tmp, _ := json.Marshal(info)
	key := "/" + s.Cluster + "/" + s.Name
	if _, err := s.client.Put(context.TODO(), key, string(tmp), clientv3.WithLease(s.leaseid)); err != nil {
		mlog.Errorf("etcd set value failed! key:%s;value:%s", key, info)
		return err
	}

	return nil
}

// 停止
func (s *EtcdService) Stop() {
	s.stop <- nil
}

// 撤销
func (s *EtcdService) revoke() error {
	_, err := s.client.Revoke(context.TODO(), s.leaseid)
	if err != nil {
		mlog.Error(err)
	}

	mlog.Infof("service:%s stop\n", s.Name)
	return nil
}



