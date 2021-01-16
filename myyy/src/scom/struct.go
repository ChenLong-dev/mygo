package scom

import (
	"github.com/ChenLong-dev/gobase/mbase/mutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ISP int

const (
	DefaultISP   ISP = 0
	ChinaTelecom ISP = 1
	ChinaUnicom  ISP = 2
	ChinaMobile  ISP = 3
)

const (
	TELECOM = "telecom"
	UNICOM  = "unicom"
	MOBILE  = "mobile"
)

var AclUrlRes map[string]string //源站是否配置ACL（"noacl"：否，"acl"：是）

const (
	OPCODE_LOGIN_REQ int32 = 10001
	OPCODE_LOGIN_RSP int32 = 10002
	OPCODE_CHECK_REQ int32 = 10011
	OPCODE_CHECK_RSP int32 = 10012
	OPCODE_PING_REQ  int32 = 10021
	OPCODE_PING_RSP  int32 = 10022
)

type LoginReq struct {
	Region string `json:"Region,omitempty"`

	Rand      string `json:"Rand,omitempty"`
	Timestamp uint64 `json:"Timestamp,omitempty"`
	Signature string `json:"Signature,omitempty"`
}
type LoginRsp struct {
	Code   int    `json:"Code,omitempty"`
	Status string `json:"Status,omitempty"`
}

///////////////////////////////////////////////////////
type Line struct {
	Addr string `json:"a,omitempty"`
	Isp  ISP    `json:"i,omitempty"`
}
type Head struct {
	Key string `json:"k,omitempty"`
	Val string `json:"v,omitempty"`
}
type Task struct {
	Url    string `json:"url,omitempty"` //"http://mirrors.sangfor.org/"
	Heads  []Head `json:"hs,omitempty"`
	Method string `json:"m,omitempty"` //GET,HEAD
	//Timeout 		int 		`json:"t,omitempty"`		//timeout, in second
	PrimaryAddrs []Line `json:"pa,omitempty"`
	SecondAddrs  []Line `json:"sa,omitempty"`
}

func (task *Task) String() string {
	return mutils.JsonPrint(task)
}

type Result struct {
	Code   int           `json:"c,omitempty"` //	返回http状态码，若超时或网络错误返回负数
	Status string        `json:"s,omitempty"`
	Delay  time.Duration `json:"d,omitempty"`
}
type TaskResult struct {
	Primarys []Result `json:"ps,omitempty"`
	Seconds  []Result `json:"ss,omitempty"`
}

func (tr *TaskResult) String() string {
	return mutils.JsonPrint(tr)
}

type CheckReq struct {
	Timeout   int       `json:"timeout,omitempty"`
	Tasks     []*Task   `json:"tasks,omitempty"`
	CheckTime time.Time `json:"checktime,omitempty"`
}
type CheckRsp struct {
	Results []*TaskResult `json:"results,omitempty"`
}

///////////////////////////////////////////////////////
type SiteIP struct {
	ExtranetIp string `bson:"extranet_ip"`
	Type       string `bson:"type"`
}

type BusinessInfos struct {
	Id              primitive.ObjectID `bson:"_id"`
	RefClouduser    string             `bson:"ref_clouduser"`
	BusinessName    string             `bson:"business_name"`
	Domain          string             `bson:"domain"`
	DomainId        primitive.ObjectID `bson:"domain_id"`
	DnsType         int                `bson:"dns_type"`
	Ip              interface{}        `bson:"ip"`
	Port            int                `bson:"port"`
	Protocol        string             `bson:"protocol"`
	Url             string             `bson:"url"`
	RefClusterNode  string             `bson:"ref_cluster_node"`
	InsiteIp        []SiteIP           `bson:"insite_ip"`
	RefReserveNode  string             `bson:"ref_reserve_node"`
	ReserveInsiteIp SiteIP             `bson:"reserve_insite_ip"`
}

type DomainInfos struct {
	Id              primitive.ObjectID `bson:"_id"`
	Domain          string             `bson:"domain"`
	BusinessName    string             `bson:"business_name"`
	IsIn            int                `bson:"is_in"`
	DnsType         int                `bson:"dns_type"`
	Ip              interface{}        `bson:"ip"`
	RefServiceGroup primitive.ObjectID `bson:"ref_service_group"`
	InsiteIp        []SiteIP           `bson:"insite_ip"`
	RefReserveNode  string             `bson:"ref_reserve_node"`
	ReserveInsiteIp interface{}        `bson:"reserve_insite_ip"`
}

type AlarmReq struct {
	Type    string   `json:"type"`
	Content string   `json:"content"`
	ToAlias []string `json:"to_alias"`
}
type AlarmRsp struct {
	Result int    `json:"result"`
	ErrMsg string `json:"errmsg"`
}

type ClusterIPInfos struct {
	NodeName        string `bson:"node_name"`
	OwnerType       string `bson:"owner_type"`
	RefServiceGroup string `bson:"ref_service_group"`
	Ip              SiteIP `bson:"ip"`
}

type LineResult struct {
	Line   Line
	Result Result
}
type RegionResult struct {
	Region      string
	LineResults []LineResult
}
type AclResults struct {
	Url         string                   `bson:"url"`
	AclResult   string                   `bson:"acl_result"`
	Results     map[string]*RegionResult `bson:"results"`
	CheckerTime time.Time                `bson:"check_time"`
	CreateTime  time.Time                `bson:"create_time"`
}
type AclRes struct {
	Url       string
	AclResult string
}

type EyeInfo struct {
	Url     string
	Freq    int32
	Trigger time.Time
}