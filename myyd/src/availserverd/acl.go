/**
* @Author: cl
* @Date: 2021/1/16 11:05
 */
package main

import (
	"context"
	"fmt"
	"github.com/ChenLong-dev/gobase/mg"
	"github.com/ChenLong-dev/gobase/mlog"
	"go.mongodb.org/mongo-driver/bson"
	"myyd/src/scom"
	"strings"
	"sync"
	"time"
)

type AclResultsRecords struct {
	sync.RWMutex
	Results map[string]*scom.AclResults
}

func NewAclResultsMan() *AclResultsRecords {
	return &AclResultsRecords{Results: make(map[string]*scom.AclResults)}
}

func (rsm *ResultsMan) AddAclPoint(pr *PointResult) {
	mlog.Tracef("pr.Time=%v", pr.CheckTime)
	rsm.Lock()
	rsm.lastPoints.PushBack(pr)
	if rsm.lastPoints.Len() > rsm.keepPointsNum {
		rsm.lastPoints.Remove(rsm.lastPoints.Front())
	}
	rsm.Unlock()

	results := rsm.GetAclResult(pr)
	if len(results) <= 0 {
		mlog.Warnf("get result is nil [pr.len:%d]\n", len(pr.Results))
		return
	}
	if err := aclResultsRecords.CheckAclResults(results); err != nil {
		mlog.Error(err)
	}
}

func (rsm *ResultsMan) GetAclResult(pr *PointResult) map[string]*scom.AclResults {
	checkAclResult := func(m map[string]*scom.RegionResult) (cloud int, cluster int) {
		for region, res := range m {
			if strings.Contains(region, "cloud") {
				for _, lr := range res.LineResults {
					if lr.Result.Code > 0 { // 只要公有云节点有一个通，就是noacl
						cloud++
					}
				}
			}
			if strings.Contains(region, "cluster") {
				for _, lr := range res.LineResults {
					if lr.Result.Code > 0 {
						cluster++
					}
				}
			}
		}
		return cloud, cluster
	}
	type Stat struct {
		clo int
		clu int
		res map[string]*scom.RegionResult
	}

	rsm.RLock()
	ss := make(map[string]*Stat)
	for e := rsm.lastPoints.Front(); e != nil; e = e.Next() {
		pr := e.Value.(*PointResult)
		for url, urlResult := range pr.Results {
			var s *Stat
			if r, ok := ss[url]; ok {
				s = r
			} else {
				s = &Stat{}
				ss[url] = s
			}
			clo, clu := checkAclResult(urlResult.Primarys)
			s.clo += clo
			s.clu += clu
		}
	}
	rsm.RUnlock()

	ar := make(map[string]*scom.AclResults)
	for url, st := range ss {
		var res string
		if st.clo > 0 && st.clo <= 5 { // 有站点防护的
			res = "xxx"
		} else if st.clo > 5 { // 只要有个公有云通，就是noacl
			res = "noacl"
		} else {
			if st.clu > 0 { // 云节点不通， 但机房节点通
				res = "acl"
			} else { // 云节点和机房节点都不通
				res = "unknown"
			}
		}
		if ures, ok := pr.Results[url]; ok {
			ar[url] = &scom.AclResults{Url: url, AclResult: res, Results: ures.Primarys,
				CheckerTime: pr.CheckTime, CreateTime: time.Now()}
		}
	}
	return ar
}

func FindAclResult() (map[string]*scom.AclResults, error) {
	collection := mg.CreateCollection("ACLResults")
	cursor, err := mg.Find(collection, bson.D{{}}, 0, 0)
	if err != nil {
		mlog.Warn(err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	if err = cursor.Err(); err != nil {
		mlog.Warn(err)
		return nil, err
	}
	aclResultMap := make(map[string]*scom.AclResults)
	for cursor.Next(context.Background()) {
		var aclResult scom.AclResults
		if err = cursor.Decode(&aclResult); err != nil {
			mlog.Warn(err)
			return nil, err
		}
		aclResultMap[aclResult.Url] = &aclResult
		//mlog.Debugf("acl:%v\n", aclResult)
	}
	return aclResultMap, nil
}

func compareAclRes(p1, p2 map[string]*scom.RegionResult) bool {
	for region, rr1 := range p1 {
		if rr2, ok := p2[region]; ok {
			return compareLines(rr1.LineResults, rr2.LineResults)
		} else {
			return false
		}
	}
	return true
}

func (arr *AclResultsRecords) CheckAclResults(aclRes map[string]*scom.AclResults) error {
	mlog.Infof("1 ====== arr.map.len:%d, aclRes.len:%d\n", len(arr.Results), len(aclRes))
	vt1 := time.Now()
	defer func() {
		vt2 := time.Now()
		mlog.Infof("4 ====== [exp:%v]\n", vt2.Sub(vt1))
	}()
	if len(arr.Results) == 0 { //第一次加载
		if aclMap, err := FindAclResult(); err == nil {
			mlog.Infof("2 === load acl cache results [len:%d]\n", len(aclMap))
			for url, res := range aclMap {
				mlog.Debugf("[%s] [res:%T] [res:%v]\n", url, res, res)
				arr.Results[url] = res
			}
		}
	}

	var AddInfoResults []interface{}
	var UpateInfoResults []*scom.AclResults
	var AddInfoRecords []interface{}
	for url, res1 := range aclRes {
		if res2, ok := arr.Results[url]; !ok {
			arr.Results[url] = res1
			AddInfoResults = append(AddInfoResults, res1)
			AddInfoRecords = append(AddInfoRecords, res1)
		} else {
			if res1.AclResult == res2.AclResult {
				// 比较细节
				if compareAclRes(res1.Results, res2.Results) {
					continue
				} else {
					arr.Results[url] = res1
					AddInfoRecords = append(AddInfoRecords, res1)
				}
			} else {
				arr.Results[url] = res1
				UpateInfoResults = append(UpateInfoResults, res1)
				AddInfoRecords = append(AddInfoRecords, res1)
			}
		}
	}
	mlog.Infof("3 ====== update or add acl [add res:%d] [update res:%d] [add rec:%d]", len(AddInfoResults),
		len(UpateInfoResults), len(AddInfoRecords))

	if len(AddInfoResults) != 0 { // 增加ACLResults
		if _, err := mg.MongoDB.Collection("ACLResults").InsertMany(context.TODO(), AddInfoResults); err != nil {
			mlog.Errorf("add acl results error: %v\n", err)
			return err
		}
	}

	if len(UpateInfoResults) != 0 { // 更新ACLResults
		for _, value := range UpateInfoResults {
			filter := bson.D{{"url", value.Url}}
			update := bson.D{{"$set", bson.D{{"acl_result", value.AclResult}, {"results", value.Results},
				{"check_time", value.CheckerTime}, {"create_time", value.CreateTime}}}}
			if _, err := mg.MongoDB.Collection("ACLResults").UpdateOne(context.TODO(), filter, update); err != nil {
				mlog.Errorf("update acl results error: %v\n", err)
				return err
			}
		}
	}

	if len(AddInfoRecords) != 0 { // 增加ACLRecords
		if _, err := mg.MongoDB.Collection("ACLRecords").InsertMany(context.TODO(), AddInfoRecords); err != nil {
			mlog.Errorf("add acl records error: %v\n", err)
			return err
		}
	}

	return nil
}

func (arr *AclResultsRecords) GetAclResult() ([]*scom.AclRes, error) {
	arr.RLock()
	defer arr.RUnlock()
	results := arr.Results
	if len(results) == 0 {
		return nil, fmt.Errorf("results is nil")
	}
	var arList []*scom.AclRes
	for url, res := range results {
		arList = append(arList, &scom.AclRes{Url: url, AclResult: res.AclResult})
	}
	return arList, nil
}
