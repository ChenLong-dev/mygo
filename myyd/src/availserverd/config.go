/**
* @Author: cl
* @Date: 2021/1/16 11:08
 */
package main

type Config struct {
	UpdateTaskIntvl  uint64
	CheckIntvl       uint64
	CheckTimeout     uint64
	KeepLastPointNum int
	MaxFlyCheckCount int64
	EntryAddr        string
	EntryPassword    string
	MdbgAddr         string
	DetectType       string
	PointBound       int
	RegionBound      int
	DetectLimit      int
	AlarmAlias       []string
	GrpcAddr         string
	MongoDB          *MongoDB
	MQ               *ConfigMq
	AlarmTemplate    string
}

type MongoDB struct {
	Host     string
	UserName string
	Password string
	DbName   string
}

type ConfigMq struct {
	TaskName     string
	QueueName    string
	RoutingKey   string
	ExchangeName string
	ExchangeType string
	BrokerUrl    string
}

////online
//var gConfig = &Config{
//	UpdateTaskIntvl:  60,
//	CheckIntvl:       10,
//	CheckTimeout:     7,
//	KeepLastPointNum: 10,
//	MaxFlyCheckCount: 20,
//	EntryAddr:        "0.0.0.0:36161",
//	EntryPassword:    "123456789",
//	MdbgAddr:         "0.0.0.0:26161",
//	DetectType:       "acl",
//	PointBound:       3,
//	RegionBound:      1,
//	DetectLimit:      0,
//	AlarmAlias:       []string{"12808","14424"},
//	GrpcAddr:         "121.46.30.195:36164",
//	AlarmTemplate:    "DetectType:[%s], DetectUrl:[%s], AlarmTime:[%s], AlarmReason:[%s]",
//	MongoDB: &MongoDB{
//		Host:     "mongodb://120.132.99.111:27017",
//		UserName: "xxx",
//		Password: "sangfor_123!Q@W#E",
//		DbName:   "SAAS_YD_DETECT",
//	},
//	MQ: &ConfigMq{
//		TaskName:     "saasyd.detect.task",
//		QueueName:    "saasyd_detect_queue",
//		RoutingKey:   "saasyd_detect_queue",
//		ExchangeName: "saasyd_detect_exchange",
//		ExchangeType: "direct",
//		BrokerUrl:    "amqp://root:sangfor123@121.46.4.209:5672/",
//	},
//}

//develop
var gConfig = &Config{
	UpdateTaskIntvl:  60,
	CheckIntvl:       10,
	CheckTimeout:     7,
	KeepLastPointNum: 10,
	MaxFlyCheckCount: 20,
	EntryAddr:        "0.0.0.0:36161",
	EntryPassword:    "123456789",
	MdbgAddr:         "0.0.0.0:26161",
	DetectType:       "acl",
	PointBound:       3,
	RegionBound:      1,
	DetectLimit:      0,
	AlarmAlias:       []string{"12808","14424"},
	GrpcAddr:         "121.46.30.195:36164",
	MongoDB: &MongoDB{
		Host:     "mongodb://10.227.63.170:27017",
		UserName: "root",
		Password: "sangfor123",
		DbName:   "saas_dc_detect",
	},
	MQ: &ConfigMq{
		TaskName:     "saasyd.detect.task",
		QueueName:    "saasyd_detect_queue",
		RoutingKey:   "saasyd_detect_queue",
		ExchangeName: "saasyd_detect_exchange",
		ExchangeType: "direct",
		BrokerUrl:    "amqp://root:sangfor123@10.227.63.170:5672/saasyd_mq_vhost",
	},
}

////test
//var gConfig = &Config{
//	UpdateTaskIntvl:  60,
//	CheckIntvl:       10,
//	CheckTimeout:     7,
//	KeepLastPointNum: 10,
//	MaxFlyCheckCount: 20,
//	EntryAddr:        "0.0.0.0:36161",
//	EntryPassword:    "123456789",
//	MdbgAddr:         "0.0.0.0:26161",
//	DetectType:       "acl",
//	PointBound:       3,
//	RegionBound:      1,
//	DetectLimit:      0,
//	AlarmAlias:       []string{"12808","14424"},
//	GrpcAddr:         "121.46.30.195:36164",
//	MongoDB: &MongoDB{
//		Host:     "mongodb://10.226.198.254:27023",
//		UserName: "root",
//		Password: "sangfor_saas_mongodb_root123",
//		DbName:   "SAAS_YD_detect",
//	},
//	MQ: &ConfigMq{
//		TaskName:     "saasyd.detect.task",
//		QueueName:    "saasyd_detect_queue",
//		RoutingKey:   "saasyd_detect_queue",
//		ExchangeName: "saasyd_detect_exchange",
//		ExchangeType: "direct",
//		BrokerUrl:    "amqp://root:sangfor_mq_root123@10.226.198.117:5672/saasyd_mq_vhost",
//	},
//}
