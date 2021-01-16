/**
* @Author: cl
* @Date: 2021/1/16 15:39
 */
package main
import (
	"context"
	"github.com/ChenLong-dev/gobase/mg"
	"github.com/ChenLong-dev/gobase/mlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	ws "myyd/src/saasydwebsite/website"
)

func (s *server) GetDomainIp(ctx context.Context, in *ws.GetDomainIPReq) (*ws.GetDomainIPRsp, error) {
	mlog.Infof("get a grpc [GetDomainIp] request: %+v", in)
	dataList, err := getDomainIP()
	if err != nil {
		return &ws.GetDomainIPRsp{
			Code: 444,
			Msg:  "failed",
			//Data: "",
		}, err
	}

	if len(dataList) == 0 {
		mlog.Error("data is nil")
		return &ws.GetDomainIPRsp{
			Code: 555,
			Msg:  "no data",
			Data: &ws.GetDomainIPRsp_Data{
				List: nil,
			},
		}, nil
	}

	var list []*ws.GetDomainIPRsp_Info
	for _, v := range dataList {
		var domainIP ws.GetDomainIPRsp_Info
		domainIP.Domain = v.Domain
		domainIP.IsIn = v.IsIn
		var ipArray []string
		if val, ok := v.Ip.(string); ok {
			ipArray = append(ipArray, val)
		} else if val, ok := v.Ip.(primitive.A); ok {
			for _, vIp := range val {
				ipArray = append(ipArray, vIp.(string))
			}
		} else {
			ipArray = append(ipArray, "")
			mlog.Warnf("xxxxx [domain:%s] [ip:%+v]\n", v.Domain, v.Ip)
		}
		domainIP.Ip = ipArray
		list = append(list, &domainIP)
	}

	data := &ws.GetDomainIPRsp_Data{
		List: list,
	}

	return &ws.GetDomainIPRsp{
		Code: 0,
		Msg:  "success",
		Data: data,
	}, nil
}

type DomainIP struct {
	Domain string      `bson:"domain"`
	Ip     interface{} `bson:"ip"`
	IsIn   uint32      `bson:"is_in"`
}

func getDomainIP() ([]*DomainIP, error) {
	collection := mg.CreateCollection("FCDomain")
	dFilter := bson.D{{"is_deleted", 0}, {"status", 3000}, {"is_in", 1}}
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
	var domainIPList []*DomainIP
	count := 0
	for cursor.Next(context.Background()) {
		var domainIP *DomainIP
		if err = cursor.Decode(&domainIP); err != nil {
			mlog.Warnf("[err:%s] [data:%+v]\n", err, *domainIP)
		}
		count++
		domainIPList = append(domainIPList, domainIP)
	}

	mlog.Infof("get data count [count:%d] [domainip_list:%d]\n", count, len(domainIPList))
	return domainIPList, nil
}

