/**
* @Author: cl
* @Date: 2021/1/4 16:07
 */
package main

import (
	"flag"
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase/mcom"
	"github.com/ChenLong-dev/gobase/sspb"
	"github.com/golang/protobuf/proto"
	"os"
)

func main() {
	ip := flag.String("i", "127.0.0.1", "connect to ip")
	port := flag.Uint("p", 0, "connect to port")
	cmd := flag.String("c", "help", "command name")
	arg := flag.String("a", "", "command arg")

	flag.Parse()

	req := &sspb.PBMDbgReq{}
	req.Cmdname = *cmd
	req.Cmdreq = *arg
	req.Uid = int32(os.Getuid())
	//req.Uname = os.Getuna

	addr := fmt.Sprintf("%s:%d", *ip, *port)
	msrv, derr := mcom.MServerDial(mcom.NewMContext(), "tcp", addr, nil, nil, nil)
	if derr != nil {
		fmt.Fprintf(os.Stderr, "MServerDial(%s) error:%v\n", addr, derr)
		os.Exit(1)
	}

	dataRsp, rerr := msrv.Rpc(mcom.NewMContext(), mcom.MakeMHead(int32(sspb.PBOpCode_OPC_MBASE_MDBG_REQ)), req, 0)
	if rerr != nil {
		fmt.Fprintf(os.Stderr, "rpc(%v) error:%v\n", req, rerr)
		os.Exit(1)
	}
	rsp := &sspb.PBMDbgRsp{}
	proto.Unmarshal(dataRsp, rsp)

	fmt.Printf("%s\n", rsp.Cmdrsp)
	os.Exit(int(rsp.Result))
}