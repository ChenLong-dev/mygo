/**
* @Author: cl
* @Date: 2021/1/16 11:18
 */
package main


import (
	"container/list"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyd/src/scom"
	"sync"
	"time"
)

const KeepPointNum = 10

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
	TargetPointNum      = 3
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

func NewResultsMan(keepPointsNum int) *ResultsMan {
	rsm := &ResultsMan{keepPointsNum: keepPointsNum,
		alarmUrlMap: make(map[string]int64)}
	rsm.lastPoints = list.New()
	return rsm
}

var defaultResultsMan *ResultsMan

var targetResultsMan *ResultsMan

var resultRecords *ResultsRecords

var targetResultRecords *ResultsTargetRecords

var aclResultsRecords *AclResultsRecords

func InitResultsMan(keepPointNum int) {
	defaultResultsMan = NewResultsMan(keepPointNum)
	targetResultsMan = NewResultsMan(TargetPointNum)
	resultRecords = &ResultsRecords{Results: make(map[string]CheckResult), Results2: make(map[string]*UrlResult)}
	targetResultRecords = &ResultsTargetRecords{Results: make(map[string]CheckResult),
		Results2: make(map[string]*UrlResult), Results3: make(map[string]*AddrTargetResult)}
	aclResultsRecords = NewAclResultsMan()
}

type ResultsRecords struct {
	sync.RWMutex
	Results  map[string]CheckResult
	Results2 map[string]*UrlResult
}

type ResultRecords struct {
	Url              string      `bson:"url"`
	CheckTime        time.Time   `bson:"check_time"`
	CreateTime       time.Time   `bson:"create_time"`
	PrePrimaryResult ResultValue `bson:"pre_primary_result"`
	NowPrimaryResult ResultValue `bson:"now_primary_result"`
	PreSecondResult  ResultValue `bson:"pre_second_result"`
	NowSecondResult  ResultValue `bson:"now_second_result"`
	Results          *UrlResult  `bson:"results"`
}