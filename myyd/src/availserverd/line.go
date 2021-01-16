/**
* @Author: cl
* @Date: 2021/1/16 11:09
 */
package main

import (
	"context"
	"fmt"
	"github.com/ChenLong-dev/gobase/mg"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyd/src/availserverd/alarm"
	"myyd/src/scom"
	"reflect"
	"time"
)

func (rsm *ResultsMan) AddLinePoint(pr *PointResult, cksLen int) {
	//mlog.Debugf("==== begin AddPoint")
	mlog.Tracef("pr.Time=%v", pr.CheckTime)
	rsm.Lock()

	rsm.lastPoints.PushBack(pr)
	if rsm.lastPoints.Len() > rsm.keepPointsNum {
		rsm.lastPoints.Remove(rsm.lastPoints.Front())
	}

	rsm.Unlock()

	//	todo 存储到数据库

	results := rsm.GetLineResults(cksLen)
	if len(results) <= 0 {
		mlog.Warnf("get result is nil [pr.len:%d]\n", len(pr.Results))
		return
	}
	mlog.Debugf("GetResults len:%d, val:%+v", len(results), results)
	vt1 := time.Now()
	alarmMap, err := resultRecords.CheckLineResults(pr, results, cksLen)
	if err != nil {
		mlog.Warn(err)
	}
	vt2 := time.Now()
	mlog.Infof("1 x=x=x=x [alarmMap.len:%d] [exp:%v]\n", len(alarmMap), vt2.Sub(vt1))
	if err := rsm.HandleLineDetect(alarmMap); err != nil {
		mlog.Error(err)
		return
	}
	vt3 := time.Now()
	mlog.Infof("2 x=x=x=x [pr.len:%d] [exp:%v]\n", len(pr.Results), vt3.Sub(vt2))
}

func (rsm *ResultsMan) HandleLineDetect(results map[string]*CheckResult) error {
	var normal, abnormal []string
	alarmInfo := make(map[string]string)
	for key, val := range results {
		mlog.Debugf("HandleLineDetect results key:%s, val:%v", key, val)
		if val.PrimaryResult == ResultValue_OK {
			//发送MQ消息
			normal = append(normal, key)
		} else if val.PrimaryResult == ResultValue_FAILED {
			//告警通知
			content := fmt.Sprintf("【X%s告警X】 [%s] [%s] [%v]", DetectType, key, "line线路检测ip访问失败",
				time.Now().Format("2006-01-02 15:04:05"))
			alarmInfo[key] = content
			//发送MQ消息
			abnormal = append(abnormal, key)
		} else {
			mlog.Debugf("results key:%s, val:%v", key, val)
			continue
		}
	}

	//发送MQ消息
	go ProducerMq(LINE_DETECT, normal, abnormal)

	//发送告警通知
	go alarm.SendAlarmMsgs(alarmInfo)

	return nil
}

func checkIPResult(url string, m map[string]*scom.RegionResult, regionBound int) ResultValue {
	var pointBound int
	pointBound = regionBound/2 + 1
	ok := 0
	failed := 0
	for _, rr := range m {
		//mlog.Infof("xxxxxx [url:%s] [r:%s] [rr:%v]", url, r, rr)
		switch rv := checkLineResult(rr.LineResults, false); rv {
		case ResultValue_OK:
			ok++
		case ResultValue_FAILED:
			failed++
		}
	}
	if ok >= pointBound {
		return ResultValue_OK
	} else if failed >= regionBound {
		return ResultValue_FAILED
	} else {
		return ResultValue_UNKNOWN
	}
}

func (rsm *ResultsMan) GetLineResults(regionBound int) (mcr map[string]*CheckResult) {
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

			rv := checkIPResult(url, urlResult.Primarys, regionBound)
			addRS(&s.primary, rv)
		}
	}

	/* 0 (OK)  1 (FAILED)
	1 1 1 1 | 4 >= 3 | -1 (FAILED)
	0 1 0 1 | 2 <  3 | 0 （UNKNOWN）
	0 0 0 0 | 4 >= 3 | 1 (OK)
	1 1 1 0 | 3 >= 3 | 1 (OK)

	*/

	mcr = make(map[string]*CheckResult)
	check := func(rs *resultStat, url string) ResultValue {
		mlog.Debugf("ddddd 1 [url:%s] [rs:%v]\n", url, rs)
		if rs.ok >= rsm.keepPointsNum {
			return ResultValue_OK
		} else if rs.failed >= rsm.keepPointsNum {
			return ResultValue_FAILED
		} else {
			return ResultValue_UNKNOWN
		}
	}
	for url, s := range ss {
		cr := &CheckResult{}
		cr.Url = url
		cr.PrimaryResult = check(&s.primary, url)
		mcr[url] = cr
		mlog.Debugf("ddddd 2 [url:%s] [rs.P:%v]\n", url, cr.PrimaryResult)
	}

	return mcr
}

func compareLines(line1, line2 []scom.LineResult) bool {
	if len(line1) != len(line2) {
		return false
	}
	count := 0
	for _, l1 := range line1 {
		for _, l2 := range line2 {
			if reflect.DeepEqual(l1.Line, l2.Line) {
				if l1.Result.Code != l2.Result.Code {
					return false
				}
				if l1.Result.Status != l2.Result.Status {
					return false
				}
			} else {
				count++
			}
		}
	}
	if count > len(line1) { //如果线路属性不同(IP地址或营运商）
		return false
	}
	return true
}

func compareLinkRes(res1, res2 *UrlResult) bool {
	res1Primarys := res1.Primarys
	res2Primarys := res2.Primarys
	for region1, regionRes1 := range res1Primarys {
		if regionRes2, ok := res2Primarys[region1]; !ok {
			return false
		} else {
			return compareLines(regionRes1.LineResults, regionRes2.LineResults)
		}
	}
	return true
}

func (rr *ResultsRecords) CheckLineResults(pr *PointResult, mcr map[string]*CheckResult, cksLen int) (map[string]*CheckResult, error) {
	checkTime := pr.CheckTime
	createTime := time.Now()
	var infoRecords []interface{}
	alarmMap := make(map[string]*CheckResult)
	resultsMap := pr.Results
	vt1 := time.Now()
	mlog.Infof("1 ====== pr.map.len:%d, mcr.len:%d, rr.map.len:%d\n", len(resultsMap), len(mcr), len(rr.Results))

	updateRes := func(url string, mcrRes *CheckResult, rrRes CheckResult, res *UrlResult, isFirst bool) {
		// 添加到对比数据缓存中
		rr.Results[url] = *mcrRes
		rr.Results2[url] = res
		// 添加到入库数据列表中
		if !isFirst && mcrRes.PrimaryResult != ResultValue_UNKNOWN {
			infoRecords = append(infoRecords, &ResultRecords{Url: url, CheckTime: checkTime, CreateTime: createTime,
				PrePrimaryResult: rrRes.PrimaryResult, PreSecondResult: rrRes.SecondResult,
				NowPrimaryResult: mcrRes.PrimaryResult, NowSecondResult: mcrRes.SecondResult, Results: res})
			alarmMap[url] = mcrRes
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
			mlog.Debugf("[ip:%s] [mcr:%v]\n", url, mcr[url])
			if mcrRes, ok = mcr[url]; ok { // 找到，检查总结果
				cipr := checkIPResult(url, res1.Primarys, cksLen)
				mlog.Debugf("cvcv 1 [ip:%s] [rr.res:%d] [mcr.res:%d] [cipr.res:%d]\n", url, rrRes.PrimaryResult, mcrRes.PrimaryResult, cipr)
				if judgeResults(rrRes.PrimaryResult, mcrRes.PrimaryResult, cipr) {
					mlog.Debugf("cvcv 2 [ip:%s] [err.res:%d] [mcr.res:%d] [cipr.res:%d]\n", url, rrRes.PrimaryResult, mcrRes.PrimaryResult, cipr)
					updateRes(url, mcrRes, rrRes, res1, false)
				} else {
					updateRes(url, mcrRes, rrRes, res1, true)
				}
			}
		}
	}
	mlog.Infof("2 ====== add line records :%d, alarm.len:%d", len(infoRecords), len(alarmMap))
	rr.RUnlock()
	if len(infoRecords) == 0 {
		return nil, fmt.Errorf("data is nil")
	}
	vt2 := time.Now()
	mlog.Infof("1 cl=cl=cl=cl [exp:%v]\n", vt2.Sub(vt1))

	if _, err := mg.MongoDB.Collection("LineRecords").InsertMany(context.TODO(), infoRecords); err != nil {
		mlog.Errorf("add line results error: %v\n", err)
		return nil, err
	}

	vt3 := time.Now()
	mlog.Infof("2 cl=cl=cl=cl [exp:%v]\n", vt3.Sub(vt2))

	return alarmMap, nil
}
