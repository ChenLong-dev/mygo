/**
* @Author: cl
* @Date: 2021/1/14 19:42
 */
package main

import (
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase"
	"github.com/ChenLong-dev/gobase/mbase/mcom"
	"github.com/ChenLong-dev/gobase/mbase/msys"
	"github.com/ChenLong-dev/gobase/mlog"
	conf "myyd/src/availd/config"
	"os"
)

const (
	LINK_DETECT = "link"
	LINE_DETECT = "line"
	ACL_DETECT  = "acl"
)

var (
	DetectType   string
	PingInterval int
)

func main() {
	//concurrent := flag.Int("concurrent", 1000, "Concurrent check")
	//region := flag.String("region", "cloud_node_gz", "deploy place")
	//password := flag.String("loginPassword", "123456789", "login password")
	//addr := flag.String("loginAddr", "121.46.4.209:36161", "login addr")
	//DetectType = flag.String("detectType", "acl", "detect line")

	concurrent := conf.Conf.AvailConf.Concurrent
	region := conf.Conf.AvailConf.Region
	password := conf.Conf.AvailConf.LoginPassword
	addr := conf.Conf.AvailConf.LoginAddr
	DetectType = conf.Conf.AvailConf.DetectType
	mdbgAddr := conf.Conf.AvailConf.MdbgAddr
	pi := conf.Conf.AvailConf.PingInterval
	if pi <= 0 {
		PingInterval = 20
	} else {
		PingInterval = pi
	}
	fmt.Printf("concurrent:[%d], region:[%s], password:[%s], addr:[%s], detectType:[%s], mdbgAddr:[%s] PingInterval:[%d]",
		concurrent, region, password, addr, DetectType, mdbgAddr, PingInterval)

	mbase.Init()

	//var mdbgAddr string
	//if DetectType == "acl" {
	//	mdbgAddr = "127.0.0.1:16161"
	//} else if DetectType == "link" {
	//	mdbgAddr = "127.0.0.1:16162"
	//} else if DetectType == "line" {
	//	mdbgAddr = "127.0.0.1:16163"
	//} else {
	//	mlog.Errorf("not detect type [type:%s]\n", *DetectType)
	//	os.Exit(1)
	//}

	mlog.Infof("concurrent:[%d], region:[%s], password:[%s], addr:[%s], detectType:[%s], mdbgAddr:[%s] PingInterval:[%d]",
		concurrent, region, password, addr, DetectType, mdbgAddr, PingInterval)

	msys.SetMaxFdSize(100000)

	err := Init(concurrent)
	if err != nil {
		mlog.Errorf("init error:%v", err)
		os.Exit(1)
	}

	err = ConnectServer(addr, region, password)
	if err != nil {
		mlog.Errorf("ConnectServer error:%v", err)
		os.Exit(1)
	}

	mcom.InitMdbg(mdbgAddr)

	select {}
}

