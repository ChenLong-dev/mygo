#!/bin/bash

SHELLPATH=$(cd `dirname $0`; pwd)
cd ${SHELLPATH}

SUP="/virus/supervisor"
EYE="/virus/cloud_waf_detect/eye"
EYESERVER="/virus/cloud_waf_detect/eyeserver"

function check_process() {
    count=`ps -ef |grep $1 |grep -v "grep" |wc -l`
    #echo $count
    if [ 0 == $count ];then
        return 0
    fi
    return 1
}

function install_eye_node() {
    if [ ! -d ${SUP}/log ]; then
        mkdir -p ${SUP}/log
    fi
    if [ ! -f ${SUP}/eye_node.ini ]; then
        cp ./node/eye_node.ini /virus/supervisor/
    else
        echo -e "\033[33m ${SUP}/eye_node.ini had exist \033[0m"
    fi
    echo -e "\033[32m === eye_node.ini load success === \033[0m"
}

function install_detect_node() {
    #supervisor安装、配置加载
    install_supervisor

    echo -e "\033[32m ################################# 安装检测节点 #################################\033[0m"
    case $1 in
    "eye_node")
        install_eye_node
    ;;
    "all")
        install_eye_node
    ;;
    esac

    #启动supervisor
    check_process supervisord
    if [ $? == 0 ]; then
        supervisord -c /etc/supervisor/supervisord.conf
    else
        echo "supervisord had exist"
    fi
    ps -ef|grep supervisord
    echo '#################### install supervisor done #####################'
}

function install_eye_server() {
    if [ ! -d ${SUP}/log ]; then
        mkdir -p ${SUP}/log
    fi
    if [ ! -f ${SUP}/eye_server.ini ]; then
        cp ./server/eye_server.ini /virus/supervisor/
    else
        echo -e "\033[33m ${SUP}/eye_server.ini had exist \033[0m"
    fi
    echo -e "\033[32m === eye_server.ini load success === \033[0m"
}

function install_detect_server() {
    #supervisor安装、配置加载
    install_supervisor

    echo -e "\033[32m ################################# 安装检测服务 #################################\033[0m"
    case $1 in
    "eye_server")
        install_eye_server
    ;;
    "all")
        install_eye_server
    ;;
    esac

    #启动supervisor
    check_process supervisord
    if [ $? == 0 ]; then
        supervisord -c /etc/supervisor/supervisord.conf
    else
        echo "supervisord had exist"
    fi
    ps -ef|grep supervisord
    echo '#################### install supervisor done #####################'
}

function print_help {
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
    echo -e "\033[35m #./install.sh {param} \033[0m"
    echo -e "\033[35m {param}: \033[0m"
    echo -e "\033[35m       eye_node         : make acl node \033[0m"
    echo -e "\033[35m       eye_server       : make acl server \033[0m"
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
}

function install_supervisor {
    # 守护进程supervisor安装
    if [ ! -f /usr/bin/supervisord ]; then
        pip3 install supervisor
        ln -s /usr/local/python3/bin/supervisord   /usr/bin/supervisord
        ln -s /usr/local/python3/bin/supervisorctl /usr/bin/supervisorctl
    else
        echo -e "\033[33m /usr/bin/supervisord had exist \033[0m"
    fi
    echo -e "\033[32m === /usr/bin/supervisord install success === \033[0m"

    # 守护进程supervisor配置
    if [ ! -d "/etc/supervisor/" ];then
        mkdir -m 755 -p /etc/supervisor/
        cp supervisord.conf /etc/supervisor/
    else
        echo -e "\033[33m /etc/supervisor/ had exist \033[0m"
    fi

    # 守护进程supervisor log目录
    if [ ! -d "/etc/supervisor/log" ];then
        mkdir -m 755 -p /etc/supervisor/log
    else
        echo -e "\033[33m /etc/supervisor/log had exist \033[0m"
    fi

    if [ ! -d "/var/run/" ];then
        mkdir -m 777 -p /var/run
    else
        echo -e "\033[33m /var/run/ had exist \033[0m"
    fi

    if [ ! -d "/var/log/" ];then
        mkdir -m 777 -p /var/log
    else
        echo -e "\033[33m /var/log/ had exist \033[0m"
    fi

    # 开机启动配置
    if [ ! -f /usr/lib/systemd/system/supervisord.service ]; then
        cp supervisord.service /usr/lib/systemd/system/
        systemctl enable supervisord
        systemctl is-enabled supervisord
    else
        echo -e "\033[33m /usr/lib/systemd/system/supervisord.service had exist \033[0m"
    fi
    echo -e "\033[32m === 开机启动配置 load success === \033[0m"
}

function main() {
    case $1 in
    "eye_node")
        install_detect_node eye_node
    ;;
    "eye_server")
        install_detect_server eye_server
    ;;
    "help")
        print_help
    ;;
    *)
      print_help
    ;;
    esac
}

main $@