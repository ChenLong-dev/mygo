/**
* @Author: cl
* @Date: 2020/12/31 16:59
 */
package main

import (
	"fmt"
	"github.com/ChenLong-dev/gobase/config"
	"github.com/ChenLong-dev/gobase/mlog"
	conf "myexample/config"
	"myexample/services"
	"os"
	"runtime"
)

var ostype = runtime.GOOS

func GetProjectPath() string {
	var projectPath string
	projectPath, _ = os.Getwd()
	fmt.Printf("projectPath:%s\n", projectPath)
	return projectPath
}

func GetConfigPath() string {
	path := GetProjectPath()
	if ostype == "windows" {
		path = path + "\\" + "config\\"
	} else if ostype == "linux" {
		path = path + "/" + "config/"
	}
	fmt.Printf("[GetConfigPath] path:%s\n", path)
	return path
}

func GetConLogPath() string {
	path := GetProjectPath()
	if ostype == "windows" {
		path = path + "\\log\\"
	} else if ostype == "linux" {
		path = path + "/log/"
	}
	fmt.Printf("[GetConLogPath] path:%s\n", path)
	return path
}


func main() {
	fmt.Println("test....")
	// 配置文件初始化
	config.InitExt(&conf.Cfg)
	fmt.Println("server start mode: ", config.Conf.AppMode)

	//日志初始化
	logCfg := conf.Cfg.Log
	_ = mlog.Init(&mlog.Params{
		Path:       logCfg.Path,
		MaxSize:    logCfg.MaxSize,
		MaxBackups: logCfg.MaxBackups,
		MaxAge:     logCfg.MaxAge,
		Level:      logCfg.Level,
	})
	mlog.Infof("test:%s", conf.Cfg.Mdbg.Host)

	err := services.InitMdbg()
	if err != nil {
		panic("InitMdbg is failed: " + err.Error())
	}
	select {

	}
}
