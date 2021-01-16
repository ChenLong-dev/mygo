package main


import (
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase/mcom"
	"github.com/ChenLong-dev/gobase/mbase/mutils"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyy/src/scom"
	"net/url"
	"strings"
)

func InitMdbg(addr string) (err error) {
	mlog.Infof("InitMdbg addr=%s", addr)

	if addr == "" {
		return nil
	}
	mcom.InitMdbg(addr)

	mcom.RegisterMdbgMml(mcom.MdbgPowerRoot, "lastPoint", "pretty=1", mdbgLastPoint)
	mcom.RegisterMdbgMml(mcom.MdbgPowerRoot, "allPoint", "detail=1,pretty=1", mdbgAllPoint)
	mcom.RegisterMdbgMml(mcom.MdbgPowerRoot, "result", "pointBound=%d,regionBound=%d", mdbgResult)
	mcom.RegisterMdbgMml(mcom.MdbgPowerRoot, "allChecker", "", mdbgAllChecker)
	mcom.RegisterMdbgMml(mcom.MdbgPowerRoot, "printMan", "", mdbgPrintMan)
	mcom.RegisterMdbgMml(mcom.MdbgPowerRoot, "setTasks", "method=HEAD/GET,urls=https://kdzl.cn|http://sangfor.com.cn,pretty=1", mdbgSetTasks)
	mcom.RegisterMdbgMml(mcom.MdbgPowerRoot, "setKeepPoingsNum", "keepPointsNum=%d", mdbgSetKeepPointsNum)

	return nil
}

func mdbgLastPoint(mml *mutils.Mml) (result int32, cmdRsp string) {
	pretty := mml.GetBool("pretty", true)
	pr := defaultResultsMan.GetLastPointResult()
	if pr == nil {
		return -1, "no data"
	} else {
		if pretty {
			return 0, mutils.JsonPrintPretty(pr)
		} else {
			return 0, mutils.JsonPrint(pr)
		}
	}
}
func mdbgAllPoint(mml *mutils.Mml) (result int32, cmdRsp string) {
	pretty := mml.GetBool("pretty", true)
	detail := mml.GetBool("detail", false)
	prs := defaultResultsMan.GetAllPointResult()

	if detail {
		if pretty {
			return 0, mutils.JsonPrintPretty(prs)
		} else {
			return 0, mutils.JsonPrint(prs)
		}
	} else {
		str := fmt.Sprintf("pointNum:%d\n", len(prs))
		for _, pr := range prs {
			str += fmt.Sprintf("    {CheckTime:%v, UrlNum:%d}\n", pr.CheckTime, len(pr.Results))
		}
		return 0, str
	}
}
func mdbgAllChecker(mml *mutils.Mml) (result int32, cmdRsp string) {
	checkers := defaultCheckers.All()

	str := fmt.Sprintf("checkers(%d):\n", len(checkers))
	for _, checker := range checkers {
		str += "    " + checker.String() + "\n"
	}

	return 0, str
}
func mdbgResult(mml *mutils.Mml) (result int32, cmdRsp string) {
	//pointBound := mml.GetInt("pointBound", 6)
	regionBound := mml.GetInt("regionBound", 1)
	pretty := mml.GetBool("pretty", true)

	rs := defaultResultsMan.GetResults(regionBound)

	if pretty {
		return 0, mutils.JsonPrintPretty(rs)
	} else {
		return 0, mutils.JsonPrint(rs)
	}
}
func mdbgPrintMan(mml *mutils.Mml) (result int32, cmdRsp string) {
	return 0, defaultMan.String()
}
func mdbgSetTasks(mml *mutils.Mml) (result int32, cmdRsp string) {
	urls := mml.GetString("urls", "")
	ss := strings.Split(urls, "|")
	method := mml.GetString("method", "GET")

	mlog.Tracef("method=%s,urls=%s,ss=%v", method, urls, ss)

	tasks := make([]*scom.Task, 0)
	for _, u := range ss {
		task := &scom.Task{Url: u, Method: method}
		if up, uperr := url.Parse(u); uperr == nil {
			line := scom.Line{Addr: up.Host}
			task.PrimaryAddrs = append(task.PrimaryAddrs, line)
			tasks = append(tasks, task)
		}
	}

	defaultMan.SetTasks(tasks)
	return 0, "ok"
}
func mdbgSetKeepPointsNum(mml *mutils.Mml) (result int32, cmdRsp string) {
	nw := mml.GetInt("keepPointsNum", 0)
	if nw <= 0 {
		return -1, "error: keepPointsNum illegal"
	} else {
		defaultResultsMan.SetKeepPointsNum(nw)
		return 0, "assign ok"
	}
}
