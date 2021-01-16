package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase/mutils"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyy/src/scom"
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

func (acr *AvailCheckRsp) String() string {
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

var runC chan *AvailCheck

func Init(maxC int) error {
	if maxC <= 0 {
		maxC = 1
	}

	runC = make(chan *AvailCheck, maxC)
	for i := 0; i < maxC; i++ {
		go runCheckC(runC)
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

func runCheckC(acs chan *AvailCheck) {
	for ac := range acs {
		ac.Rsp = checkOnce(ac.Req)
		ac.C <- struct{}{}
	}
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
	getAddr := func(addr string) string {
		ss := strings.Split(addr, ":")
		if len(ss) > 2 { //	ipv6
			mlog.Debugf("ResolveDns 1 [ipv6] [addr:%s]", addr)
			return addr
		}

		if !isIpv4(ss[0]) {
			if ip, err := mutils.ResolveDns(ss[0]); err == nil {
				if len(ss) > 1 {
					ipp := ip + ":" + ss[1]
					mlog.Debugf("ResolveDns 2 [url dns] [addr:%s] [ss:%s] [ip:%s] [port:%s] [ipp:%s]", addr, ss, ip, ss[1], ipp)
					return ipp
				}else {
					mlog.Debugf("ResolveDns 3 [url dns] [addr:%s] [ss:%s] [ip:%s]", addr, ss, ip)
					return ip
				}
			}
		}
		mlog.Debugf("ResolveDns 4 [addr:%s]", addr)
		return addr
	}

	ss1 := strings.Split(addr, "?")	// 去掉参数
	ss2 := strings.Split(ss1[0], "//") // 去掉http://或https://    www.xxx.com:80, www.xxx.com, 1.1.1.1, 1.1.1.1:80, fe80::fcfc:feff:feea:e23f:80, fe80::fcfc:feff:feea:e23f
	if len(ss2) < 2 {
		ss3 := strings.Split(ss2[0], "/")	// 去掉路径
		return getAddr(ss3[0])
	} else {
		ss3 := strings.Split(ss2[1], "/")	// 去掉路径
		return getAddr(ss3[0])
	}
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
	mlog.Debugf("xxx getIpAddr [addr.Addr:%s] [req.Addr:%v]", addr.Addr, req.Addr)

	return req
}

func getResult(ac *AvailCheck) scom.Result {
	<-ac.C
	result := scom.Result{}
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

func Check(tasks []*scom.Task, timeout time.Duration) (results []*scom.TaskResult, err error) {
	acs := make([]*AvailCheck, 0, len(tasks)*2)
	for _, task := range tasks {
		mlog.Debugf("www 1 [url:%s] [P:%v]", task.Url, task.PrimaryAddrs)
		if strings.Index(task.Url, "http://") < 0 && strings.Index(task.Url, "https://") < 0 {
			task.Url = "http://" + task.Url
		}
		if len(task.PrimaryAddrs) == 0 {
			u, ue := url.Parse(task.Url)
			if ue == nil {
				task.PrimaryAddrs = append(task.PrimaryAddrs, scom.Line{Addr: u.Host})
			}
		}
		mlog.Debugf("www 2 [url:%s] [P:%v]", task.Url, task.PrimaryAddrs)
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
		mlog.Debugf("result [url:%s] [t:%v] [tr:%v]", task.Url, task, taskResult)
	}

	return results, nil
}
