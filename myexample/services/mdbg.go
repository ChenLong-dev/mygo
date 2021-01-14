/**
* @Author: cl
* @Date: 2021/1/14 14:43
 */
package services

import (
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase/mcom"
	"github.com/ChenLong-dev/gobase/mbase/mutils"
	"github.com/ChenLong-dev/gobase/mlog"
	conf "myexample/config"
)

func InitMdbg() (err error) {
	if conf.Cfg.Mdbg.Enable != 1 {
		mlog.Infof("mdbg is disabale [%d]", conf.Cfg.Mdbg.Enable)
		return
	}
	addr := fmt.Sprintf("%s:%d", conf.Cfg.Mdbg.Host, conf.Cfg.Mdbg.Port)
	mlog.Infof("mdbg is enabale [%d], addr:%s", conf.Cfg.Mdbg.Enable, addr)

	if addr == "" {
		return
	}
	_ = mcom.InitMdbg(addr)

	mcom.RegisterMdbgMml(mcom.MdbgPowerRoot, "example", "", mdbgExample)

	return nil
}

func mdbgExample(mml *mutils.Mml) (result int32, cmdRsp string) {
	name := mml.GetString("name", "chenlong")
	mlog.Infof("1====== [mdbgLastPoint:%s] =====1", name)
	return 0, "2====== [mdbgLastPoint] =====2"
}