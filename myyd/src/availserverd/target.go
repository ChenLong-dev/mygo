/**
* @Author: cl
* @Date: 2021/1/16 11:19
 */
package main

import (
	"fmt"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyd/src/availserverd/alarm"
	"myyd/src/availserverd/config"
	"myyd/src/scom"
	"strings"
	"sync"
	"time"
)

type ResultsTargetRecords struct {
	sync.RWMutex
	Results  map[string]CheckResult
	Results2 map[string]*UrlResult
	Results3 map[string]*AddrTargetResult
}

type TargetResultRecords struct {
	Url        string      `bson:"url"`
	CheckTime  time.Time   `bson:"check_time"`
	CreateTime time.Time   `bson:"create_time"`
	Result     ResultValue `bson:"result"`
	Results    *UrlResult  `bson:"results"`
}

type AddrTargetResult struct {
	NodeName string
	Result   ResultValue
	Line     scom.Line
}

func (rsm *ResultsMan) AddTargetPoint(tr *PointResult, cksLen int, allClusterIPMap map[string][]scom.Line) {
	mlog.Tracef("tr.Time=%v", tr.CheckTime)

	rsm.Lock()
	rsm.lastPoints.PushBack(tr)
	if rsm.lastPoints.Len() > rsm.keepPointsNum {
		rsm.lastPoints.Remove(rsm.lastPoints.Front())
	}
	rsm.Unlock()

	ipnnMap := make(map[string]*AddrTargetResult)
	for nn, ll := range allClusterIPMap {
		for _, line := range ll {
			ipnnMap[line.Addr] = &AddrTargetResult{NodeName: nn, Line: line, Result: 99}
		}
	}

	results, addresults := rsm.GetTargetResults(cksLen, ipnnMap)
	if len(results) <= 0 {
		mlog.Warnf("get result is nil, [tr.len:%d]\n", len(tr.Results))
		return
	}
	mlog.Debugf("GetResults [results.len:%d], [addresults.len:%d]", len(results), len(addresults))
	vt1 := time.Now()
	// 机房上报（正常（normal）、故障（abnormal））
	_ = DetectTargetClusterReport(results)
	// 线路上报（正常（normal）、故障（abnormal））
	_ = DetectTargetLineReport(addresults)
	// 线路数据处理/线路告警
	_ = targetResultRecords.CheckAddrResults(tr, addresults, cksLen)
	// 机房数据处理/机房告警
	_, checkTargetResultsErr := targetResultRecords.CheckTargetResults(tr, results, cksLen)
	if checkTargetResultsErr != nil {
		mlog.Error(checkTargetResultsErr)
	}
	vt2 := time.Now()
	mlog.Infof("1 x=x=x=x [exp:%v]\n", vt2.Sub(vt1))
}

func checkTargetLineResult(region string, lrs []scom.LineResult) ResultValue {
	failed := 0
	ok := 0
	for _, lr := range lrs {
		if lr.Result.Code >= 200 && lr.Result.Code <= 500 {
			ok++
		} else {
			failed++
		}
	}
	mlog.Debugf("lllll [region:%s] [ll:%v] [ok:%d] [failed:%d]", region, lrs, ok, failed)
	if ok > len(lrs)/2 {
		return ResultValue_OK
	} else {
		return ResultValue_FAILED
	}
}

func checkTargetUrlResult(m map[string]*scom.RegionResult, regionBound int) ResultValue {
	var pointBound int
	pointBound = regionBound/2 + 1
	ok := 0
	failed := 0
	for r, rr := range m {
		switch rv := checkTargetLineResult(r, rr.LineResults); rv {
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

func checkTargetAddrResult(m map[string]*scom.RegionResult, regionBound int) map[string]ResultValue {
	var pointBound int
	pointBound = regionBound/2 + 1
	type resultStat struct {
		ok      int
		failed  int
		unknown int
	}
	addrResMap := make(map[string]*resultStat)
	for _, rr := range m {
		for _, lr := range rr.LineResults {
			var lrst *resultStat
			if lineRst, ok := addrResMap[lr.Line.Addr]; ok {
				lrst = lineRst
			} else {
				lrst = &resultStat{}
				addrResMap[lr.Line.Addr] = lrst
			}
			if lr.Result.Code >= 200 && lr.Result.Code < 500 {
				lrst.ok++
			} else {
				lrst.failed++
			}
		}
	}
	addrRests := make(map[string]ResultValue)
	for addr, rs := range addrResMap {
		if rs.ok >= pointBound { // 成功：需要过检测节点数的一半
			addrRests[addr] = ResultValue_OK
		} else if rs.failed >= regionBound { // 故障：需要所有检测节点都检测故障
			addrRests[addr] = ResultValue_FAILED
		} else {
			addrRests[addr] = ResultValue_UNKNOWN
		}
	}

	return addrRests
}

func (rsm *ResultsMan) GetTargetResults(regionBound int, ipnnMap map[string]*AddrTargetResult) (mcr map[string]*CheckResult, atr map[string]*AddrTargetResult) {
	pointBound := rsm.keepPointsNum/2 + 1
	mlog.Infof("keepPointsNum=%d,pointBound=%d,regionBound=%d", rsm.keepPointsNum, pointBound, regionBound)

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
	ar := make(map[string]*resultStat)
	for e := rsm.lastPoints.Front(); e != nil; e = e.Next() {
		tr := e.Value.(*PointResult)
		for url, urlResult := range tr.Results {
			// 统计机房结果
			var s *Stat
			if r, ok := ss[url]; ok {
				s = r
			} else {
				s = &Stat{}
				ss[url] = s
			}

			tur := checkTargetUrlResult(urlResult.Primarys, regionBound)
			mlog.Debugf("ccccc 1 [url:%s] [rv:%d]", url, tur)
			addRS(&s.primary, tur)
			mlog.Debugf("ccccc 1 [url:%s] [s:%d]", url, s.primary)

			// 统计线路结果
			tar := checkTargetAddrResult(urlResult.Primarys, regionBound)
			for addr, tr := range tar {
				var a *resultStat
				if w, ok := ar[addr]; ok {
					a = w
				} else {
					a = &resultStat{}
					ar[addr] = a
				}
				addRS(a, tr)
			}
		}
	}

	// 机房统计
	mcr = make(map[string]*CheckResult)
	check := func(rs *resultStat, url string) ResultValue {
		mlog.Debugf("ddddd 1 [url:%s] [rs:%v] [pointBound:%d]\n", url, rs, pointBound)
		if rs.ok >= pointBound { // 成功: 一半以上检测点成功为成功
			return ResultValue_OK
		} else if rs.failed >= regionBound { // 故障： 全部检测点都故障为故障
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
	// 线路统计
	atr = make(map[string]*AddrTargetResult)
	for addr, lst := range ar {
		ar := &AddrTargetResult{}
		ar.Result = check(lst, addr)
		if ars, ok := ipnnMap[addr]; ok {
			ar.NodeName = ars.NodeName
			ar.Line = ars.Line
		}
		atr[addr] = ar
		mlog.Debugf("ddddd 2 [addr:%s] [ar:%v]\n", addr, ar)
	}
	return mcr, atr
}

func judgeResults(previous, present, current ResultValue) bool {
	/*
		故障：1 0 -1
		     0 -1 -1
		恢复：-1 0 1
		      0 1 1
	*/
	if (previous == ResultValue_UNKNOWN && present == ResultValue_FAILED && current == ResultValue_FAILED) || // 3次故障算故障
		(previous == ResultValue_FAILED && present == ResultValue_UNKNOWN && current == ResultValue_OK) { // 1次成功算成功
		return true
	}
	return false
}

func (rtr *ResultsTargetRecords) CheckAddrResults(tr *PointResult, addresults map[string]*AddrTargetResult, cksLen int) []*AddrTargetResult {
	mlog.Debug("fffff ENTER CheckAddrResults =============")
	var infoRecords []*AddrTargetResult
	updateRes := func(addr string, atr *AddrTargetResult, first bool) {
		// 添加到对比数据缓存中
		rtr.Results3[addr] = atr
		if !first { // 只有在判定结果不为0未知的情况下进行更新
			// 添加到入库数据列表中
			infoRecords = append(infoRecords, atr)
		}
	}

	curMap := make(map[string]ResultValue)
	for _, res := range tr.Results {
		tar := checkTargetAddrResult(res.Primarys, cksLen)
		for addr, tr := range tar {
			curMap[addr] = tr
		}
	}

	rtr.RLock()
	for addr, atr := range addresults {
		//mlog.Debugf("fffff 1 [addr:%s] [atr:%v]\n", addr, atr)
		if r, ok := rtr.Results3[addr]; !ok {
			updateRes(addr, atr, true)
		} else {
			var cur ResultValue
			if cur, ok = curMap[addr]; !ok {
				continue
			}
			mlog.Debugf("fffff 1 [addr:%s] [atr.region:%s] [r.res:%d] [atr.res:%d] [cur:%d]\n", addr, atr.NodeName, r.Result, atr.Result, cur)
			if judgeResults(r.Result, atr.Result, cur) {
				mlog.Debugf("fffff 2 [addr:%s] [atr.region:%s] [r.res:%d] [atr.res:%d] [cur:%d]\n", addr, atr.NodeName, r.Result, atr.Result, cur)
				updateRes(addr, atr, false)
			} else {
				updateRes(addr, atr, true)
			}
		}
	}
	mlog.Infof("3 ====== add addr records :%d", len(infoRecords))
	rtr.RUnlock()

	if len(infoRecords) == 0 {
		return infoRecords
	}

	getisp := func(line scom.Line) string {
		if line.Isp == 1 {
			return "电信"
		} else if line.Isp == 2 {
			return "联通"
		} else if line.Isp == 3 {
			return "移动"
		} else {
			return "default"
		}
	}

	// 线路告警
	alarmInfo := make(map[string]string)
	for _, at := range infoRecords {
		var content string
		if at.Result == ResultValue_UNKNOWN {
			content = fmt.Sprintf("【Ytarget线路恢复Y】[%s] [%s] [%s] [%+v]",
				getisp(at.Line), at.Line.Addr, "故障线路恢复正常", time.Now().Format("2006-01-02 15:04:05"))
		} else if at.Result == ResultValue_FAILED {
			content = fmt.Sprintf("【Xtarget线路告警X】[%s] [%s] [%s] [%+v]",
				getisp(at.Line), at.Line.Addr, "线路故障", time.Now().Format("2006-01-02 15:04:05"))
		} else {
			mlog.Debugf("results [addr:%s]", at.Line.Addr)
			continue
		}
		alarmInfo[at.Line.Addr] = content
	}

	//发送线路告警通知
	if len(alarmInfo) > 0 {
		go alarm.SendAlarmMsgs(alarmInfo)
	}

	return infoRecords
}

func ClusterAlarm(trr []*TargetResultRecords) {
	alarmInfo := make(map[string]string)
	for _, tr := range trr {
		var content string
		nodeName := getCluster(tr.Url)
		if nodeName == "" {
			continue
		}
		if tr.Result == ResultValue_UNKNOWN {
			// 告警信息
			content = fmt.Sprintf("【Ytarget机房恢复Y】[%s] [%s] [%+v]",
				nodeName, "故障机房恢复正常", time.Now().Format("2006-01-02 15:04:05"))
		} else if tr.Result == ResultValue_FAILED {
			// 告警信息
			content = fmt.Sprintf("【Xtarget机房告警X】[%s] [%s] [%+v]",
				nodeName, "机房故障", time.Now().Format("2006-01-02 15:04:05"))
		} else {
			mlog.Debugf("results [url:%s] [nn:%s] [rs:%v]", tr.Url, nodeName, tr)
			continue
		}
		alarmInfo[nodeName] = content
	}
	//发送告警通知
	if len(alarmInfo) > 0 {
		go alarm.SendAlarmMsgs(alarmInfo)
	}
}

func (rtr *ResultsTargetRecords) CheckTargetResults(tr *PointResult, mcr map[string]*CheckResult, cksLen int) ([]*TargetResultRecords, error) {
	checkTime := tr.CheckTime
	createTime := time.Now()
	var infoRecords []*TargetResultRecords
	resultsMap := tr.Results
	vt1 := time.Now()
	mlog.Infof("1 ====== tr.map.len:%d, mcr.len:%d, rr.map.len:%d\n", len(resultsMap), len(mcr), len(rtr.Results))

	updateRes := func(url string, mcrRes *CheckResult, errRes CheckResult, res *UrlResult, first bool) {
		// 添加到对比数据缓存中
		rtr.Results[url] = *mcrRes
		rtr.Results2[url] = res
		if !first { // 只有在判定结果不为0未知的情况下进行更新
			// 添加到入库数据列表中
			infoRecords = append(infoRecords, &TargetResultRecords{Url: url, CheckTime: checkTime, CreateTime: createTime,
				Result: mcrRes.PrimaryResult, Results: res})
		}
	}

	mlog.Debugf("yyy 1 cache: [%v]\n", rtr.Results)
	rtr.RLock()
	for url, res1 := range resultsMap {
		var mcrRes *CheckResult
		if errRes, ok := rtr.Results[url]; !ok {
			if mcrRes, ok = mcr[url]; ok {
				updateRes(url, mcrRes, errRes, res1, true)
			} else {
				mlog.Warnf("1 [url:%s] not result in mcr!", url)
			}
		} else {
			if mcrRes, ok = mcr[url]; ok { // 找到，检查总结果
				cur := checkTargetUrlResult(res1.Primarys, cksLen)
				mlog.Debugf("bbbbb 1 [url:%s] [err.res:%d] [mcr.res:%d] [cur.res:%d]\n", url, errRes.PrimaryResult, mcrRes.PrimaryResult, cur)
				if judgeResults(errRes.PrimaryResult, mcrRes.PrimaryResult, cur) {
					mlog.Infof("bbbbb 2 [url:%s] [err.res:%d] [mcr.res:%d] [cur.res:%d]\n", url, errRes.PrimaryResult, mcrRes.PrimaryResult, cur)
					updateRes(url, mcrRes, errRes, res1, false)
				} else {
					updateRes(url, mcrRes, errRes, res1, true)
				}
			}
		}
	}
	mlog.Debugf("yyy 2 cache: [%v]\n", rtr.Results)
	mlog.Infof("2 ====== add target records :%d", len(infoRecords))
	rtr.RUnlock()

	if len(infoRecords) == 0 {
		return nil, fmt.Errorf("data is nil")
	}

	// 机房告警
	ClusterAlarm(infoRecords)

	vt2 := time.Now()
	mlog.Infof("1 cl=cl=cl=cl [exp:%v]\n", vt2.Sub(vt1))

	// 机房结果保存数据库
	/*var addInfoRecords []interface{}
	for _, ir := range infoRecords {
		addInfoRecords = append(addInfoRecords, ir)
	}
	if _, err := mg.MongoDB.Collection("TargetRecords").InsertMany(context.TODO(), addInfoRecords); err != nil {
		mlog.Errorf("add eye results error: %v\n", err)
		return nil, err
	}
	vt3 := time.Now()
	mlog.Infof("2 cl=cl=cl=cl [exp:%v]\n", vt3.Sub(vt2))*/

	return infoRecords, nil
}

func DetectTargetLineReport(results map[string]*AddrTargetResult) error {
	vt1 := time.Now()
	var normal, abnormal []string
	for addr, ar := range results {
		if ar.Result == ResultValue_OK {
			//发送MQ消息
			normal = append(normal, addr)
		} else if ar.Result == ResultValue_FAILED {
			//发送MQ消息
			abnormal = append(abnormal, addr)
		} else {
			mlog.Debugf("results addr:%s, ar:%v", addr, ar)
			continue
		}
	}
	vt2 := time.Now()
	mlog.Infof("1 v=v=v=v [exp:%v]\n", vt2.Sub(vt1))

	//发送MQ消息
	mlog.Debugf("DetectTargetAddrReport len(normal):%v, len(abnormal):%v", len(normal), len(abnormal))
	if len(normal)+len(abnormal) > 0 {
		go ProducerMq("target_line", normal, abnormal)
	}
	vt3 := time.Now()
	mlog.Infof("2 v=v=v=v [exp:%v]\n", vt3.Sub(vt2))
	return nil
}

func getCluster(url string) string {
	// http://bj.detect.sre.ac.cn
	if !strings.Contains(url, `.detect.sre.ac.cn`) {
		return ""
	}
	arr1 := strings.Split(url, "//")
	if len(arr1) < 2 {
		return ""
	}
	arr2 := strings.Split(arr1[1], ".")
	return "cluster_node_" + arr2[0]
}

func DetectTargetClusterReport(results map[string]*CheckResult) error {
	vt1 := time.Now()
	var normal, abnormal []string
	for url, rs := range results {
		nodeName := getCluster(url)
		if nodeName == "" {
			continue
		}
		if rs.PrimaryResult == ResultValue_OK {
			//发送MQ消息
			normal = append(normal, nodeName)
		} else if rs.PrimaryResult == ResultValue_FAILED {
			//发送MQ消息
			abnormal = append(abnormal, nodeName)
		} else {
			mlog.Debugf("results url:%s, rs:%v", url, rs)
			continue
		}
	}
	vt2 := time.Now()
	mlog.Infof("1 v=v=v=v [exp:%v]\n", vt2.Sub(vt1))
	if len(abnormal) == len(config.Conf.AvailConf.TargetUrl) { // 现在的靶机为单机，当所有靶机都访问不通的情况下可能是靶机故障，不回源，只告警
		mlog.Warnf("maybe target web is abnormal [normal.len:%d] [abnormal.len:%d] [target.len:%d]", len(normal),
			len(abnormal), len(config.Conf.AvailConf.TargetUrl))
		return nil
	}
	//发送MQ消息
	mlog.Debugf("DetectTargetClusterReport len(normal):%v, len(abnormal):%v", len(normal), len(abnormal))
	if len(normal)+len(abnormal) > 0 {
		go ProducerMq("target_cluster", normal, abnormal)
	}
	vt3 := time.Now()
	mlog.Infof("2 v=v=v=v [exp:%v]\n", vt3.Sub(vt2))
	return nil
}

func (rtr *ResultsTargetRecords) GetTargetResults() (map[string]CheckResult, error) {
	rtr.RLock()
	defer rtr.RUnlock()
	results := rtr.Results
	if len(results) == 0 {
		return nil, fmt.Errorf("results is nil")
	}
	return results, nil
}

func (rtr *ResultsTargetRecords) GetAddrResults() (map[string]*AddrTargetResult, error) {
	rtr.RLock()
	defer rtr.RUnlock()
	results := rtr.Results3
	if len(results) == 0 {
		return nil, fmt.Errorf("results is nil")
	}
	return results, nil
}
