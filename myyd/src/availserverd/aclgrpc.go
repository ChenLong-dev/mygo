/**
* @Author: cl
* @Date: 2021/1/16 11:06
 */
package main

import (
	"context"
	"github.com/ChenLong-dev/gobase/mlog"
	"google.golang.org/grpc"
	"math"
	acl "myyd/src/availserverd/saasyd_acl_query_grpc"
	"net"
)

type server struct{}

func (s *server) GetAclResults(ctx context.Context, in *acl.AclResultReq) (*acl.AclResultRsp, error) {
	mlog.Infof("[offset:%d] [limt:%d] [type:%v]\n", in.GetOffset(), in.GetLimit(), in.GetType())

	aclResList, err := aclResultsRecords.GetAclResult()
	if err != nil {
		return &acl.AclResultRsp{
			Code: 444,
			Msg:  "failed",
			//Data: nil,
		}, err
	}

	var list []*acl.AclResultRsp_Data
	for _, v := range aclResList {
		if in.GetType() == "skip_unknown" && v.AclResult == "unknown" {
			continue
		}
		list = append(list, &acl.AclResultRsp_Data{BusinessUrl: v.Url, Result: v.AclResult})
	}
	mlog.Infof("get acl results [len:%d]\n", len(list))
	if len(list) == 0 {
		return &acl.AclResultRsp{
			Code: 555,
			Msg:  "no data",
			//Data: nil,
		}, err
	}

	return &acl.AclResultRsp{
		Code: 0,
		Msg:  "success",
		Data: list,
	}, nil
}

func (s *server) GetAclAllResults(ctx context.Context, in *acl.AclResultReq) (*acl.AclResultRsp, error) {
	mlog.Infof("[offset:%d] [limt:%d]\n", in.GetOffset(), in.GetLimit())

	aclResList, err := aclResultsRecords.GetAclResult()
	if err != nil {
		return &acl.AclResultRsp{
			Code: 444,
			Msg:  "failed",
			//Data: nil,
		}, err
	}

	var list []*acl.AclResultRsp_Data
	for _, v := range aclResList {
		list = append(list, &acl.AclResultRsp_Data{BusinessUrl: v.Url, Result: v.AclResult})
	}
	mlog.Infof("get acl results [len:%d]\n", len(list))
	if len(list) == 0 {
		return &acl.AclResultRsp{
			Code: 555,
			Msg:  "no data",
			//Data: nil,
		}, err
	}

	return &acl.AclResultRsp{
		Code: 0,
		Msg:  "success",
		Data: list,
	}, nil
}

func InitAclGrpcServer() error {
	mlog.Infof("param: InitAclGrpcServer listenAddr:%s", "0.0.0.0:36166")
	listen, err := net.Listen("tcp", "0.0.0.0:36166")
	if err != nil {
		mlog.Infof("listen is failed: %v\n", err)
		return err
	}
	var options = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(1073741824),
	}
	s := grpc.NewServer(options...)

	acl.RegisterAclGrpcServerServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		mlog.Fatalf("failed to serve: %v", err)
	}

	return nil

}