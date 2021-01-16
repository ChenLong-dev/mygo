/**
* @Author: cl
* @Date: 2021/1/16 11:10
 */
package main

import (
	"context"
	"fmt"
	"github.com/ChenLong-dev/gobase/mg"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyd/src/availserverd/alarm"
	"myyd/src/scom"
	"time"
)

func (rsm *ResultsMan) AddLinkPoint(pr *PointResult) {
	//mlog.Debugf("==== begin AddPoint")
	mlog.Tracef("pr.Time=%v", pr.CheckTime)
	rsm.Lock()

	rsm.lastPoints.PushBack(pr)
	if rsm.lastPoints.Len() > rsm.keepPointsNum {
		rsm.lastPoints.Remove(rsm.lastPoints.Front())
	}

	rsm.Unlock()

	//	todo 存储到数据库

	results := rsm.GetResults(PointBound, RegionBound)
	if len(results) <= 0 {
		mlog.Warnf("get result is nil [pr.len:%d]\n", len(pr.Results))
		return
	}
	mlog.Tracef("GetResults len:%d, val:%+v", len(results), results)
	vt1 := time.Now()
	_ = resultRecords.CheckLinkResults(pr, results)
	vt2 := time.Now()
	mlog.Infof("1 x=x=x=x [results.len:%d] [exp:%v]\n", len(results), vt2.Sub(vt1))
	if err := rsm.HandleLinkDetect(results); err != nil {
		mlog.Error(err)
		return
	}
	vt3 := time.Now()
	mlog.Infof("2 x=x=x=x [pr.len:%d] [exp:%v]\n", len(pr.Results), vt3.Sub(vt2))
}

func (rsm *ResultsMan) HandleLinkDetect(results map[string]*CheckResult) error {
	vt1 := time.Now()
	var normal, abnormal []string
	alarmInfo := make(map[string]string)
	for url, val := range results {
		if val.PrimaryResult == ResultValue_OK && val.SecondResult == ResultValue_OK {
			//发送MQ消息
			normal = append(normal, url)
		} else if val.PrimaryResult == ResultValue_OK && val.SecondResult == ResultValue_FAILED {
			//发送MQ消息
			normal = append(normal, url)
			if aclRes, ok := scom.AclUrlRes[url]; ok {
				if aclRes == "acl" {
					//告警通知
					content := fmt.Sprintf("[%s], url:[%s], reason:[%s], time:[%+v]",
						DetectType, url, "备节点故障（源站已配置ACL，主节点访问正常，备节点访问失败）", time.Now().Format("2006-01-02 15:04:05"))
					alarmInfo[url] = content
				} else if aclRes == "noacl" {
					//告警通知
					content := fmt.Sprintf("[%s], url:[%s], reason:[%s], time:[%+v]",
						DetectType, url, "源站故障（源站未配置ACL，主节点访问正常，源站不可访问）", time.Now().Format("2006-01-02 15:04:05"))
					alarmInfo[url] = content
				} else {
					mlog.Debugf("scom.AclUrlRes aclRes is nil [url:%s] [aclRes:%s]", url, aclRes)
				}
			} else {
				mlog.Debugf("scom.AclUrlRes is nil [url:%s]", url)
			}
		} else if val.PrimaryResult == ResultValue_FAILED && val.SecondResult == ResultValue_OK {
			//告警通知
			content := fmt.Sprintf("[%s], url:[%s], reason:[%s], time:[%+v]",
				DetectType, url, "链路检测过云盾主节点访问源站失败", time.Now().Format("2006-01-02 15:04:05"))
			alarmInfo[url] = content
			//发送MQ消息
			abnormal = append(abnormal, url)
		} else if val.PrimaryResult == ResultValue_FAILED && val.SecondResult == ResultValue_FAILED {
			if aclRes, ok := scom.AclUrlRes[url]; ok {
				if aclRes == "acl" {
					//告警通知
					content := fmt.Sprintf("[%s], url:[%s], reason:[%s], time:[%+v]",
						DetectType, url, "主节点故障（源站已配置ACL，主节点不可访问，备节点不可访问）", time.Now().Format("2006-01-02 15:04:05"))
					alarmInfo[url] = content
				} else if aclRes == "noacl" {
					//告警通知
					content := fmt.Sprintf("[%s], url:[%s], reason:[%s], time:[%+v]",
						DetectType, url, "源站故障（源站未配置ACL，主节点不可访问，源站不可访问）", time.Now().Format("2006-01-02 15:04:05"))
					alarmInfo[url] = content
				} else {
					mlog.Debugf("scom.AclUrlRes aclRes is nil")
				}
			} else {
				mlog.Debugf("scom.AclUrlRes is nil")
			}
		} else {
			mlog.Debugf("results key:%s, val:%v", url, val)
			continue
		}
	}
	vt2 := time.Now()
	mlog.Infof("1 v=v=v=v [exp:%v]\n", vt2.Sub(vt1))

	//发送MQ消息
	mlog.Debugf("HandleLinkDetect len(normal):%v, len(abnormal):%v", len(normal), len(abnormal))
	if len(normal)+len(abnormal) > 0 {
		go ProducerMq(LINK_DETECT, normal, abnormal)
	}
	//发送告警通知
	if len(alarmInfo) > 0 {
		go alarm.SendAlarmMsgs(alarmInfo)
	}

	vt3 := time.Now()
	mlog.Infof("2 v=v=v=v [exp:%v]\n", vt3.Sub(vt2))

	return nil
}

func checkLineResult(lrs []scom.LineResult, failedOne bool) ResultValue {
	failed := 0
	ok := 0
	for _, lr := range lrs {
		/*if lr.Result.Code < 0 {
			return ResultValue_FAILED	//	只要有一条线路不可用就认为是故障
		}*/
		if lr.Result.Code < 0 {
			failed++
		} else {
			ok++
		}
	}
	//return ResultValue_OK
	mlog.Debugf("mmmmm 1 [lrs:%v] [ok:%d] [failed:%d]", lrs, ok, failed)
	if failedOne {
		if failed > 0 {
			return ResultValue_FAILED
		} else {
			return ResultValue_OK
		}
	} else {
		if ok > 0 {
			return ResultValue_OK
		} else {
			return ResultValue_FAILED
		}
	}
}

func checkUrlResult(m map[string]*scom.RegionResult, regionBound int, failedOne bool, url string) ResultValue {
	ok := 0
	failed := 0
	for r, rr := range m {
		switch rv := checkLineResult(rr.LineResults, failedOne); rv {
		case ResultValue_OK:
			ok++
		case ResultValue_FAILED:
			failed++
		}
		mlog.Debugf("mmmmm 2 [url:%s] [r:%s] [rr:%v] [ok:%d] [failed:%d]", url, r, rr, ok, failed)
	}
	mlog.Debugf("mmmmm 3 [url:%s] [ok:%d] [failed:%d] [regionBound:%d]", url, ok, failed, regionBound)
	if ok >= regionBound {
		return ResultValue_OK
	} else if failed >= regionBound {
		return ResultValue_FAILED
	} else {
		return ResultValue_UNKNOWN
	}
}

func (rsm *ResultsMan) GetResults(pointBound int, regionBound int) (mcr map[string]*CheckResult) {
	mlog.Infof("keepPointsNum=%d,pointBound=%d,regionBound=%d", rsm.keepPointsNum, pointBound, regionBound)

	if pointBound > rsm.keepPointsNum {
		pointBound = rsm.keepPointsNum
	} else if pointBound <= rsm.keepPointsNum/2 {
		pointBound = rsm.keepPointsNum/2 + 1
	}

	rsm.RLock()
	defer rsm.RUnlock()

	type resultStat struct {
		ok      int
		failed  int
		unknown int
	}
	addRS := func(rs *resultStat, rv ResultValue) {
		switch rv {
		case ResultValue_OK:
			rs.ok++
		case ResultValue_FAILED:
			rs.failed++
		case ResultValue_UNKNOWN:
			rs.unknown++
		}
	}
	type Stat struct {
		primary resultStat
		second  resultStat
	}
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

			rv := checkUrlResult(urlResult.Primarys, regionBound, true, url)
			addRS(&s.primary, rv)
			rv = checkUrlResult(urlResult.Seconds, regionBound, false, url)
			addRS(&s.second, rv)
		}
	}

	mcr = make(map[string]*CheckResult)
	check := func(rs *resultStat, url string) ResultValue {
		mlog.Debugf("mmmmm 4 [url:%s] [rs:%v] [pointBound:%d]\n", url, rs, pointBound)
		if rs.ok >= pointBound { // 成功：需要过检测次数的一半
			return ResultValue_OK
		} else if rs.failed >= pointBound { // 故障：需要过检测次数的一半
			return ResultValue_FAILED
		} else {
			return ResultValue_UNKNOWN
		}
	}
	for url, s := range ss {
		cr := &CheckResult{}
		cr.Url = url
		cr.PrimaryResult = check(&s.primary, url)
		cr.SecondResult = check(&s.second, url)

		mcr[url] = cr
	}

	return mcr
}

func (rr *ResultsRecords) CheckLinkResults(pr *PointResult, mcr map[string]*CheckResult) error {
	checkTime := pr.CheckTime
	createTime := time.Now()
	var infoRecords []interface{}
	resultsMap := pr.Results
	vt1 := time.Now()
	mlog.Infof("1 ====== pr.map.len:%d, mcr.len:%d, rr.map.len:%d\n", len(resultsMap), len(mcr), len(rr.Results))

	updateRes := func(url string, mcrRes *CheckResult, rrRes CheckResult, res *UrlResult, isFirst bool) {
		// 添加到对比数据缓存中
		rr.Results[url] = *mcrRes
		rr.Results2[url] = res
		// 添加到入库数据列表中
		if isFirst {
			infoRecords = append(infoRecords, &ResultRecords{Url: url, CheckTime: checkTime, CreateTime: createTime,
				PrePrimaryResult: 99, PreSecondResult: 99, NowPrimaryResult: mcrRes.PrimaryResult,
				NowSecondResult: mcrRes.SecondResult, Results: res})
		} else {
			infoRecords = append(infoRecords, &ResultRecords{Url: url, CheckTime: checkTime, CreateTime: createTime,
				PrePrimaryResult: rrRes.PrimaryResult, PreSecondResult: rrRes.SecondResult,
				NowPrimaryResult: mcrRes.PrimaryResult, NowSecondResult: mcrRes.SecondResult, Results: res})
		}
	}

	rr.RLock()
	for url, res1 := range resultsMap {
		var mcrRes *CheckResult
		if rrRes, ok := rr.Results[url]; !ok {
			if mcrRes, ok = mcr[url]; ok {
				updateRes(url, mcrRes, rrRes, res1, true)
			} else {
				mlog.Warnf("1 [url:%s] not result in mcr!", url)
			}
		} else {
			if mcrRes, ok = mcr[url]; ok { // 找到，检查总结果
				if rrRes.PrimaryResult != mcrRes.PrimaryResult || rrRes.SecondResult != mcrRes.SecondResult {
					updateRes(url, mcrRes, rrRes, res1, false)
					continue
				} else {
					if res2, ok := rr.Results2[url]; ok {
						if compareLinkRes(res1, res2) {
							continue
						} else {
							updateRes(url, mcrRes, rrRes, res1, false)
						}
					} else {
						mlog.Warnf("2 [url:%s] not result in rr.Results2!", url)
					}
				}
			}
		}
	}
	mlog.Infof("2 ====== add link records :%d", len(infoRecords))
	rr.RUnlock()
	if len(infoRecords) == 0 {
		return fmt.Errorf("data is nil")
	}
	vt2 := time.Now()
	mlog.Infof("1 cl=cl=cl=cl [exp:%v]\n", vt2.Sub(vt1))

	if _, err := mg.MongoDB.Collection("LinkRecords").InsertMany(context.TODO(), infoRecords); err != nil {
		mlog.Errorf("add link results error: %v\n", err)
		return err
	}

	vt3 := time.Now()
	mlog.Infof("2 cl=cl=cl=cl [exp:%v]\n", vt3.Sub(vt2))

	return nil
}

func (rr *ResultsRecords) GetResults() (map[string]CheckResult, error) {
	rr.RLock()
	defer rr.RUnlock()
	results := rr.Results
	if len(results) == 0 {
		return nil, fmt.Errorf("results is nil")
	}
	return results, nil
}
