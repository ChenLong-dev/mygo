package main

import (
	"github.com/ChenLong-dev/gobase/mbase"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyy/src/availserverd/config"
	//"myyy/src/availserverd/etcd"
	"github.com/ChenLong-dev/gobase/etcd"
	"myyy/src/availserverd/etcd/discovery"
	"os"
	"time"
)

var (
	PointBound         int
	RegionBound        int
	DetectLimit        int
	GrpcAddr           string
	GetEyeinfoGrpcAddr string
	UploadLogGrpcAddr  string
	Mng                discovery.EtcdDis
)

const (
	ClusterName         = "services"
	BusiTaskServicePath = "saasdc_usability_c/grpc"
	UpLogServicePath    = "saasdc_usability_l/grpc"
)

func InitWatchSevice() {
	etcd.NewTLSClient()
	Mng = discovery.EtcdDis{Cluster: ClusterName}
	Mng.Watch(BusiTaskServicePath)
	Mng.Watch(UpLogServicePath)
	time.Sleep(time.Second * 5)
}

func main() {
	updateTaskIntvl := config.Conf.AvailConf.UpdateTaskIntvl
	checkIntvl := config.Conf.AvailConf.CheckIntvl
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

	PointBound = config.Conf.AvailConf.PointBound
	RegionBound = config.Conf.AvailConf.RegionBound
	GrpcAddr = config.Conf.AvailConf.GrpcAddr
	GetEyeinfoGrpcAddr = config.Conf.AvailConf.GetEyeinfoGrpcAddr
	UploadLogGrpcAddr = config.Conf.AvailConf.UploadLogGrpcAddr

	DetectLimit = config.Conf.AvailConf.DetectLimit
	AlarmAlias := config.Conf.AvailConf.AlarmAlias
	AlarmIntvl := config.Conf.AvailConf.AlarmIntvl

	mbase.Init()
	mlog.Infof("updateTaskIntvl:[%d], checkIntvl:[%d], checkTimeout:[%d], keepLastPointNum:[%d], maxFlyCheckCount:[%d],"+
		" entryAddr:[%s], entryPassword:[%s], mdbgAddr:[%s], mongoHost:[%s], mongoUsername:[%s], mongoPassword:[%s], mongoDbname:[%s], "+
		"pointBound:[%d], regionBound:[%d], grpcAddr:[%s], getEyeinfoGrpcAddr:[%s], detectLimit:[%d], AlarmAlias:[%+v], "+
		"AlarmIntvl:[%d], UploadLogGrpcAddr:[%s]",
		updateTaskIntvl, checkIntvl, checkTimeout, keepLastPointNum, maxFlyCheckCount, entryAddr, entryPassword,
		mdbgAddr, mongoHost, mongoUsername, mongoPassword, mongoDbname, PointBound, RegionBound, GrpcAddr, GetEyeinfoGrpcAddr, DetectLimit,
		AlarmAlias, AlarmIntvl, UploadLogGrpcAddr)

	InitWatchSevice()

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

	InitMan(time.Duration(updateTaskIntvl)*time.Second, time.Duration(checkIntvl)*time.Second, time.Duration(checkTimeout)*time.Second, maxFlyCheckCount)

	if err := InitBusiManServer(); err != nil {
		mlog.Errorf("InitBusiManServer err:%v", err)
	}

	select {}
}
