/**
* @Author: cl
* @Date: 2021/1/14 19:40
 */
package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase/mutils"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyd/src/availd/ping"
	"myyd/src/scom"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AvailCheckReq struct {
	Req     *http.Request `json:"Req,omitempty"`
	Addr    string        `json:"Addr,omitempty"`
	IsSsl   bool          `json:"IsSsl,omitempty"`
	Timeout time.Duration `json:"Timeout,omitempty"`
}

func (ac *AvailCheckReq) String() string {
	//return mutils.JsonPrint(ac)
	if ac == nil {
		return "nil"
	}
	return fmt.Sprintf("{Req:%v,Addr:%s,IsSsl:%t,Timeout:%v}", ac.Req, ac.Addr, ac.IsSsl, ac.Timeout)
}
func NewAvailCheckReq() *AvailCheckReq {
	return &AvailCheckReq{}
}

type AvailCheckRsp struct {
	Rsp   *http.Response `json:"Rsp,omitempty"`
	Err   error          `json:"Err,omitempty"`
	Delay time.Duration  `json:"Del,omitempty"`
}

type PingCheckReq struct {
	Addr    string        `json:"Addr,omitempty"`
	Timeout time.Duration `json:"Timeout,omitempty"`
}

type PingCheckRsp struct {
	Result string `json:"Result,omitempty"`
	Err    error  `json:"Err,omitempty"`
}

func (acr *AvailCheckRsp) String() string {
	//return mutils.JsonPrint(acr)
	if acr == nil {
		return fmt.Sprintf("nil")
	}
	return fmt.Sprintf("{Rsp:%v,Err:%v, Del:%v}", acr.Rsp, acr.Err, acr.Delay)
}
func NewAvailCheckRsp(rsp *http.Response) *AvailCheckRsp {
	return &AvailCheckRsp{Rsp: rsp}
}
func NewAvailCheckRspWithError(err error) *AvailCheckRsp {
	return &AvailCheckRsp{Err: err}
}

type AvailCheck struct {
	Req *AvailCheckReq
	Rsp *AvailCheckRsp
	C   chan struct{}
}

type PingCheck struct {
	Req *PingCheckReq
	Rsp *PingCheckRsp
	C   chan struct{}
	Ips map[string]bool
}

var runC chan *AvailCheck
var runP chan *PingCheck

func Init(maxC int) error {
	if maxC <= 0 {
		maxC = 1
	}

	//	httpClient = &http.Client{}
	//	httpClient.Timeout = 10*time.Second

	/*runC = make([]chan *AvailCheckReq, maxC)
	for i := 0; i < len(runC); i++ {
		runC[i] =
	}*/

	if DetectType == LINE_DETECT {
		runP = make(chan *PingCheck, maxC)
		for i := 0; i < maxC; i++ {
			go runPing(runP)
		}
	} else {
		runC = make(chan *AvailCheck, maxC)
		for i := 0; i < maxC; i++ {
			go runCheckC(runC)
		}
	}

	return nil
}

func checkOnce(req *AvailCheckReq) (rsp *AvailCheckRsp) {
	defer func() { mlog.Tracef("req=%v,rsp=%v", req, rsp) }()

	if req.Req == nil {
		return NewAvailCheckRspWithError(fmt.Errorf("make http.Request failed"))
	}

	if strings.Index(req.Addr, ":") < 0 {
		ss := strings.Split(req.Req.Host, ":")
		if len(ss) > 1 {
			req.Addr = req.Addr + ":" + ss[len(ss)-1]
		} else if req.IsSsl {
			req.Addr += ":443"
		} else {
			req.Addr += ":80"
		}
	}

	mlog.Tracef("req=%v", req)

	deadline := time.Now().Add(req.Timeout)
	conn, err := net.DialTimeout("tcp", req.Addr, /*10*time.Second*/ req.Timeout)
	if err != nil {
		return NewAvailCheckRspWithError(err)
	}
	conn.(*net.TCPConn).SetLinger(0)

	if req.Timeout != 0 {
		conn.SetDeadline( /*time.Now().Add(req.Timeout)*/ deadline)
	}

	if req.IsSsl {
		conn = tls.Client(conn, &tls.Config{InsecureSkipVerify: true, ServerName: req.Req.Host})
	}

	defer conn.Close()
	vt1 := time.Now()
	rsp = NewAvailCheckRsp(nil)

	err = req.Req.Write(conn)
	if err == nil {
		rsp.Rsp, err = http.ReadResponse(bufio.NewReader(conn), req.Req)
		vt2 := time.Now()
		rsp.Delay = vt2.Sub(vt1)
	}
	if err != nil {
		return NewAvailCheckRspWithError(err)
	}
	return rsp
}
func checkLine(c *PingCheck) (rsp *PingCheckRsp) {
	defer func() { mlog.Tracef("req=%v,rsp=%v", c.Req, rsp) }()
	rsp = &PingCheckRsp{}
	if c.Req.Addr == "" {
		mlog.Debugf("checkLine req.Addr is empty")
		rsp.Result = "0"
		rsp.Err = fmt.Errorf("checkLine req.Addr is empty")
		return rsp
	}
	if Ping(c.Req.Addr, c.Ips) {
		rsp.Result = "1"
		rsp.Err = nil
		return rsp
	} else {
		rsp.Result = "-1"
		rsp.Err = fmt.Errorf("checkLine Ping fail")
		return rsp
	}
}
func runCheckC(acs chan *AvailCheck) {
	for ac := range acs {
		ac.Rsp = checkOnce(ac.Req)
		ac.C <- struct{}{}
	}
}
func runPing(acs chan *PingCheck) {
	for ac := range acs {
		//ac.Rsp = checkLine(ac.Req)
		ac.Rsp = checkLine(ac)
		ac.C <- struct{}{}
	}
}

var gIdAllocator int64

func Ping(target string, ips map[string]bool) bool {
	pinger, err := ping.NewIpsPinger(target, ips)
	if err != nil {
		mlog.Debugf("exec ping fail, err:%v", err)
		return false
	}

	pinger.Count = 1
	pinger.Timeout = time.Duration(time.Duration(3) * time.Second)
	pinger.SetPrivileged(true)
	pinger.Run() // blocks until finished
	stats := pinger.Statistics()

	// 有回包，就是说明IP是可用的
	if stats.PacketsRecv >= 1 {
		return true
	}
	return false
}
func isIpv4(host string) bool {
	for _, c := range host {
		if !(c == '.' || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}
func getIpAddr(addr string) string {
	ss := strings.Split(addr, ":")
	if len(ss) > 2 { //	ipv6
		return addr
	} else if !isIpv4(ss[0]) {
		if ip, err := mutils.ResolveDns(ss[0]); err == nil {
			if len(ss) > 1 {
				return ip + ":" + ss[1]
			} else {
				return ip
			}
		}
	}
	return addr
}
func NewAvailCheckReqByTask(task *scom.Task, addr *scom.Line, timeout time.Duration) (req *AvailCheckReq) {
	mlog.Tracef("task=%v", task)
	defer func() { mlog.Tracef("req=%v", req) }()

	req = &AvailCheckReq{}

	url := task.Url
	if strings.HasPrefix(url, "https://") {
		req.IsSsl = true
	} /*else if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}*/

	httpReq, err := http.NewRequest(task.Method, url, nil)
	if err != nil {
		mlog.Warnf("task:%v http.NewRequest(%s, %s, nil) error:%v", task, task.Method, url, err)
		//return nil
	} else {
		for _, head := range task.Heads {
			httpReq.Header.Add(head.Key, head.Val)
		}
		if httpReq.Header.Get("User-Agent") == "" {
			httpReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")
		}
	}

	req.Req = httpReq

	req.Timeout = time.Duration(timeout)
	req.Addr = getIpAddr(addr.Addr)

	return req
}
func getResult(ac *AvailCheck) scom.Result {
	<-ac.C
	result := scom.Result{}
	mlog.Debugf("rsp:[%v]\n", ac.Rsp)
	if ac.Rsp.Err != nil {
		result.Code = -1
		result.Status = ac.Rsp.Err.Error()
		result.Delay = ac.Rsp.Delay
	} else {
		result.Code = ac.Rsp.Rsp.StatusCode
		result.Status = ac.Rsp.Rsp.Status
		result.Delay = ac.Rsp.Delay
	}
	return result
}

func getPingResult(ac *PingCheck) scom.Result {
	<-ac.C
	result := scom.Result{}
	mlog.Debugf("rsp:[%v]\n", ac.Rsp)
	if ac.Rsp.Result == "1" {
		result.Code = 1
	} else if ac.Rsp.Result == "-1" {
		result.Code = -1
	} else {
		result.Code = 0
	}
	return result
}
func Check(tasks []*scom.Task, timeout time.Duration) (results []*scom.TaskResult, err error) {
	acs := make([]*AvailCheck, 0, len(tasks)*2)
	for _, task := range tasks {
		if strings.Index(task.Url, "http://") < 0 && strings.Index(task.Url, "https://") < 0 {
			task.Url = "http://" + task.Url
		}
		if len(task.PrimaryAddrs) == 0 {
			u, ue := url.Parse(task.Url)
			if ue == nil {
				task.PrimaryAddrs = append(task.PrimaryAddrs, scom.Line{Addr: u.Host})
			}
		}
		for _, pa := range task.PrimaryAddrs {
			ac := &AvailCheck{}
			ac.Req = NewAvailCheckReqByTask(task, &pa, timeout)
			ac.C = make(chan struct{}, 1)
			runC <- ac

			acs = append(acs, ac)
		}
		for _, sa := range task.SecondAddrs {
			ac := &AvailCheck{}
			ac.Req = NewAvailCheckReqByTask(task, &sa, timeout)
			ac.C = make(chan struct{}, 1)
			runC <- ac

			acs = append(acs, ac)
		}
	}

	i := 0
	for _, task := range tasks {
		taskResult := &scom.TaskResult{}
		for range task.PrimaryAddrs {
			result := getResult(acs[i])
			taskResult.Primarys = append(taskResult.Primarys, result)
			i++
		}
		for range task.SecondAddrs {
			result := getResult(acs[i])
			taskResult.Seconds = append(taskResult.Seconds, result)
			i++
		}
		results = append(results, taskResult)
	}

	return results, nil
}
//func CheckPing(tasks []*scom.Task, timeout time.Duration) (results []*scom.TaskResult, err error) {
//	mlog.Debugf("enter CheckPing")
//	reqIps := make(map[string]bool)
//	for _, task := range tasks {
//		reqIps[task.Url] = false
//	}
//
//	acs := make([]*PingCheck, 0, len(tasks))
//	for _, task := range tasks {
//		ac := &PingCheck{}
//		ac.Req = &PingCheckReq{}
//		ac.Req.Addr = task.Url
//		ac.C = make(chan struct{}, 1)
//		ac.Ips = reqIps
//		runP <- ac
//		acs = append(acs, ac)
//	}
//
//	for i := 0; i < len(tasks); i++ {
//		taskResult := &scom.TaskResult{}
//		result := getPingResult(acs[i])
//		taskResult.Primarys = append(taskResult.Primarys, result)
//		mlog.Debugf("[ip:%s] [res:%v]\n", acs[i].Req.Addr, result)
//		results = append(results, taskResult)
//	}
//
//	return results, nil
//}

func CheckPing2(tasks []*scom.Task, timeout time.Duration) (results []*scom.TaskResult, err error) {
	mlog.Debugf("enter CheckPing")
	reqIps := make(map[string]bool)
	for _, task := range tasks {
		reqIps[task.Url] = false
	}

	acs := make([]*PingCheck, 0, len(tasks))
	for _, task := range tasks {
		ac := &PingCheck{}
		ac.Req = &PingCheckReq{}
		ac.Req.Addr = task.Url
		ac.C = make(chan struct{}, 1)
		ac.Ips = reqIps
		runP <- ac
		acs = append(acs, ac)
	}

	for i := 0; i < len(tasks); i++ {
		taskResult := &scom.TaskResult{}
		result := getPingResult(acs[i])
		taskResult.Primarys = append(taskResult.Primarys, result)
		mlog.Debugf("[ip:%s] [res:%v]\n", acs[i].Req.Addr, result)
		results = append(results, taskResult)
	}

	return results, nil
}

func CheckPing(tasks []*scom.Task, timeout time.Duration) ([]*scom.TaskResult, error) {
	var results []*scom.TaskResult
	var ts []*scom.Task
	for i, task := range tasks {
		ts= append(ts, task)
		if i % PingInterval == 0 {
			rs, err := CheckPing2(ts, timeout)
			if err == nil {
				results = append(results, rs...)
			}
			ts = ts[:0:0]
			mlog.Infof("check ping [i:%d] [pi:%d] [ts.len:%d] [rs:%d] [results.len:%d]", i, PingInterval, len(ts), len(rs), len(results))
			time.Sleep(time.Second*3)
		}
	}
	return results, nil
}

