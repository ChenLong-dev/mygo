/**
* @Author: cl
* @Date: 2020/12/31 16:59
 */
package main

import (
	"fmt"
	"mygo/gobase/config"
	conf "myexample/config"
)

func main() {
	fmt.Println("test....")
	// 配置文件初始化
	config.InitExt(&conf.Cfg)
	fmt.Println("server start mode: ", conf.Cfg.AppMode)
}
