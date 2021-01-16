/**
* @Author: cl
* @Date: 2021/1/16 11:11
 */
package main

import (
	"context"
	"github.com/ChenLong-dev/gobase/mlog"
	"google.golang.org/grpc"
	"math"
	"myyd/src/scom"
	"net"
	link "myyd/src/availserverd/saasyd_link_query_grpc"
)

type linkServer struct{}

func copyMap(regionResult map[string]*scom.RegionResult) map[string]*link.RegionResult {
	resRegionMap := make(map[string]*link.RegionResult)

	for kR, vR := range regionResult {
		var resRegion link.RegionResult
		resRegion.Region = vR.Region

		for _, vL := range vR.LineResults {
			resRegion.LineResults = append(resRegion.LineResults, &link.RegionResult_LineResults{
				Line: &link.RegionResult_LineResults_Line{
					Addr: vL.Line.Addr,
					Isp:  int32(vL.Line.Isp),
				},
				Result: &link.RegionResult_LineResults_Result{
					Code:   int32(vL.Result.Code),
					Status: vL.Result.Status,
					Delay:  vL.Result.Delay.Microseconds(),
				},
			})
		}

		resRegionMap[kR] = &resRegion
	}

	return resRegionMap
}

func (s *linkServer) GetLastPointResult(ctx context.Context, in *link.LinkResultReq) (*link.LinkResultRsp, error) {
	mlog.Infof("[offset:%d] [limt:%d]\n", in.GetOffset(), in.GetLimit())

	var resultList = &PointResult{}
	if in.GetType() == link.Type_LinkResult {
		resultList = defaultResultsMan.GetLastPointResult()
	} else if in.GetType() == link.Type_TargetResult {
		resultList = targetResultsMan.GetLastPointResult()
	}
	if resultList == nil {
		return &link.LinkResultRsp{
			Code: 555,
			Msg:  "no data",
		}, nil
	}

	if len(resultList.Results) == 0 {
		return &link.LinkResultRsp{
			Code: 555,
			Msg:  "no data",
		}, nil
	}

	var resultRsp link.LinkResultRsp_Data
	resultRsp.Checktime = resultList.CheckTime.Unix()

	resMap := make(map[string]*link.UrlResult)
	for kR, vR := range resultList.Results {
		var urlRes link.UrlResult
		urlRes.Url = vR.Url
		urlRes.Primarys = copyMap(vR.Primarys)
		urlRes.Seconds = copyMap(vR.Seconds)
		resMap[kR] = &urlRes
	}
	resultRsp.Results = resMap

	return &link.LinkResultRsp{
		Code: 0,
		Msg:  "success",
		Data: &resultRsp,
	}, nil
}

func (s *linkServer) GetResults(ctx context.Context, in *link.ResultsReq) (*link.ResultsRsp, error) {
	mlog.Infof("[PointBound:%d] [RegionBound:%d]\n", in.GetPointBound(), in.GetRegionBound())
	var results = map[string]CheckResult{}
	var err error
	if in.GetType() == link.Type_LinkResult {
		results, err = resultRecords.GetResults()
	} else if in.GetType() == link.Type_TargetResult {
		results, err = targetResultRecords.GetTargetResults()
	}

	if err != nil {
		return &link.ResultsRsp{
			Code: 444,
			Msg:  "failed",
		}, nil
	}

	resMap := make(map[string]*link.CheckResult)

	for kR, vR := range results {
		resMap[kR] = &link.CheckResult{
			Url:           vR.Url,
			PrimaryResult: int32(vR.PrimaryResult),
			SecondResult:  int32(vR.SecondResult),
		}
	}

	if len(resMap) == 0 {
		return &link.ResultsRsp{
			Code: 555,
			Msg:  "no data",
		}, nil
	}

	return &link.ResultsRsp{
		Code: 0,
		Msg:  "success",
		Data: resMap,
	}, nil
}

func (s *linkServer) GetAddrResults(ctx context.Context, in *link.ResultsAddrReq) (*link.ResultsAddrRsp, error) {
	mlog.Infof("get a GetAddrResults req [%v]\n", in)
	//var results = map[string]AddrTargetResult{}
	results, err := targetResultRecords.GetAddrResults()
	if err != nil {
		return &link.ResultsAddrRsp{
			Code: 444,
			Msg:  "failed",
		}, nil
	}

	resMap := make(map[string]*link.CheckAddrResult)

	for kR, vR := range results {
		line := &link.Line{
			Addr: vR.Line.Addr,
			Isp:  int32(vR.Line.Isp),
		}
		resMap[kR] = &link.CheckAddrResult{
			NodeName: vR.NodeName,
			Result:   int32(vR.Result),
			Line:     line,
		}
	}

	if len(resMap) == 0 {
		return &link.ResultsAddrRsp{
			Code: 555,
			Msg:  "no data",
		}, nil
	}

	return &link.ResultsAddrRsp{
		Code: 0,
		Msg:  "success",
		Data: resMap,
	}, nil
}

func InitLinkGrpcServer(addr string) error {
	mlog.Infof("param: InitGrpcServer listenAddr:%s", addr)
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

	link.RegisterLinkGrpcServerServer(s, &linkServer{})

	if err := s.Serve(listen); err != nil {
		mlog.Fatalf("failed to serve: %v", err)
	}

	return nil

}
