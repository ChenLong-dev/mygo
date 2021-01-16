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

type Checker struct {
	Region 		string
	Msrv 		*mcom.MServer
}
func (checker *Checker) String() string {
	//return mutils.JsonPrint(checker)
	if checker == nil {
		return "{nil}"
	}
	return fmt.Sprintf("{Region:%s, Msrv:%v}", checker.Region, checker.Msrv.String())
}

func (checker *Checker) Close() {
	if checker.Msrv != nil {
		checker.Msrv.Close(mcom.NewMContext())
		checker.Msrv.Pridata = nil
	}
}

func (checker *Checker) Check(checkReq *scom.CheckReq, timeout time.Duration) (checkRsp *scom.CheckRsp, err error) {
	mlog.Tracef("checker(%v) timeout=%v", checker, timeout)
	defer func() {
		if err != nil {
			mlog.Warnf("Checker(%v) check return error:%v", checker, err)
		} else {
			mlog.Tracef("err=%v", err)
		}
	} ()

	checkData, merr := json.Marshal(checkReq)
	if merr != nil {
		return nil, fmt.Errorf("json.Marshal error:%v", merr)
	}
	gzCheckData, gzerr := mutils.GZipCompress(checkData)
	if gzerr != nil {
		return nil, fmt.Errorf("GZipCompress error:%v", gzerr)
	}
	//var rpcRsp []byte
	//var rpcErr error
	rpcRsp, rpcErr := checker.Msrv.Rpc(mcom.NewMContext(), mcom.MakeMHead(scom.OPCODE_CHECK_REQ), gzCheckData, int64(timeout/time.Millisecond))
	if rpcErr != nil {
		return nil, fmt.Errorf("rpc error:%v", rpcErr)
	}

	rsp, gzUerr := mutils.GZipUnCompress(rpcRsp)
	if gzUerr != nil {
		return nil, fmt.Errorf("GZipUnCompress error:%v", gzUerr)
	}
	checkRsp = &scom.CheckRsp{}
	umerr := json.Unmarshal(rsp, checkRsp)
	if umerr != nil {
		return nil, fmt.Errorf("Unmarshal error:%v", umerr)
	}
	return checkRsp, nil
}

func NewChecker(region string, msrv *mcom.MServer) *Checker {
	checker := &Checker{Msrv: msrv, Region: region}

	return checker
}