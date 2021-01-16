package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase/msys"
	"github.com/ChenLong-dev/gobase/mg"
	"github.com/ChenLong-dev/gobase/mlog"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	busitask "myyy/src/availserverd/saasyd_usability_business_task_grpc"
	"myyy/src/scom"
	"sync"
	"sync/atomic"
	"time"
)

type Man struct {
	sync.RWMutex
	tasks        []*scom.Task //此表根据信息查询和ACL结果来生成
	eyeInfos     map[string]*scom.EyeInfo
	updateIntvl  time.Duration
	checkIntvl   time.Duration
	checkTimeout time.Duration

	flyCheckCount  int64
	MaxFlyCheckNum int64

	updateTimer *msys.MTimer
	checkTimer  *msys.MTimer
}

const (
	GET_BUSINESS_TASK_LIMIT = 1000
)

func (man *Man) String() string {
	return fmt.Sprintf("{tasksNum:%d, updateIntvl:%v, checkIntvl:%v, checkTimeout:%v, flyCheckCount:%d}", len(man.tasks), man.updateIntvl, man.checkIntvl, man.checkTimeout, atomic.LoadInt64(&man.flyCheckCount))
}
func (man *Man) FlyCheckCount() int64 {
	return atomic.LoadInt64(&man.flyCheckCount)
}
func (man *Man) UpdateTask() {
	// 测试
	//if err := man.loadEyeInfo(); err != nil {
	//	mlog.Errorf("loadEyeInfo fail, err: %v", err)
	//	return
	//}

	// 生产
	if err := man.getBusinessTaskList(); err != nil {
		mlog.Errorf("getBusinessTaskList fail, err: %v", err)
		return
	}
}

func (man *Man) loadEyeInfo() error {
	table := "EyeInfo"
	mlog.Debugf("====begin loadEyeInfo   {table:%s} .... ", table)
	collection := mg.CreateCollection(table)
	cursor, err := mg.Find(collection, bson.D{{}}, 0, 0)
	if err != nil {
		mlog.Warn(err)
		return err
	}
	defer cursor.Close(context.Background())

	if err = cursor.Err(); err != nil {
		mlog.Warn(err)
		return err
	}
	type EyeInfoDB struct {
		Url  string `bson:"url"`
		Freq int32  `bson:"freq"`
	}
	for cursor.Next(context.Background()) {
		var eyeInfoDB EyeInfoDB
		if err = cursor.Decode(&eyeInfoDB); err != nil {
			mlog.Warn(err)
			return err
		}

		var eyeInfo *scom.EyeInfo
		if ei, ok := man.eyeInfos[eyeInfoDB.Url]; ok {
			eyeInfo = ei
		} else {
			eyeInfo = &scom.EyeInfo{}
			man.eyeInfos[eyeInfoDB.Url] = eyeInfo
		}
		eyeInfo.Url = eyeInfoDB.Url
		eyeInfo.Freq = eyeInfoDB.Freq
	}

	mlog.Debugf("loadEyeInfo len(man.eyeInfos)=%d", len(man.eyeInfos))
	if len(man.eyeInfos) == 0 {
		return fmt.Errorf("data is nil!")
	}

	mlog.Debugf("====end loadEyeInfo")
	return nil
}
func (man *Man) SetTasks(tasks []*scom.Task) {
	man.Lock()
	man.tasks = tasks
	man.Unlock()
}

type CheckerEnv struct {
	checker *Checker
	Rsp     *scom.CheckRsp
	Err     error
}

func (man *Man) checkOne(checker *Checker, checkReq *scom.CheckReq, c chan *CheckerEnv) {
	checkRsp, err := checker.Check(checkReq, /*man.checkTimeout * 3*/ 0)
	if err != nil {
		c <- &CheckerEnv{checker: checker, Rsp: nil, Err: err}
	}
	c <- &CheckerEnv{checker: checker, Rsp: checkRsp}
}

func checkTasksCache(ei *scom.EyeInfo) bool {
	now := time.Now()
	if now.Sub(ei.Trigger) > time.Second*time.Duration(ei.Freq)*60 {
		return true
	}
	return false
}

func (man *Man) Check() {
	atomic.AddInt64(&man.flyCheckCount, 1)
	defer atomic.AddInt64(&man.flyCheckCount, -1)

	var tasks []*scom.Task

	man.RLock()
	for url, ei := range man.eyeInfos {
		if !checkTasksCache(ei) {
			continue
		}
		var task scom.Task
		var line scom.Line
		task.Method = "GET"
		task.Url = url
		line.Addr = url
		line.Isp = 0
		task.PrimaryAddrs = append(task.PrimaryAddrs, line)
		tasks = append(tasks, &task)
	}
	man.RUnlock()

	if len(tasks) == 0 {
		return
	}
	checkReq := &scom.CheckReq{}
	checkReq.Tasks = tasks
	checkReq.Timeout = int(man.checkTimeout / time.Second)
	checkReq.CheckTime = time.Now()

	cks := defaultCheckers.All()
	cksLen := len(cks)
	if cksLen == 0 {
		return
	}
	mlog.Debugf("xxx 0 [req:%v]\n", checkReq)
	c := make(chan *CheckerEnv, cksLen)
	for _, ck := range cks {
		go man.checkOne(ck, checkReq, c)
	}

	pr := &PointResult{}
	pr.CheckTime = time.Now()
	pr.Results = make(map[string]*UrlResult)

	beginTime := time.Now()
	for i := 0; i < cksLen; i++ {
		checkEnv := <-c
		region := checkEnv.checker.Region
		vt1 := time.Now()
		if checkEnv.Err != nil {
			mlog.Warnf("checker(%v) return error:%v", checkEnv.checker, checkEnv.Err)
			continue
		}

		if checkEnv.Rsp == nil {
			mlog.Debugf("Check rsp nil")
			return
		}
		for i, cr := range checkEnv.Rsp.Results {
			task := checkReq.Tasks[i]

			var ur *UrlResult
			if r, ok := pr.Results[task.Url]; ok {
				ur = r
			} else {
				ur = &UrlResult{Url: task.Url, Primarys: make(map[string]*scom.RegionResult), Seconds: make(map[string]*scom.RegionResult)}
				pr.Results[task.Url] = ur
			}

			mlog.Debugf("xxx 1 [url:%s] [region:%s] [res:%v]\n", task.Url, region, cr)
			makeLineResults := func(lines []scom.Line, rs []scom.Result) (lrs []scom.LineResult) {
				if len(lines) != len(rs) {
					mlog.Warnf("region:%s, len(lines):%d, lines value:[%+v], len(rs):%d, rs value:[%+v],", region, len(lines), lines, len(rs), rs)
					return
				}
				for i, r := range rs {
					lr := scom.LineResult{Line: lines[i], Result: r}
					lrs = append(lrs, lr)
				}
				return lrs
			}
			ur.Primarys[region] = &scom.RegionResult{Region: region, LineResults: makeLineResults(task.PrimaryAddrs, cr.Primarys)}
			ur.Seconds[region] = &scom.RegionResult{Region: region, LineResults: makeLineResults(task.SecondAddrs, cr.Seconds)}
		}
		vt2 := time.Now()
		mlog.Infof("1 =xxx= [region:%s] [exp1:%v] [exp2:%v]", region, vt2.Sub(vt1), vt2.Sub(checkReq.CheckTime))
	}

	vt3 := time.Now()
	defaultResultsMan.AddPoint(pr, cksLen)
	vt5 := time.Now()
	mlog.Infof("2 =xxx= [expend:%v], [total costtime:%v]", vt5.Sub(vt3), vt5.Sub(beginTime))
}

func (man *Man) getBusinessTaskList() error {
	mlog.Debugf("===getBusinessTaskList enter")
	//pick one service in random
	var busiTaskClient busitask.BusinessTaskListClient
	var conn *grpc.ClientConn
	if etcdNode, ok := Mng.GetServiceInfoRandom(BusiTaskServicePath); ok {
		mlog.Debugf("etcdNode is %+v", etcdNode)
		var err error
		conn, err = grpc.Dial(etcdNode.Info, grpc.WithInsecure())
		if err != nil {
			mlog.Error(err)
			return err
		}
		busiTaskClient = busitask.NewBusinessTaskListClient(conn)
	} else {
		mlog.Errorf("no service is chosen: %+v", etcdNode)
		return errors.New("No service is chosen ")
	}
	defer conn.Close()

	//批量分页拉取检测任务
	var offset int32
	for idx := 0; ; idx++ {
		var rsp *busitask.BusinessTaskRsp
		var getBusinessTaskErr error
		rsp, getBusinessTaskErr = busiTaskClient.GetBusinessTask(context.Background(),
			&busitask.BusinessTaskReq{
				Offset: offset,
				Limit:  GET_BUSINESS_TASK_LIMIT,
			})
		if getBusinessTaskErr != nil {
			mlog.Error(getBusinessTaskErr)
			return getBusinessTaskErr
		}
		if rsp == nil {
			mlog.Error(getBusinessTaskErr)
			return getBusinessTaskErr
		}
		if rsp.GetCode() != 0 || rsp.GetData() == nil {
			mlog.Error(rsp.GetMessage())
			return getBusinessTaskErr
		}

		len := len(rsp.GetData().GetTaskList())
		mlog.Debugf("idx:%d, offset:%d,getBusinessTaskList grpc rsp len(rsp):%d, rsp total:%d",
			idx, offset, len, rsp.GetData().GetTotal())
		man.Lock()
		for _, val := range rsp.GetData().GetTaskList() {
			//mlog.Debugf("1111 val:%v", val)
			var eyeInfo *scom.EyeInfo
			if ei, ok := man.eyeInfos[val.GetUrl()]; ok {
				eyeInfo = ei
			} else {
				eyeInfo = &scom.EyeInfo{}
				man.eyeInfos[val.GetUrl()] = eyeInfo
			}
			eyeInfo.Url = val.GetUrl()
			eyeInfo.Freq = val.GetFrequency()
		}
		man.Unlock()

		if len == GET_BUSINESS_TASK_LIMIT {
			offset += GET_BUSINESS_TASK_LIMIT
			mlog.Debugf("getBusinessTaskList len:%d, offset:%d", len, offset)
		} else {
			mlog.Debugf("getBusinessTaskList break")
			break
		}
	}

	mlog.Debugf("===getBusinessTaskList end")
	return nil
}

func (man *Man) Start() {
	mlog.Tracef("")

	man.UpdateTask()
	man.Check()

	man.updateTimer = msys.StartMTimer(int64(man.updateIntvl/time.Millisecond), updateTimer, man)
	man.checkTimer = msys.StartMTimer(int64(man.checkIntvl/time.Millisecond), checkTimer, man)
}

func updateTimer(pridata interface{}) {
	man := pridata.(*Man)
	man.UpdateTask()
	man.updateTimer = msys.StartMTimer(man.updateTimer.Elapsed(), updateTimer, man)
}
func checkTimer(pridata interface{}) {
	man := pridata.(*Man)
	if man.FlyCheckCount() < man.MaxFlyCheckNum {
		mlog.Debugf("checkTimer man.FlyCheckCount() < man.MaxFlyCheckNum")
		go man.Check()
	} else {
		mlog.Debugf("checkTimer else")
		man.Check()
	}

	man.checkTimer = msys.StartMTimer(man.checkTimer.Elapsed(), checkTimer, man)
}
func NewMan(updateTaskIntvl time.Duration, checkIntvl time.Duration, checkTimeout time.Duration, maxFlyCheckCount int64) *Man {
	man := &Man{updateIntvl: updateTaskIntvl, checkIntvl: checkIntvl, checkTimeout: checkTimeout, MaxFlyCheckNum: maxFlyCheckCount,
		eyeInfos: make(map[string]*scom.EyeInfo)}

	return man
}

var defaultMan *Man

func InitMan(updateTaskIntvl time.Duration, checkIntvl time.Duration, checkTimeout time.Duration, maxFlyCheckCount int64) {
	defaultMan = NewMan(updateTaskIntvl, checkIntvl, checkTimeout, maxFlyCheckCount)
	defaultMan.Start()
}
func initMongoDB(host, username, password, dbname string) error {
	if err := mg.Connect(host, username, password, dbname); err != nil {
		mlog.Errorf("init MongoDB is failed, err:%+v\n", err)
		return err
	}
	mlog.Info("init MongoDB is success ...")
	return nil
}

func (em Man) CheckEyeInfosCache(eyeInfos map[string]scom.EyeInfo) {
	em.Lock()
	defer em.Unlock()
	now := time.Now()
	for url, eyeInfo := range eyeInfos {
		if ei, ok := em.eyeInfos[url]; ok {
			if eyeInfo.Freq != ei.Freq {
				em.eyeInfos[url] = &scom.EyeInfo{Url: url, Freq: eyeInfo.Freq, Trigger: now}
			} else {
				continue
			}
		} else {
			em.eyeInfos[url] = &scom.EyeInfo{Url: url, Freq: eyeInfo.Freq, Trigger: now}
		}
	}
}

func (em *Man) AddEyeInfosCache(eyeInfo scom.EyeInfo) error {
	em.Lock()
	defer em.Unlock()
	now := time.Now()
	if ei, ok := em.eyeInfos[eyeInfo.Url]; ok {
		return fmt.Errorf("[url:%s] is exit [%v]", eyeInfo.Url, ei)
	} else {
		em.eyeInfos[eyeInfo.Url] = &scom.EyeInfo{Url: eyeInfo.Url, Freq: eyeInfo.Freq, Trigger: now}
	}
	mlog.Infof("add eye info is success! [eyeInfo:%v]\n", eyeInfo)
	return nil
}

func (em *Man) DelEyeInfosCache(url string) error {
	em.Lock()
	defer em.Unlock()
	if _, ok := em.eyeInfos[url]; ok {
		delete(em.eyeInfos, url)
	} else {
		return fmt.Errorf("not data in em! [url:%s]\n", url)
	}
	mlog.Infof("del eye info is success! [url:%s]\n", url)
	return nil
}

func (em *Man) EditEyeInfosCache(eyeInfo scom.EyeInfo) error {
	em.Lock()
	defer em.Unlock()
	now := time.Now()
	if ei, ok := em.eyeInfos[eyeInfo.Url]; ok {
		ei.Freq = eyeInfo.Freq
		ei.Trigger = now
	} else {
		//em.eyeInfos[eyeInfo.Url] = &scom.EyeInfo{Url: eyeInfo.Url, Freq: eyeInfo.Freq, Trigger: now}
		return fmt.Errorf("not data in em! [url:%s]\n", eyeInfo.Url)
	}
	mlog.Infof("edit eye info is success! [eyeInfo:%v]\n", eyeInfo)

	return nil
}