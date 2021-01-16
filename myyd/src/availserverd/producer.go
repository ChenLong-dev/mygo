/**
* @Author: cl
* @Date: 2021/1/16 11:16
 */
package main

import (
	"github.com/ChenLong-dev/gobase/mlog"
	"github.com/ChenLong-dev/gobase/mq"
	"myyd/src/availserverd/config"
	"time"
)

type RespPro struct {
	msg map[string]interface{}
}

// 实现发送者
func (t *RespPro) MsgContent() map[string]interface{} {
	return t.msg
}

func ProducerMq(typ string, normal, abnormal []string) {
	if len(normal) + len(abnormal) <=0 {
		return
	}
	vt1 := time.Now()
	defer func() {
		vt2 := time.Now()
		mlog.Infof("mq v=v=v=v [normal.len:%d] [abnormal.len:%d] [exp:%v]\n", len(normal), len(abnormal), vt2.Sub(vt1))
	}()

	//指定发送队列与任务名
	queueExchange := &mq.QueueExchange{
		config.Conf.Mq.TaskName,
		config.Conf.Mq.QueueName,
		config.Conf.Mq.RoutingKey,
		config.Conf.Mq.ExchangeName,
		config.Conf.Mq.ExchangeType,
	}
	// 发送消息到mq
	mq := mq.New(queueExchange, config.Conf.Mq.BrokerUrl)

	body := make(map[string]interface{})
	body["type"] = typ
	body["normal"] = normal
	body["abnormal"] = abnormal

	t := &RespPro{body,}
	mq.RegisterProducer(t)
	mq.Start()
}
// 实现接收者
func (t *RespPro) Consumer(dataByte []byte) error {
	mapData, err := mq.GetData(dataByte)
	if err != nil {
		return err
	}

	for key, value := range mapData {
		mlog.Infof("1 --- Consumer recv a msg: %+v, %+v\n", key, value)
	}

	return nil
}

func ConsummerMq(configMq *ConfigMq) {
	//指定发送队列与任务名
	queueExchange := &mq.QueueExchange{
		configMq.TaskName,
		configMq.QueueName,
		configMq.RoutingKey,
		configMq.ExchangeName,
		configMq.ExchangeType,
	}
	// 发送消息到mq
	//mq := mq.New(queueExchange, "amqp://root:sangfor123@10.227.63.170:5672/")
	mq := mq.New(queueExchange, configMq.BrokerUrl)

	t := &RespPro{}
	mq.RegisterReceiver(t)
	mq.Start()
}