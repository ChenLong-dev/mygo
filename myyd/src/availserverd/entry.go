/**
* @Author: cl
* @Date: 2021/1/16 11:08
 */
package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase/mcom"
	"github.com/ChenLong-dev/gobase/mbase/msys"
	"github.com/ChenLong-dev/gobase/mbase/mutils"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyd/src/scom"
	"time"
)

type accessChecker struct {
	msrv 	*mcom.MServer
	timer   *msys.MTimer
}

func entryNotify(mctx *mcom.MContext, msrv *mcom.MServer) {
	mctx.Tracef("msrv=%v", msrv)

	achecker := &accessChecker{msrv: msrv}
	msrv.Pridata = achecker

	tf := func(pridata interface{}) {
		if achecker.msrv != nil {
			achecker.msrv.Close(mcom.NewMContext())
		}
	}

	achecker.timer = msys.StartMTimer(10*1000, tf, achecker)

	msrv.EnableGoCallback(true)
}
func entryError(mctx *mcom.MContext, msrv *mcom.MServer, err error) {
	mctx.Tracef("msrv=%v,err=%v", msrv, err)

	if achecker, ok := msrv.Pridata.(*accessChecker); ok {
		achecker.timer.Stop()
		msrv.Pridata = nil
		achecker.msrv = nil
		msrv.Close(mctx)
	} else if checker, ok := msrv.Pridata.(*Checker); ok {
		defaultCheckers.Del(checker.Region)
		checker.Close()
	}
}
const AvailTimeDifference uint64 = 300
func handleLogin(mctx *mcom.MContext, msrv *mcom.MServer, data []byte) {
	mctx.Tracef("msrv=%v", msrv)

	loginReq := &scom.LoginReq{}
	err := json.Unmarshal(data, loginReq)
	if err != nil {
		mlog.Warnf("json.Unmarshal error:%v", err)
		msrv.Reply(mctx, err, nil)
		return
	}

	timestamp := uint64(time.Now().Second())
	sign := scom.Signature(loginReq.Region, loginReq.Rand, loginReq.Timestamp, gConfig.EntryPassword)
	if sign != loginReq.Signature {
		err = fmt.Errorf("region or password error")
	} else if timestamp > loginReq.Timestamp + AvailTimeDifference || loginReq.Timestamp > timestamp + AvailTimeDifference {
		err = fmt.Errorf("timestamp error")
	}
	if err != nil {
		mlog.Warnf("region(%s) login error:%v, login info(%v)", loginReq.Region, err, mutils.JsonPrint(loginReq))
		msrv.Reply(mctx, err, nil)
		return
	}

	if achecker, ok := msrv.Pridata.(*accessChecker); ok {
		achecker.timer.Stop()
		achecker.msrv = nil
		msrv.Pridata = nil
	}

	checker := NewChecker(loginReq.Region, msrv)
	defaultCheckers.Add(checker)
	msrv.Pridata = checker

	loginRsp := &scom.LoginRsp{Code: 0, Status: "ok"}
	loginData, _ := json.Marshal(loginRsp)
	msrv.Reply(mctx, nil, loginData)
}
func entryHandle(mctx *mcom.MContext, msrv *mcom.MServer, data []byte) {
	mctx.Tracef("msrv=%v, len(data)=%d", msrv, len(data))

	switch mctx.OpCode() {
	case scom.OPCODE_LOGIN_REQ:
		handleLogin(mctx, msrv, data)
	default:
		mlog.Warnf("unsupport opcode(%d), mctx=%v,msrv=%v,len(data)=%d", mctx.OpCode(), mctx, msrv, len(data))
	}
}


func InitEntry(addr string) (err error) {
	mlog.Tracef("addr=%s", addr)
	defer func() {mlog.Tracef("err=%v", err)} ()

	_, merr := mcom.MServerListen(mcom.NewMContext(), "tcp", addr, entryNotify, entryHandle, entryError/*, mcom.MServerOptDisableGo()*/)
	return merr
}