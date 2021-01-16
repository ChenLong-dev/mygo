/**
* @Author: cl
* @Date: 2021/1/16 15:35
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
	"time"
)

func (s *server) GetBusinessInfos(ctx context.Context, in *ws.GetBusinessInfosReq) (*ws.GetBusinessInfosRsp, error) {
	mlog.Infof("get a grpc [GetBusinessInfos] request: %+v", in)

	mlog.Infof("[refCluterNode:%s] [refClouduser:%s] [offset:%d] [limt:%d] [type:%s]\n",
		in.GetRefClusterNode(), in.GetRefClouduser(), in.GetOffset(), in.GetLimit(), in.GetMsgType())

	dataList, err := getBusinessInfo(in.GetRefClusterNode(), in.GetRefClouduser(), int64(in.GetOffset()),
		int64(in.GetLimit()))
	if err != nil {
		return &ws.GetBusinessInfosRsp{
			Code: 444,
			Msg:  "failed",
			//Data: "",
		}, err
	}

	vt1 := time.Now()
	defer func() {
		vt2 := time.Now()
		mlog.Infof("3 ==xxx== [vt1:%v] [vt2:%v] [expend:%v]", vt1, vt2, vt2.Sub(vt1))
	}()

	if len(dataList) == 0 {
		mlog.Error("data is nil")
		return &ws.GetBusinessInfosRsp{
			Code: 555,
			Msg:  "no data",
			Data: &ws.GetBusinessInfosRsp_Data{
				List: nil,
			},
		}, nil
	}

	var list []*ws.GetBusinessInfosRsp_BusinessInfo
	for _, v := range dataList {
		var ipArray []string
		if val, ok := v.Ip.(string); ok {
			ipArray = append(ipArray, val)
		} else if val, ok := v.Ip.(primitive.A); ok {
			for _, vIp := range val {
				ipArray = append(ipArray, vIp.(string))
			}
		} else {
			ipArray = append(ipArray, "")
			mlog.Warnf("xxxxx [url:%s] [ip:%+v]\n", v.Url, v.Ip)
		}
		var insiteIpArray []*ws.GetBusinessInfosRsp_BusinessInfo_InsiteIp
		for _, vInsite := range v.InsiteIp {
			insiteIpArray = append(insiteIpArray, &ws.GetBusinessInfosRsp_BusinessInfo_InsiteIp{
				ExtranetIp: vInsite.ExtranetIp,
				Type:       vInsite.Type,
			})
		}
		reserveInsiteIp := v.ReserveInsiteIp
		pReserveInsiteIp := &ws.GetBusinessInfosRsp_BusinessInfo_InsiteIp{
			ExtranetIp: reserveInsiteIp.ExtranetIp,
			Type:       reserveInsiteIp.Type,
		}

		businessInfo := &ws.GetBusinessInfosRsp_BusinessInfo{
			RefClouduser:    v.RefClouduser,
			Domain:          v.BusinessName,
			Url:             v.Url,
			BusinessName:    v.BusinessName,
			Port:            uint32(v.Port),
			Protocol:        v.Protocol,
			Ip:              ipArray,
			RefClusterNode:  v.RefClusterNode,
			InsiteIp:        insiteIpArray,
			RefReserveNode:  v.RefReserveNode,
			ReserveInsiteIp: pReserveInsiteIp,
			DnsType:         uint32(v.DnsType),
			CheckIntvl:      uint32(v.CheckIntvl),
			CheckLevel:      uint32(v.CheckLevel),
		}
		list = append(list, businessInfo)

	}
	data := &ws.GetBusinessInfosRsp_Data{
		List: list,
	}

	return &ws.GetBusinessInfosRsp{
		Code: 0,
		Msg:  "success",
		Data: data,
	}, nil
}

func getDomainInfos(refClusterNode, refClouduser string) (map[string]*scom.DomainInfos, error) {
	vt1 := time.Now()
	defer func() {
		vt2 := time.Now()
		mlog.Infof("1 ==xxx== [vt1:%v] [vt2:%v] [expend:%v]", vt1, vt2, vt2.Sub(vt1))
	}()
	collection := mg.CreateCollection("FCDomain")
	dFilter := bson.D{{"end_time", bson.M{"$gte": time.Now().Unix()}}, {"is_deleted", 0}, {"status", 3000}, {"is_offline", bson.M{"$ne": 1}},
		{"ref_cluster_node", bson.M{"$nin": []string{"cluster_node_sh", "cluster_node_xj_telecom", "cluster_node_fjqz_ChinaMobile"}}},
		{"$or", []interface{}{bson.D{{"is_in", 1}}, bson.D{{"is_in", 0}, {"dns_bypass_status", 1}, {"dns_type", 1}}}},
	}
	//dFilter := bson.D{{"is_deleted", 0}, {"status", 3000}, {"is_offline", bson.M{"$ne": 1}}}
	if refClusterNode != "" {
		dFilter = append(dFilter, bson.E{"ref_cluster_node", refClusterNode})
	}

	if refClusterNode != "" {
		dFilter = append(dFilter, bson.E{"ref_clouduser", refClouduser})
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
	domainInfosMap := make(map[string]*scom.DomainInfos)
	count := 0
	for cursor.Next(context.Background()) {
		var domainInfos *scom.DomainInfos
		if err = cursor.Decode(&domainInfos); err != nil {
			mlog.Warnf("[err:%s] [data:%+v]\n", err, domainInfos)
		}
		count++
		domainInfosMap[domainInfos.Domain] = domainInfos
		mlog.Debugf("get domain infos : [data: %v]", domainInfos)
	}
	if len(domainInfosMap) == 0 {
		return nil, fmt.Errorf("get domain data is nill! [%+v]\n", dFilter)
	}
	return domainInfosMap, nil
}

func getFCWhiteList() (map[string]string, error) {
	whiteMap := make(map[string]string)
	collection := mg.CreateCollection("FCWhite")
	cursor, err := mg.Find(collection, bson.D{{}}, 0, 0)
	if err != nil {
		mlog.Warn(err)
		return whiteMap, err
	}
	defer cursor.Close(context.Background())
	if err = cursor.Err(); err != nil {
		mlog.Warn(err)
		return whiteMap, err
	}
	type white struct {
		Url string `bson:"url"`
	}

	count := 0
	for cursor.Next(context.Background()) {
		var whiteUrl white
		if err = cursor.Decode(&whiteUrl); err != nil {
			mlog.Warnf("[err:%s] [data:%+v]\n", err, whiteUrl)
		}
		count++
		whiteMap[whiteUrl.Url] = "1"
	}
	mlog.Debugf("get white url infos : [data len: %d], [count: %d]", len(whiteMap), count)

	return whiteMap, nil
}

func getBusinessInfo(refClusterNode, refClouduser string, offset, limit int64) ([]*scom.BusinessInfos, error) {
	domainInfosMap, err := getDomainInfos(refClusterNode, refClouduser)
	if err != nil {
		return nil, err
	}
	mlog.Infof("get domain infos: [len:%d]\n", len(domainInfosMap))
	vt1 := time.Now()
	defer func() {
		vt2 := time.Now()
		mlog.Infof("2 ==xxx== [vt1:%v] [vt2:%v] [expend:%v]", vt1, vt2, vt2.Sub(vt1))
	}()

	var domainIDList []primitive.ObjectID
	fcWhiteList, _ := getFCWhiteList()

	for _, val := range domainInfosMap {
		if _, ok := fcWhiteList[val.Domain]; ok {
			continue
		}
		domainIDList = append(domainIDList, val.Id)
	}
	collection := mg.CreateCollection("FCWebsite")
	dFilter := bson.M{"domain_id": bson.M{"$in": domainIDList}}
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
	var businessInfosList []*scom.BusinessInfos
	count := 0
	for cursor.Next(context.Background()) {
		//vt := time.Now()
		var businessInfos *scom.BusinessInfos
		if err = cursor.Decode(&businessInfos); err != nil {
			mlog.Warnf("[err:%s] [data:%+v]\n", err, *businessInfos)
		}
		count++
		if businessInfos.Port == 0 {
			continue
		}
		for _, val := range businessInfosList {
			if val.Url == businessInfos.Url {
				mlog.Warnf("xxxxx [url:%+v] [%+v]\n", val.Url, *businessInfos)
			}
		}
		//mlog.Debugf("xxxx [%+v] [%v]\n", businessInfos, time.Now().Sub(vt))
		businessInfosList = append(businessInfosList, businessInfos)
	}

	for _, val := range businessInfosList {
		if value, ok := domainInfosMap[val.Domain]; ok {
			val.BusinessName = value.BusinessName
			val.DnsType = value.DnsType
			val.Ip = value.Ip
			val.InsiteIp = value.InsiteIp
			val.RefReserveNode = value.RefReserveNode
			if _, ok := value.ReserveInsiteIp.(string); ok {
				val.ReserveInsiteIp = scom.SiteIP{}
				mlog.Debugf("reserve insite ip is string [%v]\n", value.ReserveInsiteIp)
			} else if siteIp, ok := value.ReserveInsiteIp.(primitive.D); ok {
				for _, v := range siteIp {
					if v.Key == "extranet_ip" {
						val.ReserveInsiteIp.ExtranetIp = v.Value.(string)
					} else if v.Key == "type" {
						val.ReserveInsiteIp.Type = v.Value.(string)
					}
				}
				//mlog.Debugf("xxxxxxx  [primitive.D] [%v]\n", value.ReserveInsiteIp)
			} else {
				val.ReserveInsiteIp = scom.SiteIP{}
				mlog.Debugf("reserve insite ip is other [%v]\n", value.ReserveInsiteIp)
			}
			//mlog.Debugf("xxx [%+v] [%T]\n", val, value.ReserveInsiteIp)
		}
	}

	mlog.Infof("get data count [domain_map:%d] [count:%d] [business_list:%d]\n",
		len(domainInfosMap), count, len(businessInfosList))
	return businessInfosList, nil
}
