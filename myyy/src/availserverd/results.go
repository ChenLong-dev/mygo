package main

import (
	"container/list"
	"context"
	"fmt"
	"github.com/ChenLong-dev/gobase/mg"
	"github.com/ChenLong-dev/gobase/mlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	logup "myyy/src/availserverd/saasyd_usability_log_up_grpc"
	"myyy/src/scom"
	"reflect"
	"strconv"
	"sync"
	"time"
)

const KeepPointNum = 10

const (
	RETRY        = 3 //重试次数
	EYERECORDS   = "EyeRecords"
	RETRYRECORDS = "RetryRecords"
)

type UrlResult struct {
	Url      string
	Primarys map[string]*scom.RegionResult
	Seconds  map[string]*scom.RegionResult
}
type PointResult struct {
	CheckTime time.Time
	Results   map[string]*UrlResult
}

type ResultValue int

const (
	ResultValue_OK      = 1  //	正常
	ResultValue_FAILED  = -1 //	故障
	ResultValue_UNKNOWN = 0  //	无法确认
)

type CheckResult struct {
	Url           string
	PrimaryResult ResultValue
	SecondResult  ResultValue
}

type ResultsMan struct {
	sync.RWMutex
	keepPointsNum int
	lastPoints    *list.List
	alarmUrlMap   map[string]int64
}

func (rsm *ResultsMan) KeepPointsNum() int {
	return rsm.keepPointsNum
}
func (rsm *ResultsMan) SetKeepPointsNum(nw int) {
	mlog.Infof("set new keepPointsNum value:%d", nw)
	rsm.keepPointsNum = nw
}
func (rsm *ResultsMan) PointCount() int {
	return rsm.lastPoints.Len()
}

func (rsm *ResultsMan) AddPoint(pr *PointResult, cksLen int) {
	mlog.Tracef("pr.Time=%v", pr.CheckTime)

	rsm.Lock()
	rsm.lastPoints.PushBack(pr)
	if rsm.lastPoints.Len() > rsm.keepPointsNum {
		rsm.lastPoints.Remove(rsm.lastPoints.Front())
	}
	rsm.Unlock()

	results := rsm.GetResults(cksLen)
	if len(results) <= 0 {
		mlog.Warnf("get result is nil, [pr.len:%d]\n", len(pr.Results))
		return
	}
	mlog.Tracef("GetResults len:%d, val:%+v", len(results), results)
	vt1 := time.Now()
	_, checkEyeResultsErr := resultRecords.CheckEyeResults(pr, results, cksLen)
	if checkEyeResultsErr != nil {
		mlog.Error(checkEyeResultsErr)
	}
	vt2 := time.Now()
	mlog.Infof("1 x=x=x=x [exp:%v]\n", vt2.Sub(vt1))
	if err := resultRecords.uploadUsabilityLog2(); err != nil {
		mlog.Error(err)
		return
	}
	//if err := defaultResultsMan.uploadUsabilityLog(r); err != nil {
	//	mlog.Error(err)
	//	return
	//}
	vt3 := time.Now()
	mlog.Infof("2 x=x=x=x [exp:%v]\n", vt3.Sub(vt2))
}

func (rsm *ResultsMan) uploadUsabilityLog(eyeLogResults []*EyeResultRecords) error {
	mlog.Debugf("===uploadUsabilityLog enter")
	var logupClient logup.UsabilityLogUpClient
	var conn *grpc.ClientConn
	for index := 0; index < 3; {
		var err error
		if etcdNode, ok := Mng.GetServiceInfoRandom(UpLogServicePath); ok {
			mlog.Debugf("etcdNode is %+v", etcdNode)
			conn, err = grpc.Dial(etcdNode.Info, grpc.WithInsecure())
			if err != nil {
				index++
				mlog.Errorf("GetServiceInfoRandom fail 1, retry index:%d, err:%+v", index, err)
				time.Sleep(time.Millisecond * 500)
				continue
			}
			logupClient = logup.NewUsabilityLogUpClient(conn)
		} else {
			index++
			mlog.Errorf("GetServiceInfoRandom fail 2, retry index:%d, no service is chosen::%+v", index, etcdNode)
			time.Sleep(time.Millisecond * 500)
			continue
		}
		break
	}
	defer conn.Close()

	sendGrpc := func(rr *EyeResultRecords) (*logup.UsabilityLogRsp, error) {
		mlog.Debugf("111 uploadUsabilityLog val:%+v", rr)
		var nodeList []*logup.UsabilityLogReq_NodeList
		//var nodeInfo logup.UsabilityLogReq_NodeList
		if rr.Results != nil {
			mlog.Debugf("222 uploadUsabilityLog val.Results:%+v", rr.Results)
			for _, regionResult := range rr.Results.Primarys {
				var nodeInfo logup.UsabilityLogReq_NodeList
				var lineResultArr []*logup.UsabilityLogReq_NodeList_LineResult
				nodeInfo.Node = regionResult.Region
				for _, lineVal := range regionResult.LineResults {
					lineResult := &logup.UsabilityLogReq_NodeList_LineResult{
						Addr:   lineVal.Line.Addr,
						Isp:    strconv.Itoa(int(lineVal.Line.Isp)),
						Code:   int32(lineVal.Result.Code),
						Status: lineVal.Result.Status,
					}
					lineResultArr = append(lineResultArr, lineResult)
				}
				nodeInfo.LineResult = lineResultArr
				nodeList = append(nodeList, &nodeInfo)
			}
			//nodeList = append(nodeList, &nodeInfo)
			mlog.Debugf("333 uploadUsabilityLog nodeList:%+v", nodeList)
		}
		var resStr string
		if int(rr.Result) == 1 {
			resStr = "recover"
		} else if int(rr.Result) == -1 {
			resStr = "disconnect"
		} else {
			//resStr = "unknown"
			return nil, fmt.Errorf("unknown")
		}

		var rsp *logup.UsabilityLogRsp
		var getLogupErr error
		rsp, getLogupErr = logupClient.UploadUsabilityLog(context.Background(),
			&logup.UsabilityLogReq{
				Url:        rr.Url,
				Result:     resStr,
				HappenTime: rr.CheckTime.Unix(),
				NodeList:   nodeList,
			})
		mlog.Debugf("444 uploadUsabilityLog client.UploadUsabilityLog rsp:%+v", rsp)
		return rsp, getLogupErr
	}

	//start the grpc call
	for _, val := range eyeLogResults {
		for index := 0; index < RETRY; {
			rsp, getLogupErr := sendGrpc(val)
			if getLogupErr != nil {
				index++
				mlog.Errorf("uploadUsabilityLog fail 1, retry index:%d, err:%+v", index, getLogupErr)
				time.Sleep(time.Millisecond * 500)
				continue
			}
			if rsp == nil {
				mlog.Errorf("getLogupErr:%+v", getLogupErr)
				break
			}
			if rsp.GetCode() != 0 {
				index++
				mlog.Errorf("uploadUsabilityLog fail 2, retry index:%d, err:%+v", index, getLogupErr)
				time.Sleep(time.Millisecond * 500)
				continue
			}
			break
		}
	}

	mlog.Debugf("===uploadUsabilityLog end")
	return nil
}

func (rsm *ResultsMan) GetLastPointResult() (pr *PointResult) {
	rsm.RLock()
	defer rsm.RUnlock()

	e := rsm.lastPoints.Back()
	if e == nil {
		return nil
	}
	return e.Value.(*PointResult)
}
func (rsm *ResultsMan) GetAllPointResult() (prs []*PointResult) {
	rsm.RLock()
	defer rsm.RUnlock()

	for e := rsm.lastPoints.Front(); e != nil; e = e.Next() {
		pr := e.Value.(*PointResult)
		prs = append(prs, pr)
	}
	return prs
}
func checkLineResult(lrs []scom.LineResult) ResultValue {
	failed := 0
	ok := 0
	for _, lr := range lrs {
		if lr.Result.Code >= 200 && lr.Result.Code <= 403 {
			ok++
		} else {
			failed++
		}
	}
	if ok > 0 {
		return ResultValue_OK
	} else {
		return ResultValue_FAILED
	}
}
func checkUrlResult(m map[string]*scom.RegionResult, regionBound int) ResultValue {
	var pointBound int
	pointBound = regionBound/2 + 1
	ok := 0
	failed := 0
	for _, rr := range m {
		switch rv := checkLineResult(rr.LineResults); rv {
		case ResultValue_OK:
			ok++
		case ResultValue_FAILED:
			failed++
		}
	}
	mlog.Debugf("aaaa [regionBound:%d] [pointBound:%d] [ok:%d] [failed:%d]", regionBound, pointBound, ok, failed)
	if ok >= pointBound { // 成功：需要过检测节点数的一半
		return ResultValue_OK
	} else if failed >= regionBound { // 故障：需要所有检测节点都检测故障
		return ResultValue_FAILED
	} else {
		return ResultValue_UNKNOWN
	}
}
func (rsm *ResultsMan) GetResults(regionBound int) (mcr map[string]*CheckResult) {
	//var pointBound int
	//pointBound = regionBound/2 + 1
	mlog.Infof("keepPointsNum=%d,regionBound=%d", rsm.keepPointsNum, regionBound)

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

			rv := checkUrlResult(urlResult.Primarys, regionBound)
			//mlog.Debugf("ccccc 1 [url:%s] [rv:%d]", url, rv)
			addRS(&s.primary, rv)
			mlog.Debugf("ccccc 1 [url:%s] [s:%d]", url, s.primary)
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

func NewResultsMan(keepPointsNum int) *ResultsMan {
	rsm := &ResultsMan{keepPointsNum: keepPointsNum,
		alarmUrlMap: make(map[string]int64)}
	rsm.lastPoints = list.New()
	return rsm
}

var defaultResultsMan *ResultsMan

var resultRecords *EyeResultsRecords

func InitResultsMan(keepPointNum int) {
	defaultResultsMan = NewResultsMan(keepPointNum)
	resultRecords = &EyeResultsRecords{Results: make(map[string]CheckResult), Results2: make(map[string]*UrlResult)}
	resultRecords.SendMsgList = list.New()
}

type EyeResultsRecords struct {
	sync.RWMutex
	Results     map[string]CheckResult
	Results2    map[string]*UrlResult
	SendMsgList *list.List
}

type EyeResultRecords struct {
	Url        string      `bson:"url"`
	CheckTime  time.Time   `bson:"check_time"`
	CreateTime time.Time   `bson:"create_time"`
	Result     ResultValue `bson:"result"`
	Results    *UrlResult  `bson:"results"`
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

func compareRes(res1, res2 *UrlResult) bool {
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

func (errs EyeResultsRecords) CheckEyeResults(pr *PointResult, mcr map[string]*CheckResult, cksLen int) ([]*EyeResultRecords, error) {
	checkTime := pr.CheckTime
	createTime := time.Now()
	var infoRecords []*EyeResultRecords
	resultsMap := pr.Results
	vt1 := time.Now()
	mlog.Infof("1 ====== pr.map.len:%d, mcr.len:%d, rr.map.len:%d\n", len(resultsMap), len(mcr), len(errs.Results))

	updateRes := func(url string, mcrRes *CheckResult, errRes CheckResult, res *UrlResult, first bool) {
		// 添加到对比数据缓存中
		errs.Results[url] = *mcrRes
		errs.Results2[url] = res
		if !first && mcrRes.PrimaryResult != ResultValue_UNKNOWN { // 只有在判定结果不为0未知的情况下进行更新
			// 添加到入库数据列表中
			infoRecords = append(infoRecords, &EyeResultRecords{Url: url, CheckTime: checkTime, CreateTime: createTime,
				Result: mcrRes.PrimaryResult, Results: res})
			errs.SendMsgList.PushBack(&SendMsg{Url: url, CheckTime: checkTime, Result: mcrRes.PrimaryResult, Results: res})
		}
	}

	mlog.Debugf("yyy 1 cache: [%v]\n", errs.Results)
	errs.RLock()
	for url, res1 := range resultsMap {
		var mcrRes *CheckResult
		if errRes, ok := errs.Results[url]; !ok {
			if mcrRes, ok = mcr[url]; ok {
				updateRes(url, mcrRes, errRes, res1, true)
			} else {
				mlog.Warnf("1 [url:%s] not result in mcr!", url)
			}
		} else {
			if mcrRes, ok = mcr[url]; ok { // 找到，检查总结果
				cur := checkUrlResult(res1.Primarys, cksLen)
				mlog.Debugf("bbbbb 1 [url:%s] [err.res:%d] [mcr.res:%d] [cur.res:%d]\n", url, errRes.PrimaryResult, mcrRes.PrimaryResult, cur)
				if (errRes.PrimaryResult == ResultValue_UNKNOWN && mcrRes.PrimaryResult == ResultValue_FAILED) || // 之前（缓存）不是故障，判断两次结果后都是故障，上报故障；
					(errRes.PrimaryResult != ResultValue_OK && cur == ResultValue_OK) { // 之前（缓存）是故障，当次检测结果是正常，上报恢复故障
					mlog.Debugf("bbbbb 2 [url:%s] [err.res:%d] [mcr.res:%d] [cur.res:%d]\n", url, errRes.PrimaryResult, mcrRes.PrimaryResult, cur)
					updateRes(url, mcrRes, errRes, res1, false)
				} else {
					updateRes(url, mcrRes, errRes, res1, true)
				}
			}
		}
	}
	mlog.Debugf("yyy 2 cache: [%v]\n", errs.Results)
	mlog.Infof("2 ====== add eye records [info.len:%d] [sendmsg.len:%d] ", len(infoRecords), errs.SendMsgList.Len())
	errs.RUnlock()

	if len(infoRecords) == 0 {
		return nil, fmt.Errorf("data is nil")
	}
	// 添加到发送队列中

	vt2 := time.Now()
	mlog.Infof("1 cl=cl=cl=cl [exp:%v]\n", vt2.Sub(vt1))
	var addInfoRecords []interface{}
	for _, ir := range infoRecords {
		addInfoRecords = append(addInfoRecords, ir)
	}
	if _, err := mg.MongoDB.Collection(EYERECORDS).InsertMany(context.TODO(), addInfoRecords); err != nil {
		mlog.Errorf("add eye results error: %v\n", err)
		return nil, err
	}
	vt3 := time.Now()
	mlog.Infof("2 cl=cl=cl=cl [exp:%v]\n", vt3.Sub(vt2))

	return infoRecords, nil
}

type SendMsg struct {
	Id        primitive.ObjectID `bson:"_id"`
	Url       string             `bson:"url"`
	CheckTime time.Time          `bson:"check_time"`
	Result    ResultValue        `bson:"result"`
	Results   *UrlResult         `bson:"results"`
}

func (errs *EyeResultsRecords) uploadUsabilityLog2() error {
	mlog.Debugf("===uploadUsabilityLog enter")
	//vt1 := time.Now()
	//defer func() {
	//	vt2 := time.Now()
	//	mlog.Infof("upload log [exp:%v]", vt2.Sub(vt1))
	//}()

	_ = errs.retryRecordsSelect()
	//if err := errs.retryRecordsSelect(); err != nil {
	//	mlog.Warnf("get retry records is failed! [err:%v]", err)
	//}
	if errs.SendMsgList.Len() == 0 {
		return nil
	}
	mlog.Debugf("get send msg [len:%d]", errs.SendMsgList.Len())

	// 判断是否能够连接数据中心
	var logupClient logup.UsabilityLogUpClient
	var conn *grpc.ClientConn
	for index := 0; index < 3; {
		var err error
		if etcdNode, ok := Mng.GetServiceInfoRandom(UpLogServicePath); ok {
			mlog.Debugf("etcdNode is %+v", etcdNode)
			conn, err = grpc.Dial(etcdNode.Info, grpc.WithInsecure())
			if err != nil {
				index++
				mlog.Errorf("GetServiceInfoRandom fail 1, retry index:%d, err:%+v", index, err)
				time.Sleep(time.Millisecond * 500)
				continue
			}
			logupClient = logup.NewUsabilityLogUpClient(conn)
		} else {
			index++
			mlog.Errorf("GetServiceInfoRandom fail 2, retry index:%d, no service is chosen::%+v", index, etcdNode)
			time.Sleep(time.Millisecond * 500)
			continue
		}
		break
	}
	defer conn.Close()

	mlog.Debugf("etcd init success ...")
	var SendMsgList []*SendMsg
	errs.RLock()
	var n *list.Element
	for e := errs.SendMsgList.Front(); e != nil; e = n {
		n = e.Next()
		var sm *SendMsg
		sm = e.Value.(*SendMsg)
		SendMsgList = append(SendMsgList, sm)
		errs.SendMsgList.Remove(e)
	}
	errs.RUnlock()

	if !reflect.ValueOf(logupClient).IsValid() {
		mlog.Warnf("get a no valid logupClient")
		for _, smval := range SendMsgList {
			err := processRetryRecords(false, smval)
			if err != nil {
				mlog.Error(err)
			}
		}
		return fmt.Errorf("get a nil logupClient, all records into RetryRecords [sm.len:%d]", len(SendMsgList))
	}

	sendGrpc := func(sm *SendMsg) (*logup.UsabilityLogRsp, error) {
		mlog.Debugf("111 uploadUsabilityLog Results:%+v", sm)
		var nodeList []*logup.UsabilityLogReq_NodeList
		for _, regionResult := range sm.Results.Primarys {
			var nodeInfo logup.UsabilityLogReq_NodeList
			var lineResultArr []*logup.UsabilityLogReq_NodeList_LineResult
			nodeInfo.Node = regionResult.Region
			for _, lineVal := range regionResult.LineResults {
				lineResult := &logup.UsabilityLogReq_NodeList_LineResult{
					Addr:   lineVal.Line.Addr,
					Isp:    strconv.Itoa(int(lineVal.Line.Isp)),
					Code:   int32(lineVal.Result.Code),
					Status: lineVal.Result.Status,
				}
				lineResultArr = append(lineResultArr, lineResult)
			}
			nodeInfo.LineResult = lineResultArr
			nodeList = append(nodeList, &nodeInfo)
			mlog.Debugf("222 uploadUsabilityLog nodeList:%+v", nodeList)
		}

		var resStr string
		if int(sm.Result) == 1 {
			resStr = "recover"
		} else if int(sm.Result) == -1 {
			resStr = "disconnect"
		} else {
			//resStr = "unknown"
			return nil, fmt.Errorf("unknown")
		}

		var rsp *logup.UsabilityLogRsp
		var getLogupErr error
		rsp, getLogupErr = logupClient.UploadUsabilityLog(context.Background(),
			&logup.UsabilityLogReq{
				Url:        sm.Url,
				Result:     resStr,
				HappenTime: sm.CheckTime.Unix(),
				NodeList:   nodeList,
			})
		mlog.Debugf("uploadUsabilityLog [rsp:%+v] [rsp.code:%d] [rsp:%s]", rsp, rsp.GetCode(), rsp.GetMessage())
		return rsp, getLogupErr
	}

	for _, smval := range SendMsgList {
		flag := false
		for index := 0; index < RETRY; {
			rsp, getLogupErr := sendGrpc(smval)
			if getLogupErr != nil {
				index++
				mlog.Errorf("uploadUsabilityLog fail 1, retry index:%d, err:%+v, [url:%s]", index, getLogupErr, smval.Url)
				time.Sleep(time.Millisecond * 500)
				continue
			}
			if rsp == nil {
				mlog.Errorf("getLogupErr:%+v", getLogupErr)
				flag = true
				break
			}
			if rsp.GetCode() != 0 {
				index++
				mlog.Errorf("uploadUsabilityLog fail 2, retry index:%d, err:%+v, [url:%s]", index, getLogupErr, smval.Url)
				time.Sleep(time.Millisecond * 500)
				continue
			}
			flag = true
			break
		}
		mlog.Infof("send result [flag:%t] [url:%s] [iszero:%v] [id:%v]", flag, smval.Url, smval.Id.IsZero(), smval.Id)
		err := processRetryRecords(flag, smval)
		if err != nil {
			mlog.Warn(err)
		}
	}

	mlog.Debugf("===uploadUsabilityLog end")
	return nil
}

type RetryRecords struct {
	Id         primitive.ObjectID `bson:"_id"`
	Url        string             `bson:"url`
	SendMsg    SendMsg            `bson:"send_msg"`
	CheckTime  time.Time          `bson:"check_time"`
	CreateTime time.Time          `bson:"create_time"`
	UpdateTime time.Time          `bson:"update_time"`
	SendFlag   bool               `bson:"send_flag"`
}

func processRetryRecords(flag bool, smval *SendMsg) error {
	if !flag { // 发送失败，存入数据库，等待再次发送
		// TODO：存入数据库中
		mlog.Debugf("rrrrr 1 insert [url:%s] [smval:%v] [iszero:%v]", smval.Url, smval, smval.Id.IsZero())
		if smval.Id.IsZero() == true {
			if err := retryRecordsInsert(smval); err != nil {
				mlog.Warnf("insert retry records is failed! [err:%v]", err)
				return err
			}
		}
	} else { // 发送成功，查看是否是重发数据，如果是更新数据库状态，如果不是直接返回
		// TODO：更新数据库状态
		mlog.Debugf("rrrrr 2 update [url:%s] [smval:%v] [iszero:%v]", smval.Url, smval, smval.Id.IsZero())
		if smval.Id.IsZero() == false {
			if err := retryRecordsUpdate(smval); err != nil {
				mlog.Warnf("update retry records is failed! [err:%v]", err)
				return err
			}
		}
	}
	return nil
}

func retryRecordsInsert(sm *SendMsg) error {
	mlog.Debugf("insert a retry records [sm:%v]", *sm)
	now := time.Now()
	rri := RetryRecords{Id: primitive.NewObjectID(), Url: sm.Url, SendMsg: *sm, CheckTime: sm.CheckTime, SendFlag: false, CreateTime: now, UpdateTime: now}
	if _, err := mg.MongoDB.Collection(RETRYRECORDS).InsertOne(context.TODO(), rri); err != nil {
		mlog.Warnf("add retry records error: %v\n", err)
		return err
	}
	return nil
}

func retryRecordsUpdate(sm *SendMsg) error {
	mlog.Debugf("update a retry records [sm:%v]", *sm)
	now := time.Now()
	filter := bson.D{{"_id", sm.Id}}
	update := bson.D{{"$set", bson.D{{"send_flag", true}, {"update_time", now}}}}
	if _, err := mg.MongoDB.Collection(RETRYRECORDS).UpdateOne(context.TODO(), filter, update); err != nil {
		mlog.Errorf("update retry records error: %v\n", err)
		return err
	}
	return nil
}

func (errs *EyeResultsRecords) retryRecordsSelect() error {
	mlog.Debug("select 500 retry records")
	collection := mg.CreateCollection(RETRYRECORDS)
	filter := bson.D{{"send_flag", false}}
	cursor, err := mg.Find(collection, filter, 500, 0)
	if err != nil {
		mlog.Warn(err)
		return err
	}
	defer cursor.Close(context.Background())
	if err = cursor.Err(); err != nil {
		mlog.Warn(err)
		return err
	}
	errs.Lock()
	len1 := errs.SendMsgList.Len()
	for cursor.Next(context.Background()) {
		var rr *RetryRecords
		if err = cursor.Decode(&rr); err != nil {
			mlog.Warn(err)
			return err
		}
		sm := rr.SendMsg
		sm.Id = rr.Id
		errs.SendMsgList.PushBack(&sm)
		mlog.Debug("rrrrr 2 [id:%v] [rr:%v] [sm:%t]", rr.Id, rr, &sm)
	}
	len2 := errs.SendMsgList.Len()
	errs.Unlock()
	mlog.Infof("add retry records to list, [len1:%d] [len2:%d] [add.len:%d]", len1, len2, len2-len1)
	return nil
}