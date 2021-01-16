/**
* @Author: cl
* @Date: 2021/1/16 11:12
 */
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase/msys"
	"github.com/ChenLong-dev/gobase/mg"
	"github.com/ChenLong-dev/gobase/mlog"
	"google.golang.org/grpc"
	"myyd/src/availserverd/config"
	acl "myyd/src/availserverd/saasyd_acl_query_grpc"
	"myyd/src/availserverd/website"
	"myyd/src/scom"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Man struct {
	sync.RWMutex
	tasks            []*scom.Task      //此表根据信息查询和ACL结果来生成
	aclResults       map[string]string //源站是否配置ACL（"noacl"：否，"acl"：是）
	updateIntvl      time.Duration
	checkIntvl       time.Duration
	checkTimeout     time.Duration
	urlLastCheckTime map[string]int64
	allClusterIPMap  map[string][]scom.Line

	flyCheckCount  int64
	MaxFlyCheckNum int64

	updateTimer *msys.MTimer
	checkTimer  *msys.MTimer
}

var offset uint32 = 0

const (
	RETRY = 100 //重试次数
)

func (man *Man) String() string {
	return fmt.Sprintf("{tasksNum:%d, updateIntvl:%v, checkIntvl:%v, checkTimeout:%v, flyCheckCount:%d}", len(man.tasks), man.updateIntvl, man.checkIntvl, man.checkTimeout, atomic.LoadInt64(&man.flyCheckCount))
}
func (man *Man) FlyCheckCount() int64 {
	return atomic.LoadInt64(&man.flyCheckCount)
}
func (man *Man) UpdateTask() {
	//	todo 从信息查询模块和ACL模块获取信息来更新tasks
	mlog.Debugf("====begin updatetask")
	//加载ACL检测结果数据
	if DetectType == LINK_DETECT {
		for index := 0; index < RETRY; {
			len := len(scom.AclUrlRes)
			aclErr := man.loadAclResult()
			if aclErr != nil && len <= 0 {
				index++
				mlog.Errorf("loadAclResult fail, retry index:%d, err:%+v", index, aclErr)
				time.Sleep(time.Second * 10)
				continue
			}
			break
		}
		if err := man.loadDetectInfo(); err != nil {
			mlog.Errorf("loadDetectInfo fail, err: %v", err)
			return
		}
	}
	//加载检测客户信息数据
	if DetectType == ACL_DETECT {
		if err := man.loadDetectInfo(); err != nil {
			mlog.Errorf("loadDetectInfo fail, err: %v", err)
			return
		}
	}
	//加载线路节点任务
	if DetectType == LINE_DETECT {
		if err := man.loadLineTask(); err != nil {
			mlog.Errorf("loadDetectInfo fail, err: %v", err)
			return
		}
	}
	mlog.Debugf("====end updatetask")
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

func (man *Man) isCheck(urlCheckIntvl int64, urlLastCheckTime int64) bool {
	mlog.Debugf("xxx urlCheckIntvl:%d, urlLastCheckTime:%d, CheckIntvl:%d", urlCheckIntvl, urlLastCheckTime, CheckIntvl)
	if urlCheckIntvl != 0 && uint64(urlCheckIntvl) != CheckIntvl {
		if time.Now().Sub(time.Unix(urlLastCheckTime, 0)) < time.Second*time.Duration(urlCheckIntvl) {
			return false
		}
	}
	return true
}

func (man *Man) appendTargetTasks(tasks *[]*scom.Task, allClusterIPMap map[string][]scom.Line) int {
	num := 0
	for _, v := range config.Conf.AvailConf.TargetUrl {
		if !strings.Contains(v, `.detect.sre.ac.cn`) {
			continue
		}
		// http://bj.detect.sre.ac.cn		cluster_node_gz
		arr1 := strings.Split(v, "//")
		if len(arr1) < 2 {
			continue
		}
		arr2 := strings.Split(arr1[1], ".")
		nodeName := "cluster_node_" + arr2[0]
		if allClusterIP, ok := allClusterIPMap[nodeName]; !ok {
			mlog.Warnf("not found the nodename [%s] in allClusterIPMap, [v:%s]", nodeName, v)
			continue
		} else {
			var task scom.Task
			task.Method = "GET"
			task.Url = v
			task.PrimaryAddrs = append(task.PrimaryAddrs, allClusterIP...)
			*tasks = append(*tasks, &task)
			num = num + len(allClusterIP)
			mlog.Infof("add target ip [url:%s] [nn:%s] [num:%d-len:%d]", v, nodeName, num, len(allClusterIP))
		}
	}
	return num
}

func (man *Man) checkIsInTargetUrls(url string) bool {
	for _, v := range config.Conf.AvailConf.TargetUrl {
		if v == url {
			return true
		}
	}
	return false
}

func (man *Man) Check() {
	atomic.AddInt64(&man.flyCheckCount, 1)
	defer atomic.AddInt64(&man.flyCheckCount, -1)

	man.Lock()
	originTasks := man.tasks
	tLen1 := len(originTasks)
	var tasks = make([]*scom.Task, 0)
	allClusterIPMap := man.allClusterIPMap
	if DetectType == "link" {
		for i := 0; i < tLen1; i++ {
			mlog.Debugf("xxx :%+v", originTasks[i])
			if _, ok := man.urlLastCheckTime[originTasks[i].Url];
				originTasks[i].CheckLevel == 0 || (ok && !man.isCheck(originTasks[i].CheckIntvl, man.urlLastCheckTime[originTasks[i].Url])) {
				//tasks = append(tasks[:i], tasks[i+1:]...)
				continue
			} else {
				man.urlLastCheckTime[originTasks[i].Url] = time.Now().Unix()
				tasks = append(tasks, originTasks[i])
			}
		}
	}
	man.Unlock()

	tLen2 := len(tasks)

	num := 0
	if DetectType == "link" && len(allClusterIPMap) > 0 { // 添加靶机检测
		num = man.appendTargetTasks(&tasks, allClusterIPMap)
	}
	tLen3 := len(tasks)
	mlog.Debugf("[tasks info:%v]\n", tasks)
	mlog.Infof("[len1:%d] [len2:%d] [len3:%d] [target.task:%d] [target.num:%d]\n", tLen1, tLen2, tLen3,
		len(config.Conf.AvailConf.TargetUrl), num)
	if tLen3 == 0 {
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

	c := make(chan *CheckerEnv, cksLen)
	for _, ck := range cks {
		go man.checkOne(ck, checkReq, c)
		if DetectType == "acl" {
			time.Sleep(time.Second * 5)
		}
	}

	pr := &PointResult{}
	pr.CheckTime = time.Now()
	pr.Results = make(map[string]*UrlResult)

	tr := &PointResult{}
	tr.CheckTime = pr.CheckTime
	tr.Results = make(map[string]*UrlResult)

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
			if man.checkIsInTargetUrls(task.Url) { // 靶机结果
				if r, ok := tr.Results[task.Url]; ok {
					ur = r
				} else {
					ur = &UrlResult{Url: task.Url, Primarys: make(map[string]*scom.RegionResult), Seconds: make(map[string]*scom.RegionResult)}
					tr.Results[task.Url] = ur
				}
			} else { // 单业务结果
				if r, ok := pr.Results[task.Url]; ok {
					ur = r
				} else {
					ur = &UrlResult{Url: task.Url, Primarys: make(map[string]*scom.RegionResult), Seconds: make(map[string]*scom.RegionResult)}
					pr.Results[task.Url] = ur
				}
			}

			makeLineResults := func(lines []scom.Line, rs []scom.Result) (lrs []scom.LineResult) {
				if len(lines) != len(rs) {
					mlog.Warnf("lineres url:%s, region:%s, len(lines):%d, lines value:[%+v], len(rs):%d, rs value:[%+v],", task.Url, region, len(lines), lines, len(rs), rs)
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
			mlog.Tracef("3 =xxxxxxx= result info: [url: %v] [region: %v] [primary result:%v],", task.Url, region, cr.Primarys)
		}
		vt2 := time.Now()
		mlog.Infof("1 =xxx= [region:%s] [exp1:%v] [exp2:%v]", region, vt2.Sub(vt1), vt2.Sub(checkReq.CheckTime))
	}

	mlog.Debugf("4 =xxxxxx= [pr.result.len: %d] [tr.result.len: %d] ", len(pr.Results), len(tr.Results))

	vt3 := time.Now()
	if DetectType == ACL_DETECT {
		defaultResultsMan.AddAclPoint(pr)
	} else if DetectType == LINK_DETECT {
		defaultResultsMan.AddLinkPoint(pr)
		targetResultsMan.AddTargetPoint(tr, cksLen, allClusterIPMap)
	} else if DetectType == LINE_DETECT {
		defaultResultsMan.AddLinePoint(pr, cksLen)
	}

	vt5 := time.Now()
	mlog.Infof("2 =xxx= [exp:%v], [total exp:%v]", vt5.Sub(vt3), vt5.Sub(beginTime))
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
		aclResults: make(map[string]string), urlLastCheckTime: make(map[string]int64)}
	//man.tasks = getZiboTasks()

	scom.AclUrlRes = make(map[string]string)
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
func (man *Man) loadDetectInfo() error {
	conn, err := grpc.Dial(GrpcAddr, grpc.WithInsecure())
	if err != nil {
		mlog.Error(err)
		return err
	}
	defer conn.Close()
	client := website.NewWebsiteClient(conn)
	//批量拉取需要链路检测的客户信息
	var rsp *website.GetBusinessInfosRsp
	var getBusinessInfoErr error
	rsp, getBusinessInfoErr = client.GetBusinessInfos(context.Background(),
		&website.GetBusinessInfosReq{
			RefClusterNode: "",
			RefClouduser:   "",
			Offset:         offset,
			Limit:          uint32(DetectLimit),
			MsgType:        "xxxxx",
		})
	if getBusinessInfoErr != nil {
		mlog.Error(getBusinessInfoErr)
		return getBusinessInfoErr
	}
	if rsp == nil {
		mlog.Error(getBusinessInfoErr)
		return getBusinessInfoErr
	}
	if rsp.GetCode() != 0 || rsp.GetData() == nil {
		mlog.Error(rsp.GetMsg())
		return getBusinessInfoErr
	}

	len := uint32(len(rsp.GetData().GetList()))
	if len <= 0 {
		mlog.Debugf("get list len : %d", len)
		offset = 0
		return nil
	} else {
		offset = offset + uint32(DetectLimit)
	}
	mlog.Infof("get acl business list [len:%d], [limit:%d], [offset:%d]", len, gConfig.DetectLimit, offset)
	if DetectType == ACL_DETECT {
		if err := man.loadAclTasks(rsp.GetData().GetList()); err != nil {
			mlog.Errorf("loadAclTasks fail, err: %v", err)
			return err
		}
	} else if DetectType == LINK_DETECT {
		if err := man.loadLinkTasks(rsp.GetData().GetList()); err != nil {
			mlog.Errorf("loadAclTasks fail, err: %v", err)
			return err
		}
	} else {
		msg := fmt.Sprintf("detect type err, DetectType:%s", DetectType)
		mlog.Error(msg)
		return errors.New(msg)
	}

	return nil
}

func (man *Man) loadAclResult() error {
	mlog.Debugf("enter loadAclResult grpc addr:%s", AclGrpcAddr)
	conn, err := grpc.Dial(AclGrpcAddr, grpc.WithInsecure())
	if err != nil {
		mlog.Error(err)
		return err
	}
	defer conn.Close()
	client := acl.NewAclGrpcServerClient(conn)
	//批量拉取需要链路检测的ACL信息
	var rsp *acl.AclResultRsp
	var getAclErr error
	rsp, getAclErr = client.GetAclResults(context.Background(),
		&acl.AclResultReq{
			Offset: offset,
			Limit:  uint32(DetectLimit),
			Type:   "skip_unknown",
		})
	if getAclErr != nil {
		mlog.Error(getAclErr)
		return getAclErr
	}
	if rsp == nil {
		mlog.Error(getAclErr)
		return getAclErr
	}
	if rsp.GetCode() != 0 {
		mlog.Error(rsp.GetMsg())
		return getAclErr
	}

	mlog.Debugf("loadAclResult grpc rsp:%+v", rsp)
	len := uint32(len(rsp.GetData()))
	if len <= 0 {
		mlog.Debugf("get list len : %d", len)
		offset = 0
		return nil
	} else {
		offset = offset + uint32(DetectLimit)
	}

	mlog.Tracef("list [len:%d], [limit:%d], [offset:%d]", len, gConfig.DetectLimit, offset)
	for _, data := range rsp.GetData() {
		if data.GetResult() == "acl" {
			scom.AclUrlRes[data.GetBusinessUrl()] = "acl"
		} else if data.GetResult() == "noacl" {
			scom.AclUrlRes[data.GetBusinessUrl()] = "noacl"
		}
	}

	mlog.Debugf("loadAclResult scom.AclUrlRes:%+v", scom.AclUrlRes)

	return nil
}

func (man *Man) loadAllClusterIP() error {
	mlog.Debugf("===loadAllClusterIP enter")
	conn, err := grpc.Dial(GrpcAddr, grpc.WithInsecure())
	if err != nil {
		mlog.Error(err)
		return err
	}
	defer conn.Close()
	client := website.NewWebsiteClient(conn)
	//批量拉取所有检测IP
	var rsp *website.GetClusterIPInfosRsp
	var getIpInfoErr error
	rsp, getIpInfoErr = client.GetClusterIpInfos(context.Background(),
		&website.GetClusterIPInfosReq{
			Offset:  0,
			Limit:   0,
			MsgType: "target",
		})
	if getIpInfoErr != nil {
		mlog.Error(getIpInfoErr)
		return getIpInfoErr
	}
	if rsp == nil {
		mlog.Error(getIpInfoErr)
		return getIpInfoErr
	}
	if rsp.GetCode() != 0 || rsp.GetData() == nil {
		mlog.Error(rsp.GetMsg())
		return getIpInfoErr
	}

	mlog.Debugf("loadAllClusterIP grpc rsp len(rsp.GetData().GetList()):%v", len(rsp.GetData().GetList()))

	aci := make(map[string][]scom.Line)
	for _, value := range rsp.GetData().GetList() {
		//if value.GetType() == 4 { // 排除独享IP，type=4
		//	continue
		//}
		nodeName := value.GetNodeName()
		var line scom.Line
		line.Addr = value.GetIp().GetExtranetIp()
		if line.Addr != "" {
			if value.GetIp().GetType() == scom.TELECOM {
				line.Isp = 1
			} else if value.GetIp().GetType() == scom.UNICOM {
				line.Isp = 2
			} else if value.GetIp().GetType() == scom.MOBILE {
				line.Isp = 3
			} else {
				line.Isp = 0
			}
			if ll, ok := aci[nodeName]; ok {
				ll = append(ll, line)
				aci[nodeName] = ll
			} else {
				var al []scom.Line
				al = append(al, line)
				aci[nodeName] = al
			}
		}
	}

	count := 1
	for nn, l := range aci {
		mlog.Infof("[%d] cluster [nn:%s] ip count [l.len:%d]\n", count, nn, len(l))
		count++
	}

	man.Lock()
	man.allClusterIPMap = aci
	man.Unlock()

	//mlog.Debugf("[allClusterIP.len: %d] [allClusterIP: %v]", len(allClusterIP), allClusterIP)
	mlog.Debugf("===loadAllClusterIP end")
	return nil
}

func (man *Man) loadLinkTasks(list []*website.GetBusinessInfosRsp_BusinessInfo) error {
	mlog.Debugf("====begin loadLinkTasks")
	//检测机房节点信息存到allClusterIPMap
	if err := man.loadAllClusterIP(); err != nil {
		mlog.Errorf("loadAllClusterIPMap fail, err: %v", err)
	}

	//检测客户信息存储到Task结构数组
	mlog.Debugf("list len: %d", len(list))
	var tasks []*scom.Task
	for _, data := range list {
		var task scom.Task
		task.Method = "GET"
		task.Url = data.GetUrl()
		task.CheckIntvl = int64(data.GetCheckIntvl())
		task.CheckLevel = int(data.GetCheckLevel())
		for _, value := range data.GetInsiteIp() {
			var line scom.Line
			line.Addr = value.GetExtranetIp()
			if line.Addr != "" {
				if value.GetType() == scom.TELECOM {
					line.Isp = 1
				} else if value.GetType() == scom.UNICOM {
					line.Isp = 2
				} else if value.GetType() == scom.MOBILE {
					line.Isp = 3
				} else {
					line.Isp = 0
				}
				task.PrimaryAddrs = append(task.PrimaryAddrs, line)
			}
		}
		mlog.Debugf("ttttt [task:%v]", task)

		if aclVal, ok := scom.AclUrlRes[data.GetUrl()]; ok {
			if aclVal == "acl" {
				var line scom.Line
				reserveIp := data.GetReserveInsiteIp()
				line.Addr = reserveIp.GetExtranetIp()
				if line.Addr != "" {
					if reserveIp.GetType() == scom.TELECOM {
						line.Isp = 1
					} else if reserveIp.GetType() == scom.UNICOM {
						line.Isp = 2
					} else if reserveIp.GetType() == scom.MOBILE {
						line.Isp = 3
					} else {
						line.Isp = 0
					}

					//if reserveIp.ExtranetIp == "" {
					//	//告警通知
					//	content := fmt.Sprintf("[%s], url:[%s], reason:[%s], time:[%+v]",
					//		DetectType, data.Url, "未配置备节点IP", time.Now().Format("2006-01-02 15:04:05"))
					//	alarm.SendAlarmMsg(data.Url, content)
					//}
					task.SecondAddrs = append(task.SecondAddrs, line)
				}
			} else if aclVal == "noacl" {
				for _, value := range data.GetIp() {
					var line scom.Line
					line.Addr = value
					if line.Addr != "" {
						line.Isp = 0
						task.SecondAddrs = append(task.SecondAddrs, line)
					}
				}
			} else {
				continue
			}
			if len(task.PrimaryAddrs) == 0 {
				continue
			}
			tasks = append(tasks, &task)
		}
	}

	man.Lock()
	man.tasks = tasks
	man.Unlock()

	mlog.Debugf("loadLinkTasks man.tasks: %v", man.tasks)
	mlog.Debugf("====end loadLinkTasks")

	return nil
}
func (man *Man) loadAclTasks(list []*website.GetBusinessInfosRsp_BusinessInfo) error {
	mlog.Debugf("enter loadAclTasks")
	//检测客户信息存储到Task结构数组
	var tasks []*scom.Task
	for _, data := range list {
		var task scom.Task
		task.Method = "GET"
		task.Url = data.GetUrl()
		task.CheckIntvl = int64(data.GetCheckIntvl())
		task.CheckLevel = int(data.GetCheckLevel())
		for _, value := range data.GetIp() {
			var line scom.Line
			line.Addr = value
			line.Isp = 0
			task.PrimaryAddrs = append(task.PrimaryAddrs, line)
		}
		tasks = append(tasks, &task)
	}
	man.Lock()
	man.tasks = tasks
	man.Unlock()

	return nil
}

func (man *Man) loadLineTask() error {
	mlog.Debugf("===loadLineTask enter")
	conn, err := grpc.Dial(GrpcAddr, grpc.WithInsecure())
	if err != nil {
		mlog.Error(err)
		return err
	}
	defer conn.Close()
	client := website.NewWebsiteClient(conn)
	//批量拉取线路检测IP
	var rsp *website.GetClusterIPInfosRsp
	var getIpInfoErr error
	rsp, getIpInfoErr = client.GetClusterIpInfos(context.Background(),
		&website.GetClusterIPInfosReq{
			Offset:  0,
			Limit:   0,
			MsgType: "line",
		})
	if getIpInfoErr != nil {
		mlog.Error(getIpInfoErr)
		return getIpInfoErr
	}
	if rsp == nil {
		mlog.Error(getIpInfoErr)
		return getIpInfoErr
	}
	if rsp.GetCode() != 0 || rsp.GetData() == nil {
		mlog.Error(rsp.GetMsg())
		return getIpInfoErr
	}

	mlog.Infof("loadLineTask grpc rsp len(rsp.GetData().GetList()):%v", len(rsp.GetData().GetList()))
	var tasks []*scom.Task
	for _, value := range rsp.GetData().GetList() {
		var task scom.Task
		var line scom.Line
		task.Url = value.GetIp().GetExtranetIp()
		line.Addr = value.GetIp().GetExtranetIp()
		task.PrimaryAddrs = append(task.PrimaryAddrs, line)
		tasks = append(tasks, &task)
	}

	man.Lock()
	man.tasks = tasks
	man.Unlock()

	mlog.Debugf("===loadLineTask end")
	//mlog.Debugf("loadLineTask man.tasks: %v", man.tasks)
	return nil
}
