package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase/mcom"
	"github.com/ChenLong-dev/gobase/mbase/mutils"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyy/src/scom"
	"time"
)

var gMserver *mcom.MServer
var gRegion string
var gPassword string

func login(region, password string) (err error) {
	mlog.Tracef("")
	defer func() {
		if err != nil {
			mlog.Warnf("login error:%v", err)
		} else {
			mlog.Tracef("login error:%v", err)
		}
	} ()

	loginReq := &scom.LoginReq{Region: region}
	loginReq.Region = region
	loginReq.Rand = mutils.RandomString(12, true, true, true)
	loginReq.Timestamp = uint64(time.Now().Second())
	//loginReq.Signature = mutils.SHA1Hex(loginReq.Region + loginReq.Rand + fmt.Sprint(loginReq.Timestamp) + password)
	loginReq.Signature = scom.Signature(loginReq.Region, loginReq.Rand, loginReq.Timestamp, password)

	data, _ := json.Marshal(loginReq)

	rsp, rerr := gMserver.Rpc(mcom.NewMContext(), mcom.MakeMHead(scom.OPCODE_LOGIN_REQ), data, 10*1000)
	if rerr != nil {
		return rerr
	}

	loginRsp := &scom.LoginRsp{}
	err = json.Unmarshal(rsp, loginRsp)
	if err != nil {
		return err
	}
	if loginRsp.Code < 0 {
		return fmt.Errorf("login failed, rsp=%s", mutils.JsonPrint(loginRsp))
	}

	return nil
}

func ConnectServer(addr string, region, password string) (err error) {
	mlog.Tracef("addr=%s", addr)
	defer func() {mlog.Tracef("err=%v", err)} ()

	gRegion = region
	gPassword = password
	ctx := mcom.NewMContext()

	gMserver, err = mcom.MServerDial(ctx, "tcp", addr, reConnect, handle, errorNotify)
	if err != nil {
		return err
	}

	return login(region, password)
}

func reConnect(mctx *mcom.MContext, msrv *mcom.MServer) {
	mlog.Infof("reconnect server ok!")

	login(gRegion, gPassword)
}

func handleCheck(mctx *mcom.MContext, msrv *mcom.MServer, data []byte) (rdata []byte, rerr error) {
	mctx.Tracef("len(data)=%d", len(data))
	vt1 := time.Now()
	defer func() {
		if rerr != nil {
			mlog.Warnf("len(rdata)=%d,rerr=%v", len(rdata), rerr)
		} else {
			mctx.Tracef("len(rdata)=%d,rerr=%v", len(rdata), rerr)
		}
		vt2 := time.Now()
		mlog.Infof("x=x=x=x [exp:%v]\n", vt2.Sub(vt1))
	} ()

	checkData, gzerr := mutils.GZipUnCompress(data)
	if gzerr != nil {
		return nil, fmt.Errorf("GZipUnCompress error:%v", gzerr)
	}

	checkReq := &scom.CheckReq{}
	umerr := json.Unmarshal(checkData, checkReq)
	if umerr != nil {
		return nil, fmt.Errorf("json.Unmarshal error:%v", umerr)
	}

	results, cherr := Check(checkReq.Tasks, time.Duration(checkReq.Timeout)*time.Second)
	if cherr != nil {
		return nil, fmt.Errorf("Check error:%v", cherr)
	}

	checkRsp := &scom.CheckRsp{Results: results}
	rspData, merr := json.Marshal(checkRsp)
	if merr != nil {
		return nil, fmt.Errorf("json.Marshal error:%v", merr)
	}

	gzRspData, gzCerr := mutils.GZipCompress(rspData)
	if gzCerr != nil {
		return nil, fmt.Errorf("GZipCompress error:%v", gzCerr)
	}

	return gzRspData, nil
}

func handle(mctx *mcom.MContext, msrv *mcom.MServer, data []byte) {
	switch {
	case mctx.OpCode() == scom.OPCODE_CHECK_REQ:
		rdata, rerr := handleCheck(mctx, msrv, data)
		msrv.Reply(mctx, rerr, rdata)
	default:
		mlog.Warnf("unknown opcode(%d), discard it!", mctx.OpCode())
	}
}
func errorNotify(mctx *mcom.MContext, msrv *mcom.MServer, err error) {
	mlog.Warnf("server disconnect...")
}

