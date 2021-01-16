/**
* @Author: cl
* @Date: 2021/1/16 15:38
 */
package main

import (
	"context"
	"fmt"
	"github.com/ChenLong-dev/gobase/mg"
	"github.com/ChenLong-dev/gobase/mlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	ws "myyd/src/saasydwebsite/website"
	"myyd/src/scom"
)

func (s *server) GetClusterIpInfos(ctx context.Context, in *ws.GetClusterIPInfosReq) (*ws.GetClusterIPInfosRsp, error) {
	mlog.Infof("get a grpc [GetClusterIpInfos] request: %+v", in)
	offset := in.GetOffset()
	limit := in.GetLimit()
	msgType := in.GetMsgType()
	mlog.Infof("[offset:%d] [limt:%d] [msg_type:%s]\n", offset, limit, msgType)

	dataList, err := getClusterIPInfo(int64(offset), int64(limit), msgType)
	if err != nil {
		return &ws.GetClusterIPInfosRsp{
			Code: 444,
			Msg:  "failed",
			//Data: "",
		}, err
	}

	if len(dataList) == 0 {
		mlog.Error("data is nil")
		return &ws.GetClusterIPInfosRsp{
			Code: 555,
			Msg:  "no data",
			Data: &ws.GetClusterIPInfosRsp_Data{
				List: nil,
			},
		}, nil
	}

	var list []*ws.GetClusterIPInfosRsp_ClusterIPInfo
	for _, v := range dataList {
		clusterIPInfo := &ws.GetClusterIPInfosRsp_ClusterIPInfo{
			NodeName:        v.NodeName,
			OwnerType:       v.OwnerType,
			RefServiceGroup: v.RefServiceGroup,
			Ip: &ws.GetClusterIPInfosRsp_ClusterIPInfo_Ip{
				ExtranetIp: v.Ip.ExtranetIp,
				Type:       v.Ip.Type,
			},
			Type: v.Type,
		}
		list = append(list, clusterIPInfo)
	}
	data := &ws.GetClusterIPInfosRsp_Data{
		List: list,
	}

	return &ws.GetClusterIPInfosRsp{
		Code: 0,
		Msg:  "success",
		Data: data,
	}, nil
}

func getClusterIPInfo(offset, limit int64, msgType string) ([]*scom.ClusterIPInfos, error) {
	collection := mg.CreateCollection("FCClusterIP")
	noInRegion := []string{"cluster_node_xj_telecom", "cluster_node_sh", "cluster_node_fjqz_ChinaMobile", "cluster_node_xa"}
	dFilter := bson.M{"type": bson.M{"$in": []int{3, 4}}, "node_name": bson.M{"$nin": noInRegion}}
	cursor, err := mg.Find(collection, dFilter, limit, offset)
	if err != nil {
		mlog.Warn(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	if err = cursor.Err(); err != nil {
		mlog.Warn(err)
		return nil, err
	}
	var clusterIPInfosList []*scom.ClusterIPInfos
	count := 0
	for cursor.Next(context.Background()) {
		var clusterIPInfos *scom.ClusterIPInfos
		if err = cursor.Decode(&clusterIPInfos); err != nil {
			mlog.Warnf("[err:%s] [data:%+v]\n", err, *clusterIPInfos)
		}
		count++
		clusterIPInfosList = append(clusterIPInfosList, clusterIPInfos)
	}

	if msgType == "line" || msgType == "target" {
		return getMsgTypeInfos(clusterIPInfosList, msgType)
	}

	mlog.Infof("get data count [msgtype:%s] [count:%d] [clusterip_list:%d]\n", msgType, count, len(clusterIPInfosList))
	return clusterIPInfosList, nil
}

type OwnerVsInsiteIpInfos struct {
	IpId  primitive.ObjectID `bson:"ip_id"`
	Owner string             `bson:"owner"`
	Type  int                `bson:"type"`
}

func getOwnerVsInsiteIp(msgType string) (map[primitive.ObjectID]*OwnerVsInsiteIpInfos, error) {
	collection := mg.CreateCollection("FCOwnerVsInsiteIp")
	var dFilter bson.D
	// type: 1:独享IP账户，2:独享IP URL，3:独享用户的备用节点IP
	if msgType == "line" {
		dFilter = bson.D{{"type", 2}, {"is_deleted", 0}}
	} else if msgType == "target" {
		dFilter = bson.D{{"type", 1}, {"is_deleted", 0}}
	} else {
		mlog.Warnf("not found type [%s]", msgType)
		return nil, fmt.Errorf("not found type [%s]", msgType)
	}

	cursor, err := mg.Find(collection, dFilter, 0, 0)
	if err != nil {
		mlog.Warn(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	if err = cursor.Err(); err != nil {
		mlog.Warn(err)
		return nil, err
	}
	OwnerVsInsiteIpInfosMap := make(map[primitive.ObjectID]*OwnerVsInsiteIpInfos)
	count := 0
	for cursor.Next(context.Background()) {
		var clusterIPInfos *OwnerVsInsiteIpInfos
		if err = cursor.Decode(&clusterIPInfos); err != nil {
			mlog.Warnf("[err:%s] [data:%+v]\n", err, *clusterIPInfos)
		}
		count++
		OwnerVsInsiteIpInfosMap[clusterIPInfos.IpId] = clusterIPInfos
	}
	return OwnerVsInsiteIpInfosMap, nil
}

func getMsgTypeInfos(ciis []*scom.ClusterIPInfos, msgType string) ([]*scom.ClusterIPInfos, error) {
	OwnerVsInsiteIpInfosMap, err := getOwnerVsInsiteIp(msgType)
	if err != nil {
		mlog.Warnf("find OwnerVsInsiteIp is failed [%v]", err)
		return nil, err
	}
	var resList []*scom.ClusterIPInfos
	getTarget := func(cii *scom.ClusterIPInfos) {
		if cii.Type == 3 {
			resList = append(resList, cii)
		} else { // type = 4
			if _, ok := OwnerVsInsiteIpInfosMap[cii.Id]; ok {
				resList = append(resList, cii)
			}
		}
	}

	getLine := func(cii *scom.ClusterIPInfos) {
		if cii.Type == 4 {
			if _, ok := OwnerVsInsiteIpInfosMap[cii.Id]; ok {
				resList = append(resList, cii)
			}
		}
	}

	for _, cii := range ciis {
		if msgType == "target" {
			getTarget(cii)
		} else { // line
			getLine(cii)
		}
	}

	mlog.Infof("get msgtype data [msgtype:%s] [ciis.len:%d] [owner.len:%d] [len:%d]", msgType, len(ciis), len(OwnerVsInsiteIpInfosMap), len(resList))
	return resList, nil
}