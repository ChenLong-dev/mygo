/**
* @Author: cl
* @Date: 2021/1/16 11:20
 */
package main

import (
	"encoding/json"
	"github.com/ChenLong-dev/gobase/mbase/mutils"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyd/src/scom"
)

func getZiboTasks() (tasks []*scom.Task) {
	ziboJson := `[{"url":"http://jgswglj.zibo.gov.cn:80","ip":["120.220.22.15"]}]`
	var ziboList []map[string]interface{}
	err := json.Unmarshal([]byte(ziboJson), &ziboList)
	if err != nil {
		mlog.Errorf("json.Unmarshal error:%v", err)
		return nil
	}
	for _, zurl := range ziboList {
		task := &scom.Task{}
		task.Url = zurl["url"].(string)
		task.Method = "HEAD"
		if ip, ok := zurl["ip"].(string); ok {
			task.PrimaryAddrs = append(task.PrimaryAddrs, scom.Line{Addr: ip})
		} else {
			for _, sip := range zurl["ip"].([]interface{}) {
				task.PrimaryAddrs = append(task.PrimaryAddrs, scom.Line{Addr: sip.(string)})
			}
		}

		tasks = append(tasks, task)
	}

	mlog.Tracef("tasks=%v", mutils.JsonPrintPretty(tasks))

	return tasks
}
