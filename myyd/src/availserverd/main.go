/**
* @Author: cl
* @Date: 2021/1/16 11:11
 */
package main

import (
	"github.com/ChenLong-dev/gobase/mbase"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyd/src/availserverd/config"
	"os"
	"time"
)

const (
	LINK_DETECT = "link"
	LINE_DETECT = "line"
	ACL_DETECT  = "acl"
	LINK_ADDR   = "0.0.0.0:36167"
	LINE_ADDR   = "0.0.0.0:36168"
)

var (
	PointBound    int
	RegionBound   int
	DetectType    string
	DetectLimit   int
	GrpcAddr      string
	AclGrpcAddr   string
	AlarmTemplate string
	CheckIntvl    uint64
)

func main() {
	updateTaskIntvl := config.Conf.AvailConf.UpdateTaskIntvl
	CheckIntvl = config.Conf.AvailConf.CheckIntvl
	checkTimeout := config.Conf.AvailConf.CheckTimeout
	keepLastPointNum := config.Conf.AvailConf.KeepLastPointNum
	maxFlyCheckCount := config.Conf.AvailConf.MaxFlyCheckCount
	entryAddr := config.Conf.AvailConf.EntryAddr
	entryPassword := config.Conf.AvailConf.EntryPassword
	mdbgAddr := config.Conf.AvailConf.MdbgAddr

	mongoHost := config.Conf.Mongo.Host
	mongoUsername := config.Conf.Mongo.UserName
	mongoPassword := config.Conf.Mongo.Password
	mongoDbname := config.Conf.Mongo.DbName

	DetectType = config.Conf.AvailConf.DetectType
	PointBound = config.Conf.AvailConf.PointBound
	RegionBound = config.Conf.AvailConf.RegionBound
	GrpcAddr = config.Conf.AvailConf.GrpcAddr
	AclGrpcAddr = config.Conf.AvailConf.AclGrpcAddr
	DetectLimit = config.Conf.AvailConf.DetectLimit
	AlarmTemplate = config.Conf.AvailConf.AlarmTemplate
	AlarmAlias := config.Conf.AvailConf.AlarmAlias
	AlarmIntvl := config.Conf.AvailConf.AlarmIntvl

	mbase.Init()
	mlog.Infof("updateTaskIntvl:[%d], CheckIntvl:[%d], checkTimeout:[%d], keepLastPointNum:[%d], maxFlyCheckCount:[%d],"+
		" entryAddr:[%s], entryPassword:[%s], mdbgAddr:[%s], mongoHost:[%s], mongoUsername:[%s], mongoPassword:[%s], mongoDbname:[%s], "+
		"detectType:[%s], pointBound:[%d], regionBound:[%d], grpcAddr:[%s], detectLimit:[%d], AlarmAlias:[%+v], AlarmTemplate:[%s], "+
		"AlarmIntvl:[%d], aclgrpcAddr:[%s], mqurl:[%s]",
		updateTaskIntvl, CheckIntvl, checkTimeout, keepLastPointNum, maxFlyCheckCount, entryAddr, entryPassword,
		mdbgAddr, mongoHost, mongoUsername, mongoPassword, mongoDbname, DetectType, PointBound, RegionBound, GrpcAddr, DetectLimit,
		AlarmAlias, AlarmTemplate, AlarmIntvl, AclGrpcAddr, config.Conf.Mq.BrokerUrl)

	if err := initMongoDB(mongoHost, mongoUsername, mongoPassword, mongoDbname); err != nil {
		mlog.Errorf("initMongoDB fail, error:%v", err)
		os.Exit(1)
	}

	InitChecks()

	if err := InitEntry(entryAddr); err != nil {
		mlog.Errorf("InitEntry(%s, %s) error:%v", entryAddr, entryPassword)
		os.Exit(1)
	}

	InitResultsMan(keepLastPointNum)

	InitMdbg(mdbgAddr)

	InitMan(time.Duration(updateTaskIntvl)*time.Second, time.Duration(CheckIntvl)*time.Second, time.Duration(checkTimeout)*time.Second, maxFlyCheckCount)

	if DetectType == ACL_DETECT {
		if err := InitAclGrpcServer(); err != nil {
			mlog.Errorf("InitAclGrpcServer err:%v", err)
		}
	}

	if DetectType == LINK_DETECT {
		if err := InitLinkGrpcServer(LINK_ADDR); err != nil {
			mlog.Errorf("InitLinkGrpcServer err:%v", err)
		}
	}

	if DetectType == LINE_DETECT {
		if err := InitLinkGrpcServer(LINE_ADDR); err != nil {
			mlog.Errorf("InitLineGrpcServer err:%v", err)
		}
	}

	select {}
}

