/**
* @Author: cl
* @Date: 2021/1/16 16:40
 */
package discovery

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/ChenLong-dev/gobase/mlog"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"math/rand"
	"myyy/src/availserverd/config"
	"strings"
	"time"
)

// Etcd注册的节点，一个节点代表一个client
type EtcdNode struct {
	State   bool
	Cluster string // 集群
	Key     string // key
	Info    string // 节点信息
}

type EtcdMaster struct {
	Cluster string // 集群
	Path    string // 路径
	Nodes   map[string]*EtcdNode
	Client  *clientv3.Client
}

type EtcdDis struct {
	Cluster string
	// 监听相关的服务
	MapWatch map[string]*EtcdMaster
}

func (d *EtcdDis) Watch(service string) error {
	var w *EtcdMaster
	var e error
	if w, e = NewMaster(d.Cluster, service); e != nil {
		mlog.Errorf("Watch Service:%s Failed!Error:%s", service, e.Error())
		return e
	}

	if d.MapWatch == nil {
		d.MapWatch = make(map[string]*EtcdMaster)
	}

	if _, ok := d.MapWatch[service]; ok {
		mlog.Infof("Service:%s Have Watch!\n", service)
		return nil
	}

	mlog.Debugf("w is %v", w.Nodes)
	d.MapWatch[service] = w

	return nil
}

// 获取服务的节点信息-随机获取
func (d *EtcdDis) GetServiceInfoRandom(service string) (EtcdNode, bool) {
	if d.MapWatch == nil {
		mlog.Warn("MapWatch is nil")
		return EtcdNode{}, false
	}

	if v, ok := d.MapWatch[service]; ok {
		mlog.Debugf("V is %v, ok is %v\n", v, ok)
		if v != nil {
			if n, ok1 := v.GetNodeRandom(); ok1 {
				mlog.Debugf("n is %v, ok1 is %v\n", v, ok1)
				return n, true
			}
			mlog.Debugf("n is %v\n", v)
		}
	} else {
		mlog.Errorf("Service:%s Not Be Watched!\n", service)
	}

	return EtcdNode{}, false
}

// 获取服务的节点信息-全部获取
func (d *EtcdDis) GetServiceInfoAllNode(service string) ([]EtcdNode, bool) {
	if d.MapWatch == nil {
		mlog.Warn("MapWatch is nil")
		return []EtcdNode{}, false
	}

	if v, ok := d.MapWatch[service]; ok {
		if v != nil {
			return v.GetAllNodes(), true
		}
	} else {
		mlog.Errorf("Service:%s Not Be Watched!\n", service)
	}

	return []EtcdNode{}, false
}

// 拆分service name、key；返回bool true表示成功;false表示失败
func SplitServiceNameKey(dir string) (string, string, bool) {
	if idx := strings.Index(dir, "/"); -1 != idx {
		name := dir[:idx]
		key := dir[idx+1:]
		return name, key, true
	}
	return "", "", false
}

func NewMaster(cluster string, watchPath string) (*EtcdMaster, error) {
	cert, err := tls.LoadX509KeyPair(config.Conf.Etcd.EtcdCert, config.Conf.Etcd.EtcdCertKey)
	if err != nil {
		mlog.Error(err)
		return nil, err
	}
	caData, e := ioutil.ReadFile(config.Conf.Etcd.EtcdCa)
	if e != nil {
		return nil,e
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Conf.Etcd.Hosts,
		DialTimeout: time.Second,
		Username:    config.Conf.Etcd.UserName,
		Password:    config.Conf.Etcd.Password,
		TLS:         tlsConfig,
	})


	if nil != err {
		mlog.Error(err)
		return nil, err
	}

	master := &EtcdMaster{
		Cluster: cluster,
		Path:    watchPath,
		Nodes:   make(map[string]*EtcdNode),
		Client:  cli,
	}

	// 监听观察节点
	go master.WatchNodes()

	mlog.Debugf("master is %v; master.nodes is %v", master, master.Nodes)
	return master, err
}

func NewEtcdNode(ev *clientv3.Event) string {
	return string(ev.Kv.Value)
}

// 监听观察节点
func (m *EtcdMaster) WatchNodes() {
	// 查看之前存在的节点
	mlog.Debugf("grpc's etcd addr:%v", "/"+m.Cluster+"/"+m.Path)
	resp, err := m.Client.Get(context.Background(), "/"+m.Cluster+"/"+m.Path, clientv3.WithPrefix())
	mlog.Debugf("etcd's responese to grpc:%v", resp)
	if nil != err {
		mlog.Error(err)
	} else {
		mlog.Debugf("resp.Kvs is :%v", resp.Kvs)
		for _, ev := range resp.Kvs {
			mlog.Infof("add dir:%q, value:%q\n", ev.Key, ev.Value)
			m.addNode(string(ev.Key), string(ev.Value))
			mlog.Debugf("m.nodes are %v", m.Nodes)
		}
	}

	rch := m.Client.Watch(context.Background(), "/"+m.Cluster+"/"+m.Path, clientv3.WithPrefix(), clientv3.WithPrevKV())
	mlog.Debugf("etcd's rch to grpc:%T %v", rch, rch)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			mlog.Debugf("etcd's ev to grpc:%T %v", ev, ev)
			switch ev.Type {
			case clientv3.EventTypePut:
				mlog.Infof("[%s] dir:%q, value:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := NewEtcdNode(ev)
				m.addNode(string(ev.Kv.Key), info)
			case clientv3.EventTypeDelete:
				mlog.Infof("[%s] dir:%q, value:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				k := ev.Kv.Key
				if len(ev.Kv.Key) > (len(m.Cluster) + 1) {
					k = ev.Kv.Key[len(m.Cluster)+1:]
				}
				delete(m.Nodes, string(k))
			default:
				mlog.Infof("[%s] dir:%q, value:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
	mlog.Debugf("m.nodes are %v", m.Nodes)
}

// 添加节点
func (m *EtcdMaster) addNode(key, info string) {
	k := key
	if len(key) > (len(m.Cluster) + 1) {
		k = key[len(m.Cluster)+1:]
	}

	node := &EtcdNode{
		State:   true,
		Cluster: m.Cluster,
		Key:     k,
		Info:    info,
	}

	m.Nodes[node.Key] = node
}

// 获取该集群下所有的节点
func (m *EtcdMaster) GetAllNodes() []EtcdNode {
	var temp []EtcdNode
	for _, v := range m.Nodes {
		if nil != v {
			temp = append(temp, *v)
		}
	}
	return temp
}

func (m *EtcdMaster) GetNodeRandom() (EtcdNode, bool) {
	count := len(m.Nodes)
	mlog.Debugf("count is %v:\n", count)
	// 该集群不存在节点时，直接返回false
	if count == 0 {
		return EtcdNode{}, false
	}
	idx := rand.Intn(count)
	for _, v := range m.Nodes {
		if idx == 0 {
			return *v, true
		}
		idx = idx - 1
	}
	return EtcdNode{}, false
}