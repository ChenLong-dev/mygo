package main

import (
	"context"
	"fmt"
	"github.com/ChenLong-dev/gobase/mlog"
	"google.golang.org/grpc"
	"math"
	"myyy/src/availserverd/etcd/register"

	//"myyy/src/availserverd/etcd"
	//"myyy/src/availserverd/etcd/register"
	"github.com/ChenLong-dev/gobase/etcd"
	//"github.com/ChenLong-dev/gobase/register"
	busiman "myyy/src/availserverd/saasyd_usability_business_manager_grpc"
	"myyy/src/scom"
	"net"
	"strconv"
	"time"
)

const (
	FREQ_MIN_LIMIT          = 1 //默认分钟为单位
	RegClusterName          = "services"
	UsabilityManServicePath = "saasdc_usability_manager/grpc"
	UsabilityManServicePort = 36768
)

type server struct{}

func detectRegister(key string) {
	for ;; {
		time.Sleep(time.Second * 10)
		key1 := "/"+ RegClusterName + "/" + UsabilityManServicePath + "/" + key
		resp, err := etcd.Get(key1)
		mlog.Debugf("get detect register xxx 1 [resp:%s] [err:%v] [key1:%s]", string(resp), err, key1)
		if resp != nil && string(resp) != "" {
			continue
		}
		_, _ = InitRegSevice()
	}
}

func InitRegSevice() (string, error) {
	ip, e := scom.GetLocalIPv4()
	if e != nil {
		return "", e
	}
	host := fmt.Sprintf("%s:%d", ip, UsabilityManServicePort)
	mlog.Debugf("InitRegSevice host:%s", host)
	key := ip+":"+strconv.Itoa(UsabilityManServicePort)
	for ;; {
		Reg := register.EtcdRegister{Cluster: RegClusterName}
		if err := Reg.Register(UsabilityManServicePath, key, key); err != nil {
			mlog.Error(err)
			time.Sleep(time.Second * 1)
			continue
		}
		break
	}

	return key, nil
}

func (s *server) BusinessTaskAdd(ctx context.Context, in *busiman.BusinessTaskAddReq) (*busiman.BusinessTaskAddRsp, error) {
	mlog.Debugf("BusinessTaskAdd enter")
	if in.GetUrl() == "" {
		return &busiman.BusinessTaskAddRsp{
			Code:    -1,
			Message: "business url is empty",
		}, nil
	}
	if in.GetFreq() < FREQ_MIN_LIMIT {
		return &busiman.BusinessTaskAddRsp{
			Code:    -1,
			Message: "freq less than min-limit",
		}, nil
	}

	var eyeInfo scom.EyeInfo
	eyeInfo.Url = in.GetUrl()
	eyeInfo.Freq = in.GetFreq()
	if err := defaultMan.AddEyeInfosCache(eyeInfo); err != nil {
		mlog.Warnf("BusinessTaskAdd [%v]\n", err)
		return &busiman.BusinessTaskAddRsp{
			Code:    -1,
			Message: fmt.Sprintf("%v", err),
		}, nil
	}

	return &busiman.BusinessTaskAddRsp{
		Code:    0,
		Message: "success",
	}, nil
}

func (s *server) BusinessTaskDel(ctx context.Context, in *busiman.BusinessTaskDelReq) (*busiman.BusinessTaskDelRsp, error) {
	mlog.Debugf("BusinessTaskDel enter")
	if in.GetUrl() == "" {
		return &busiman.BusinessTaskDelRsp{
			Code:    -1,
			Message: "business url is empty",
		}, nil
	}

	if err := defaultMan.DelEyeInfosCache(in.GetUrl()); err != nil {
		mlog.Warnf("BusinessTaskDel [%v]\n", err)
		return &busiman.BusinessTaskDelRsp{
			Code:    -1,
			Message: fmt.Sprintf("%v", err),
		}, nil
	}

	return &busiman.BusinessTaskDelRsp{
		Code:    0,
		Message: "success",
	}, nil
}

func (s *server) BusinessTaskEdit(ctx context.Context, in *busiman.BusinessTaskEditReq) (*busiman.BusinessTaskEditRsp, error) {
	mlog.Debugf("BusinessTaskEdit enter")
	if in.GetUrl() == "" {
		return &busiman.BusinessTaskEditRsp{
			Code:    -1,
			Message: "business url is empty",
		}, nil
	}

	if in.GetFreq() < FREQ_MIN_LIMIT {
		return &busiman.BusinessTaskEditRsp{
			Code:    -1,
			Message: "freq less than min-limit",
		}, nil
	}

	var eyeInfo scom.EyeInfo
	eyeInfo.Url = in.GetUrl()
	eyeInfo.Freq = in.GetFreq()
	if err := defaultMan.EditEyeInfosCache(eyeInfo); err != nil {
		mlog.Warnf("BusinessTaskEdit [%v]\n", err)
		return &busiman.BusinessTaskEditRsp{
			Code:    -1,
			Message: fmt.Sprintf("%v", err),
		}, nil
	}

	return &busiman.BusinessTaskEditRsp{
		Code:    0,
		Message: "success",
	}, nil
}

func InitBusiManServer() error {
	addr := fmt.Sprintf("0.0.0.0:%d", UsabilityManServicePort)
	mlog.Infof("InitBusiManServer listenAddr:%s", addr)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		mlog.Infof("listen is failed: %v\n", err)
		return err
	}
	var options = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(1073741824),
	}
	s := grpc.NewServer(options...)

	busiman.RegisterBusinessTaskManagerServer(s, &server{})

	var key string
	var err1 error
	if key, err1 = InitRegSevice(); err != nil {
		mlog.Errorf("register saasdc_usability_manager grpc serevice fail, err:%v, key:%s", err1, key)
		return err1
	}
	mlog.Debugf("xxx 1 [%s]", RegClusterName + "/" + UsabilityManServicePath + "/" + key)
	go detectRegister(key)

	if err := s.Serve(listen); err != nil {
		mlog.Fatalf("failed to serve: %v", err)
	}

	return nil
}

