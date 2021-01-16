package main

import (
	"github.com/ChenLong-dev/gobase/mbase"
	"github.com/ChenLong-dev/gobase/mbase/msys"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyy/src/availd/config"
	"os"
)

const (
	LINK_DETECT = "link"
	LINE_DETECT = "line"
	ACL_DETECT  = "acl"
)

var (
	DetectType string
)

func main() {
	concurrent := config.Conf.AvailConf.Concurrent
	region := config.Conf.AvailConf.Region
	password := config.Conf.AvailConf.LoginPassword
	addr := config.Conf.AvailConf.LoginAddr
	DetectType = config.Conf.AvailConf.DetectType
	//mdbgAddr := config.Conf.AvailConf.MdbgAddr

	mbase.Init()

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

	//mcom.InitMdbg(mdbgAddr)

	select {}
}
