# 安装步骤

#### 1、拷贝安装包至安装服务器/tmp路径下

#### 2、解压安装包：tar -zxvf target_1594361816.tar.gz

```
[root@node1 package]# tar -zxvf target_1607936463.tar.gz
./target/
./target/pop_access_server/
./target/pop_access_server/conf/
./target/pop_access_server/conf/dev.toml
./target/pop_access_server/conf/prod.toml
./target/pop_access_server/conf/test.toml
./target/pop_access_server/install.sh
./target/pop_access_server/pop_access_server
./target/md5sum.info
./target/install.sh
```

以上为解压内容

#### 3、检测服务安装

##### 1）解压当前目录并进入target目录
```
[root@node1 package]# cd target
```

##### 2）执行检测服务安装脚本：sh install.sh -h
```
[root@node1 target]# sh install.sh
 ######################### install   #########################
 ######################### 帮助 #########################
 #./install.sh {param}
 {param}:
        -b         : 安装pop_access_server执行文件
        -d         : 安装pop_access_server执行文件和安装或替换dev配置文件
        -t         : 安装pop_access_server执行文件和安装或替换test配置文件
        -p         : 安装pop_access_server执行文件和安装或替换prod配置文件
        -a         : 安装pop_access_server执行文件和所有配置文件
 ######################### 帮助 #########################

```
##### 3）安装：sh install.sh -t
```
[root@node1 target]# sh install.sh -t
 ######################### install -t  #########################
f41e464698e41aa71c1dfc248199b561  /data/sa/service/pop_access_server/pop_access_server
 === /data/sa/service/pop_access_server/pop_access_server install success ===
2d087f01e1127b7159a9d054525357d2  /data/sa/service/pop_access_server/conf/test.toml
 === /data/sa/service/pop_access_server/conf/test.toml install success ===

```
##### 4）设置环境变量，如测试环境test：
```
[root@node1 target]# export SAAS_COMMON_CONFIG="test"
```



##### 5）修改和确认配置文件相关内容：
```
# This is a TOML document.
# export SAAS_COMMON_CONFIG="test"

#dev | test | beta | prod，应用运行模式
AppMode = "test"

[Mdbg]
Host = "0.0.0.0"
Port = 12345
Enable = 1  # 1:开启  0：关闭

# PopServer 服务端口配置
[PopServer]
Host = "10.107.88.5"
Port = 18400
SrvLockKey = "lockkey"
EtcdCluster = "services"
ProjectName = "v1.0.0"
LimitNum = 10

# ETCD 配置
[Etcd]
hosts = ["10.227.30.129:2379", "10.227.30.130:2379", "10.227.30.131:2379"]
username = "root"
password = "saas_etcd_root123"
etcdCert = "/saasdata/etcd/data/ssl/server.pem"
etcdCertKey = "/saasdata/etcd/data/ssl/server-key.pem"
etcdCa = "/saasdata/etcd/data/ssl/ca.pem"

# MongoDB配置
[Mongo]
Host = "mongodb://10.107.30.110:27017,10.107.30.111:27018,10.107.30.112:27019"
UserName = "root"
Password = "saas_mongodb_c_root123"
DbName = "SA_POP_ACCESS_SERVER"

# Redis配置
[Redis]
Network = "127.0.0.1"
Host = ""
Password = ""

# 日志相关配置
[Log]
Path = "/var/logs/sa/pop_access_server/pop_access_server.log"
Level = 2           # 日志级别 [1:TRACE 跟踪,2:DEBUG 调试,3:INFO 信息,4:WARN 警告,5:ERROR 一般错误,6:FATAL 致命错误]
MaxSize = 2         # MB
MaxBackups = 3      # 备份个数
MaxAge = 10         # 保存时间,天

```

注：修改测试环境对应的mongo以及etcd参数，其他不用修改，数据库名固定为：SAAS_yy_detect

##### 6）supervisor启动：cd target/supervisor

```
[root@node1 supervisor]# sh  install.sh
 ################################# 安装检测supervisor服务 #################################
WARNING: Running pip install with root privileges is generally not a good idea. Try `pip3 install --user` instead.
Looking in indexes: http://mirrors.xxx.org/pypi/simple, http://mirrors.xxx.org/sfpypi/simple
Requirement already satisfied: supervisor in /usr/local/lib/python3.8/site-packages (4.2.1)
 === /usr/bin/supervisord install success ===
 /etc/supervisor/ had exist
 /etc/supervisor/log had exist
 /var/run/ had exist
 /var/log/ had exist
 /usr/lib/systemd/system/supervisord.service had exist
 === 开机启动配置 load success ===
 ################################# 安装检测pop_access_server服务 #################################
 /data/supervisor/pop_access_server.ini had exist
 === pop_access_server.ini load success ===
 supervisord had exist!
root      1780  1728  0 20:28 pts/0    00:00:00 grep supervisord
root      2145     1  0 19:45 ?        00:00:00 /usr/bin/python3.8 /usr/local/bin/supervisord -c /etc/supervisor/supervisord.conf
 #################### install supervisor done #####################
```

启动命令：sh install.sh

#### 5、supervisor

##### 1）supervsior安装

supervsior使用的python环境为python3.7以上

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

6、数据库添加

数据库名：popaccessserver

由于有事务，必须先创建以下两个表：   
```$xslt
db.createCollection("TenantVersionRecord")   
db.createCollection("TenantConfigRecord")
```


7、ETCD相关设置

##### 1）popaccessserver MongoDB配置文件设置路径(key)及设置（value）
```$xslt
/saas_config/popaccessserver/MongoDB/

{
    "host": "mongodb://10.107.30.110:27017,10.107.30.111:27018,10.107.30.112:27019",
    "username": "root",
    "password": "saas_mongodb_c_root123",
    "dbname": "popaccessserver",
    "authdb": "admin"
}
```

##### 2）popaccessserver注册路径(key)及设置（value）
```$xslt
/services/v1.0.0/pop_access_server/grpc/10.107.88.10:18400

"10.107.88.10:18400"
```

##### 3）popaccessserver服务lockkey
```$xslt
/services/v1.0.0/pop_access_server/lockkey

占用
```

##### 4）etcd ssl文件
```$xslt
hosts = ["10.107.30.165:2379"]
username = "root"
salt = "SA1.0.0"
etcdCert = "/saasdata/etcd/data/ssl/server.pem"
etcdCertKey = "/saasdata/etcd/data/ssl/server-key.pem"
etcdCa = "/saasdata/etcd/data/ssl/ca.pem"
```