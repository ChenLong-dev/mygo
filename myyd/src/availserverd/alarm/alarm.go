/**
* @Author: cl
* @Date: 2021/1/16 11:03
 */
package alarm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChenLong-dev/gobase/mlog"
	"io/ioutil"
	"myyd/src/availserverd/config"
	"myyd/src/scom"
	"net/http"
	"strings"
	"sync"
	"time"
)

var AlarmUrlMap map[string]int64
var lock sync.RWMutex

func init() {
	AlarmUrlMap = make(map[string]int64)
}

func SendAlarmMsg(key, content string) error {
	//url: https://api.kd77.cn/cgi-bin/im/send?access_token=9800010000000000102700000000000048514e5801008084
	//	body = {
	//		'type': 'text',
	//			'content':  msg,
	//			'to_alias': user_list
	//	}
	//
	//	content里面是告警内容
	//
	//	to_alias里面是工号列表
	lock.Lock()
	//过滤重复告警
	if val, ok := AlarmUrlMap[key]; ok {
		timestamp := time.Now().Unix()
		//mlog.Debugf("SendAlarmMsg timestamp:%v, val:%v, config.Conf.AvailConf.AlarmIntvl:%v", timestamp, val, config.Conf.AvailConf.AlarmIntvl)
		if timestamp < val+config.Conf.AvailConf.AlarmIntvl {
			lock.Unlock()
			return nil
		}
	}
	AlarmUrlMap[key] = time.Now().Unix()
	lock.Unlock()

	//组包发送MOA告警
	client := &http.Client{}
	//reqUrl := fmt.Sprintf("https://api.kd77.cn/cgi-bin/im/send?access_token=%s", "9800010000000000102700000000000048514e5801008084")
	reqUrl := fmt.Sprintf("https://api.kdzl.cn/cgi-bin/im/send?access_token=%s", "ce0f0100000000001027000000000000f3a7da71020050ab347a8f8c")
	body := scom.AlarmReq{
		Type:    "text",
		Content: content,
		ToAlias: config.Conf.AvailConf.AlarmAlias,
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(string(jsonBody)))
	if err != nil {
		mlog.Error(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := client.Do(req)
	if respErr != nil {
		//mlog.Error(err)
		return respErr
	}

	defer resp.Body.Close()
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		mlog.Error(err)
		return readErr
	}
	var rsp scom.AlarmRsp
	unmarshallErr := json.Unmarshal(respBody, &rsp)
	if unmarshallErr != nil {
		mlog.Error(err)
		return unmarshallErr
	}
	if rsp.Result != 0 {
		mlog.Errorf(rsp.ErrMsg)
		return errors.New(rsp.ErrMsg)
	}
	mlog.Infof("alarm [%s] [%s] [%v]", key, content, rsp)

	return nil
}

func SendAlarmMsgs(alarmInfo map[string]string) error {
	for url, content := range alarmInfo {
		SendAlarmMsg(url, content)
	}

	return nil
}
