# 云眼可用性安装步骤

#### **1、拷贝安装包至安装服务器/tmp路径下**

#### **2、解压安装包：**tar -zxvf target_1594361816.tar.gz

```
[root@localhost package]# tar -zxvf target_1594363215.tar.gz
./target/
./target/tools/
./target/tools/mdbg
./target/supervisor/
./target/supervisor/install.sh
./target/supervisor/node/
./target/supervisor/node/eye_node.ini
./target/supervisor/server/
./target/supervisor/server/eye_server.ini
./target/supervisor/supervisord.conf
./target/supervisor/supervisord.service
./target/conf/
./target/conf/node/
./target/conf/node/server.toml
./target/conf/server/
./target/conf/server/server.toml
./target/detect/
./target/detect/node/
./target/detect/node/eye_node
./target/detect/server/
./target/detect/server/eye_server
./target/detect/md5sum.info
```

以上为解压内容

#### 3、检测服务安装

##### 1）解压当前目录并进入target/detect目录

##### 2）执行检测服务安装脚本：sh install.sh eye_server

```
[root@localhost detect]# sh install.sh eye_server
 ################################# 安装检测服务 #################################
4f7220261c8a4f9afc1db1d49a4f9fbe  /virus/cloud_waf_detect/eyeserver/eye_server
 === EYESERVER install success ===
```

##### 3）拷贝配置文件至安装目录：cp ../conf/server/server.toml /virus/cloud_waf_detect/eyeserver/

```
[root@localhost detect]# cp ../conf/server/server.toml /virus/cloud_waf_detect/eyeserver/
```

##### 4）修改配置文件相关内容：

```
## server conf
[mongo]
Host = "mongodb://10.107.30.110:27017"
UserName = "root"
Password = "saas_mongodb_c_root123"
DbName = "SAAS_yy_detect"

[etcd]
# hosts = ["etcd_gz.com:2379", "etcd_sh.com:2379", "etcd_bj.com:2379"]
hosts = ["10.227.30.129:2379", "10.227.30.130:2379", "10.227.30.131:2379"]
username = "root"
password = "saas_etcd_root123"
etcdCert = "/saasdata/etcd/data/ssl/server.pem"
etcdCertKey = "/saasdata/etcd/data/ssl/server-key.pem"
etcdCa = "/saasdata/etcd/data/ssl/ca.pem"

[availconf]
UpdateTaskIntvl = 3600
CheckIntvl = 10
CheckTimeout = 7
KeepLastPointNum = 3
MaxFlyCheckCount = 20
EntryAddr = "0.0.0.0:36767"
EntryPassword = "123456789"
MdbgAddr = "0.0.0.0:26767"
DetectType = "acl"
PointBound = 3
RegionBound = 1
DetectLimit = 0
AlarmAlias = []
GrpcAddr = "10.227.63.72:36164"
AlarmIntvl = 3600
```

注：修改测试环境对应的mongo以及etcd参数，其他不用修改，数据库名固定为：SAAS_yy_detect

##### 5）supervisor启动：cd target/supervisor

```
[root@localhost supervisor]# sh install.sh
 ######################### 帮助 #########################
 #./install.sh {param}
 {param}:
       eye_node         : make acl node
       eye_server       : make acl server
 ######################### 帮助 #########################
```

启动命令：sh install.sh eye_server

#### 4、检测点安装：

##### 1）解压当前目录并进入target/detect目录

##### 2）执行检测服务安装脚本：sh install.sh eye_node

```
[root@localhost detect]# sh install.sh eye_node
 ################################# 安装检测节点 #################################
df48014676f6faa49a7e974a749d42d7  /virus/cloud_waf_detect/eye/eye_node
 === EYE node install success ===
```

##### 3）拷贝配置文件至安装目录：cp ../conf/server/server.toml /virus/cloud_waf_detect/eye/

```
[root@localhost detect]# cp ../conf/node/server.toml /virus/cloud_waf_detect/eye/
```

##### 4）修改配置文件相关内容：

```
## server conf
[availconf]
region = "cloud_node_bj"
loginAddr = "10.107.63.72:36767"
loginPassword = "123456789"
concurrent = 8000
detectType = "eye"
mdbgAddr = "127.0.0.1:16767"

```

注：需修改对应节点名称region，检测服务监听的IP和端口loginAddr，其他不用修改

##### 5）supervisor启动：cd target/supervisor

```
[root@localhost supervisor]# sh install.sh
 ######################### 帮助 #########################
 #./install.sh {param}
 {param}:
       eye_node         : make acl node
       eye_server       : make acl server
 ######################### 帮助 #########################
```

启动命令：sh install.sh eye_node

#### 5、supervisor

##### 1）supervsior安装

supervsior使用的python环境为python3.7

内网安装需要修改pip源（http://mirrors.sangfor.org/help/2018/01/13/pypi.html）

或直接下载安装：

```
pip3 install -i http://mirrors.sangfor.org/pypi/simple --trusted-host mirrors.sangfor.org supervisor
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

6、数据库添加

数据库名：SAAS_yy_detect