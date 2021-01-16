# 云眼可用性安装步骤

#### **1、拷贝安装包至安装服务器/tmp路径下**

#### **2、解压安装包：**tar -zxvf target_1595232657.tar.gz

```
[root@localhost package]# tar -zxvf target_1595232657.tar.gz
./target/
./target/detect/
./target/detect/install.sh
./target/detect/server.toml
./target/detect/node/
./target/detect/node/acl_node
./target/detect/node/link_node
./target/detect/node/line_node
./target/detect/server/
./target/detect/server/acl_server
./target/detect/server/link_server
./target/detect/server/line_server
./target/detect/server/saasyd_website
./target/detect/md5sum.info
./target/tools/
./target/tools/mdbg
./target/supervisor/
./target/supervisor/install.sh
./target/supervisor/node/
./target/supervisor/node/acl_node.ini
./target/supervisor/node/line_node.ini
./target/supervisor/node/link_node.ini
./target/supervisor/server/
./target/supervisor/server/acl_server.ini
./target/supervisor/server/line_server.ini
./target/supervisor/server/link_server.ini
./target/supervisor/server/saasyd_website.ini
./target/supervisor/supervisord.conf
./target/supervisor/supervisord.service
./target/conf/
./target/conf/node/
./target/conf/node/server.toml
./target/conf/server/
./target/conf/server/server.toml
```

以上为解压内容

#### 3、检测服务安装

##### 1）解压当前目录并进入target/detect目录

##### 2）执行检测服务安装脚本：sh install.sh detect_server

```
[root@localhost detect]# sh install.sh
 ######################### 帮助 #########################
 #./install.sh {param}
 {param}:
       acl             : make acl node
       link            : make link node
       line            : make line node
       detect_node     : make all node
       aclserver       : make acl server
       lineserver      : make link server
       linkserver      : make line server
       detect_server   : make all server
       website         : make website server
       all             : make all node/server/website
 ######################### 帮助 #########################
```

##### 3）拷贝配置文件至安装目录：cp ../conf/server/server.toml /virus/cloud_waf_detect/acl_server/

```
[root@localhost detect]# cp ../conf/server/server.toml /virus/cloud_waf_detect/acl_server/

[root@localhost detect]# cp ../conf/server/server.toml /virus/cloud_waf_detect/line_server/

[root@localhost detect]# cp ../conf/server/server.toml /virus/cloud_waf_detect/link_server/
```

##### 4）修改配置文件相关内容：

acl_server:

```
[root@localhost cloud_waf_detect]# cat aclserver/server.toml
## server conf
[mongo]
Host = "mongodb://120.132.99.111:27017"
UserName = "saas_detect"
Password = "u7dT1GF7gN^}rp8=?Lp!"
DbName = "SAAS_yd_detect"

[mq]
TaskName = "saasyd.detect.task"
QueueName = "saasyd_detect_queue"
RoutingKey = "saasyd_detect_queue"
ExchangeName = "saasyd_detect_exchange"
ExchangeType = "direct"
BrokerUrl = "mqp://root:RXa8pZZgUBLCfFdd24KK@121.46.4.209:5672/"

[availconf]
UpdateTaskIntvl = 3600
CheckIntvl = 180
CheckTimeout = 10
KeepLastPointNum = 30
MaxFlyCheckCount = 20
EntryAddr = "0.0.0.0:36161"
EntryPassword = "123456789"
MdbgAddr = "0.0.0.0:26161"
DetectType = "acl"
PointBound = 3
RegionBound = 3
DetectLimit = 0
AlarmAlias = []
GrpcAddr = "121.46.30.195:36164"
AlarmIntvl = 1800
```

line_server:

```
[root@localhost cloud_waf_detect]# cat lineserver/server.toml
## server conf
[mongo]
Host = "mongodb://120.132.99.111:27017"
UserName = "saas_detect"
Password = "u7dT1GF7gN^}rp8=?Lp!"
DbName = "SAAS_yd_detect"

[mq]
TaskName = "saasyd.detect.task"
QueueName = "saasyd_detect_queue"
RoutingKey = "saasyd_detect_queue"
ExchangeName = "saasyd_detect_exchange"
ExchangeType = "direct"
BrokerUrl = "amqp://root:RXa8pZZgUBLCfFdd24KK@121.46.4.209:5672/saasyd_mq_vhost"

[availconf]
UpdateTaskIntvl = 3600
CheckIntvl = 10
CheckTimeout = 7
KeepLastPointNum = 3
MaxFlyCheckCount = 20
EntryAddr = "0.0.0.0:36163"
EntryPassword = "123456789"
MdbgAddr = "0.0.0.0:26163"
DetectType = "line"
PointBound = 3
RegionBound = 3
DetectLimit = 0
AlarmAlias = []
GrpcAddr = "121.46.30.195:36164"
AlarmIntvl = 86400

```

link_server:

```
[root@localhost cloud_waf_detect]# cat linkserver/server.toml
## server conf
[mongo]
Host = "mongodb://120.132.99.111:27017"
UserName = "saas_detect"
Password = "u7dT1GF7gN^}rp8=?Lp!"
DbName = "SAAS_yd_detect"

[mq]
TaskName = "saasyd.detect.task"
QueueName = "saasyd_detect_queue"
RoutingKey = "saasyd_detect_queue"
ExchangeName = "saasyd_detect_exchange"
ExchangeType = "direct"
BrokerUrl = "amqp://root:RXa8pZZgUBLCfFdd24KK@121.46.4.209:5672/saasyd_mq_vhost"

[availconf]
UpdateTaskIntvl = 3600
CheckIntvl = 10
CheckTimeout = 7
KeepLastPointNum = 10
MaxFlyCheckCount = 20
EntryAddr = "0.0.0.0:36162"
EntryPassword = "123456789"
MdbgAddr = "0.0.0.0:26162"
DetectType = "link"
PointBound = 3
RegionBound = 3
DetectLimit = 0
AlarmAlias = ["12808"]
GrpcAddr = "121.46.30.195:36164"
AclGrpcAddr = "127.0.0.1:36166"
AlarmIntvl = 600
#TargetUrl = ["http://doc.sre.ac.cn"]
TargetUrl = ["http://bj.detect.sre.ac.cn","http://hd.detect.sre.ac.cn","http://gz.detect.sre.ac.cn","http://fs.detect.sre.ac.cn","http://sz.detect.sre.ac.cn"]
```

注：修改测试环境对应的mongo以及mq参数，其他不用修改，数据库固定为：SAAS_yd_detect

##### 5）supervisor启动：cd target/supervisor

```
[root@localhost supervisor]# sh install.sh
 ######################### 帮助 #########################
 #./install.sh {param}
 {param}:
       acl_node        : make acl node
       link_node       : make link node
       line_node       : make line node
       detect_node     : make all node
       acl_server      : make acl server
       line_server     : make link server
       link_server     : make line server
       detect_server   : make all server
       website         : make website server
       all             : make all node/server/website
 ######################### 帮助 #########################
```

启动命令：sh install.sh detect_server

#### 4、检测点安装：

##### 1）解压当前目录并进入target/detect目录

##### 2）执行检测服务安装脚本：sh install.sh eye_node

```
[root@localhost detect]# sh install.sh
 ######################### 帮助 #########################
 #./install.sh {param}
 {param}:
       acl             : make acl node
       link            : make link node
       line            : make line node
       detect_node     : make all node
       aclserver       : make acl server
       lineserver      : make link server
       linkserver      : make line server
       detect_server   : make all server
       website         : make website server
       all             : make all node/server/website
 ######################### 帮助 #########################
```

##### 3）拷贝配置文件至安装目录：cp ../conf/server/server.toml /virus/cloud_waf_detect/acl/

```
[root@localhost detect]# cp ../conf/node/server.toml /virus/cloud_waf_detect/acl/

[root@localhost detect]# cp ../conf/node/server.toml /virus/cloud_waf_detect/line/

[root@localhost detect]# cp ../conf/node/server.toml /virus/cloud_waf_detect/link/
```

##### 4）修改配置文件相关内容：

acl_node:

```
xhunter:/virus/cloud_waf_detect #cat acl/server.toml
## server conf
[availconf]
region = "cloud_node_sh"
loginAddr = "121.46.4.209:36161"
loginPassword = "123456789"
concurrent = 4000
detectType = "acl"
mdbgAddr = "127.0.0.1:16161"
```



```
xhunter:/virus/cloud_waf_detect #cat acl/server.toml
## server conf
[availconf]
region = "cloud_node_bj"
loginAddr = "121.46.4.209:36161"
loginPassword = "123456789"
concurrent = 4000
detectType = "acl"
mdbgAddr = "127.0.0.1:16161"
```



```
xhunter:/virus/cloud_waf_detect #cat acl/server.toml
## server conf
[availconf]
region = "cloud_node_gz"
loginAddr = "121.46.4.209:36161"
loginPassword = "123456789"
concurrent = 4000
detectType = "acl"
mdbgAddr = "127.0.0.1:16161"

```

注：需修改对应节点名称region，检测服务监听的IP和端口loginAddr，其他不用修改

##### 5）supervisor启动：cd target/supervisor

```
[root@localhost supervisor]# sh install.sh
 ######################### 帮助 #########################
 #./install.sh {param}
 {param}:
       acl_node        : make acl node
       link_node       : make link node
       line_node       : make line node
       detect_node     : make all node
       acl_server      : make acl server
       line_server     : make link server
       link_server     : make line server
       detect_server   : make all server
       website         : make website server
       all             : make all node/server/website
 ######################### 帮助 #########################
```

启动命令：sh install.sh detect_node

#### 5、supervisor

##### 1）supervsior安装

supervsior使用的python环境为python3.7

内网安装需要修改pip源（http://mirrors.xxx.org/help/2018/01/13/pypi.html）

或直接下载安装：

```
pip3 install -i http://mirrors.xxx.org/pypi/simple --trusted-host mirrors.xxx.org supervisor
```

##### 2）supervsior相关命令

```
supervisorctl stop program_name  # 停止某一个进程，program_name 为 [program:x] 里的 x
supervisorctl start program_name  # 启动某个进程
supervisorctl restart program_name  # 重启某个进程
supervisorctl stop groupworker:  # 结束所有属于名为 groupworker 这个分组的进程 (start，restart 同理)
supervisorctl stop groupworker:name1  # 结束 groupworker:name1 这个进程 (start，restart 同理)
supervisorctl stop all  # 停止全部进程，注：start、restartUnlinking stale socket /tmp/supervisor.sock
、stop 都不会载入最新的配置文件
supervisorctl reload  # 载入最新的配置文件，停止原有进程并按新的配置启动、管理所有进程
supervisorctl update  # 根据最新的配置文件，启动新配置或有改动的进程，配置没有改动的进程不会受影响而重启
```

#### 6、数据库添加

数据库名：SAAS_yd_detect



#### 7、MQ添加

 [rabbit_localhost_2020-7-21.json](..\..\..\..\..\..\01_WorkSpace\2020年7月\rabbit_localhost_2020-7-21.json) 

rabbit_localhost_2020-7-21.json

```
{"rabbit_version":"3.6.15","users":[{"name":"guest","password_hash":"GqUhE1McM8o9cLCGyCZJOsBUhUvGkAXvW6epXCsFYXNVGE3+","hashing_algorithm":"rabbit_password_hashing_sha256","tags":"administrator"},{"name":"root","password_hash":"tpPqYfvECZ6u4ZsrzVvlMlWfrwxscZoZlt+4Q8mrCNIMLQTX","hashing_algorithm":"rabbit_password_hashing_sha256","tags":"administrator"}],"vhosts":[{"name":"saasdns_mq_vhost"},{"name":"saasyd_mq_vhost"},{"name":"/"}],"permissions":[{"user":"guest","vhost":"/","configure":".*","write":".*","read":".*"},{"user":"root","vhost":"saasyd_mq_vhost","configure":".*","write":".*","read":".*"},{"user":"root","vhost":"/","configure":".*","write":".*","read":".*"},{"user":"root","vhost":"saasdns_mq_vhost","configure":".*","write":".*","read":".*"}],"parameters":[],"global_parameters":[{"name":"cluster_name","value":"rabbit@localhost"}],"policies":[],"queues":[{"name":"celery@bind_task_sz.celery.pidbox","vhost":"saasdns_mq_vhost","durable":false,"auto_delete":true,"arguments":{"x-expires":10000,"x-message-ttl":300000}},{"name":"celeryev.0a849d9d-a86f-4db8-b5d5-65ddcceb9029","vhost":"saasdns_mq_vhost","durable":false,"auto_delete":true,"arguments":{"x-expires":60000,"x-message-ttl":5000}},{"name":"saasdns_bind_queue_gz","vhost":"saasdns_mq_vhost","durable":true,"auto_delete":false,"arguments":{"x-dead-letter-exchange":"dead-letter-exchange","x-dead-letter-routing-key":"dead"}},{"name":"saasdns_bind_queue_sz","vhost":"saasdns_mq_vhost","durable":true,"auto_delete":false,"arguments":{"x-dead-letter-exchange":"dead-letter-exchange","x-dead-letter-routing-key":"dead"}},{"name":"celeryev.513c2314-7fc4-464a-9826-64ebd28928cd","vhost":"saasdns_mq_vhost","durable":false,"auto_delete":true,"arguments":{"x-expires":60000,"x-message-ttl":5000}},{"name":"celery@bind_task_gz.celery.pidbox","vhost":"saasdns_mq_vhost","durable":false,"auto_delete":true,"arguments":{"x-expires":10000,"x-message-ttl":300000}},{"name":"saasyd_detect_queue","vhost":"saasyd_mq_vhost","durable":true,"auto_delete":false,"arguments":{"x-dead-letter-exchange":"dead-letter-exchange","x-dead-letter-routing-key":"dead"}},{"name":"celeryev.05337a16-9a7d-4f49-8f1b-5f1e91679c07","vhost":"saasyd_mq_vhost","durable":false,"auto_delete":true,"arguments":{"x-expires":60000,"x-message-ttl":5000}},{"name":"celery@link_statistics_task.celery.pidbox","vhost":"saasyd_mq_vhost","durable":false,"auto_delete":true,"arguments":{"x-expires":10000,"x-message-ttl":300000}},{"name":"celery@worker_drainage.celery.pidbox","vhost":"saasyd_mq_vhost","durable":false,"auto_delete":true,"arguments":{"x-expires":10000,"x-message-ttl":300000}},{"name":"statistics_tasks_beat_queue","vhost":"saasyd_mq_vhost","durable":true,"auto_delete":false,"arguments":{"x-dead-letter-exchange":"dead-letter-exchange","x-dead-letter-routing-key":"dead"}},{"name":"celeryev.d592f5d7-5f47-4704-8b48-d36fefea232f","vhost":"saasyd_mq_vhost","durable":false,"auto_delete":true,"arguments":{"x-expires":60000,"x-message-ttl":5000}},{"name":"celery","vhost":"saasyd_mq_vhost","durable":true,"auto_delete":false,"arguments":{}}],"exchanges":[{"name":"celeryev","vhost":"saasdns_mq_vhost","type":"topic","durable":true,"auto_delete":false,"internal":false,"arguments":{}},{"name":"celery.pidbox","vhost":"saasdns_mq_vhost","type":"fanout","durable":false,"auto_delete":false,"internal":false,"arguments":{}},{"name":"saasdns_bind_exchange","vhost":"saasdns_mq_vhost","type":"topic","durable":true,"auto_delete":false,"internal":false,"arguments":{}},{"name":"reply.celery.pidbox","vhost":"saasdns_mq_vhost","type":"direct","durable":false,"auto_delete":false,"internal":false,"arguments":{}},{"name":"saasyd_detect_exchange","vhost":"saasyd_mq_vhost","type":"direct","durable":true,"auto_delete":false,"internal":false,"arguments":{}},{"name":"celeryev","vhost":"saasyd_mq_vhost","type":"topic","durable":true,"auto_delete":false,"internal":false,"arguments":{}},{"name":"default_exchange","vhost":"saasyd_mq_vhost","type":"direct","durable":true,"auto_delete":false,"internal":false,"arguments":{}},{"name":"reply.celery.pidbox","vhost":"saasyd_mq_vhost","type":"direct","durable":false,"auto_delete":false,"internal":false,"arguments":{}},{"name":"celery.pidbox","vhost":"saasyd_mq_vhost","type":"fanout","durable":false,"auto_delete":false,"internal":false,"arguments":{}},{"name":"celery","vhost":"saasyd_mq_vhost","type":"direct","durable":true,"auto_delete":false,"internal":false,"arguments":{}}],"bindings":[{"source":"celery.pidbox","vhost":"saasdns_mq_vhost","destination":"celery@bind_task_gz.celery.pidbox","destination_type":"queue","routing_key":"","arguments":{}},{"source":"celery.pidbox","vhost":"saasdns_mq_vhost","destination":"celery@bind_task_sz.celery.pidbox","destination_type":"queue","routing_key":"","arguments":{}},{"source":"celeryev","vhost":"saasdns_mq_vhost","destination":"celeryev.0a849d9d-a86f-4db8-b5d5-65ddcceb9029","destination_type":"queue","routing_key":"worker.#","arguments":{}},{"source":"celeryev","vhost":"saasdns_mq_vhost","destination":"celeryev.513c2314-7fc4-464a-9826-64ebd28928cd","destination_type":"queue","routing_key":"worker.#","arguments":{}},{"source":"saasdns_bind_exchange","vhost":"saasdns_mq_vhost","destination":"saasdns_bind_queue_gz","destination_type":"queue","routing_key":"saasdns.bind.*","arguments":{}},{"source":"saasdns_bind_exchange","vhost":"saasdns_mq_vhost","destination":"saasdns_bind_queue_sz","destination_type":"queue","routing_key":"saasdns.bind.*","arguments":{}},{"source":"celery","vhost":"saasyd_mq_vhost","destination":"celery","destination_type":"queue","routing_key":"celery","arguments":{}},{"source":"celery.pidbox","vhost":"saasyd_mq_vhost","destination":"celery@link_statistics_task.celery.pidbox","destination_type":"queue","routing_key":"","arguments":{}},{"source":"celery.pidbox","vhost":"saasyd_mq_vhost","destination":"celery@worker_drainage.celery.pidbox","destination_type":"queue","routing_key":"","arguments":{}},{"source":"celeryev","vhost":"saasyd_mq_vhost","destination":"celeryev.05337a16-9a7d-4f49-8f1b-5f1e91679c07","destination_type":"queue","routing_key":"worker.#","arguments":{}},{"source":"celeryev","vhost":"saasyd_mq_vhost","destination":"celeryev.d592f5d7-5f47-4704-8b48-d36fefea232f","destination_type":"queue","routing_key":"worker.#","arguments":{}},{"source":"default_exchange","vhost":"saasyd_mq_vhost","destination":"statistics_tasks_beat_queue","destination_type":"queue","routing_key":"statistics_tasks_beat_queue","arguments":{}},{"source":"saasyd_detect_exchange","vhost":"saasyd_mq_vhost","destination":"saasyd_detect_queue","destination_type":"queue","routing_key":"saasyd_detect_queue","arguments":{}}]}
```