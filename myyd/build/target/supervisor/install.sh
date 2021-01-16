#!/bin/bash

SHELLPATH=$(cd `dirname $0`; pwd)
cd ${SHELLPATH}

SUP="/virus/supervisor"
ACL="/virus/cloud_waf_detect/acl"
LINK="/virus/cloud_waf_detect/link"
LINE="/virus/cloud_waf_detect/line"
ACLSERVER="/virus/cloud_waf_detect/aclserver"
LINKSERVER="/virus/cloud_waf_detect/linkserver"
LINESERVER="/virus/cloud_waf_detect/lineserver"
WEBSITE="/virus/cloud_waf_detect/saas_services/saasyd_website/"

function check_process() {
    count=`ps -ef |grep $1 |grep -v "grep" |wc -l`
    #echo $count
    if [ 0 == $count ];then
        return 0
    fi
    return 1
}

function install_acl_node() {
    if [ ! -d ${SUP}/log ]; then
        mkdir -p ${SUP}/log
    fi
    if [ ! -f ${SUP}/acl_node.ini ]; then
        cp ./node/acl_node.ini /virus/supervisor/
    else
        echo -e "\033[33m ${SUP}/acl_node.ini had exist \033[0m"
    fi
    echo -e "\033[32m === acl_node.ini load success === \033[0m"
}

function install_link_node() {
    if [ ! -d ${SUP}/log ]; then
        mkdir -p ${SUP}/log
    fi
    if [ ! -f ${SUP}/link_node.ini ]; then
        cp ./node/link_node.ini /virus/supervisor/
    else
        echo -e "\033[33m ${SUP}/link_node.ini had exist \033[0m"
    fi
    echo -e "\033[32m === link_node.ini load success === \033[0m"
}

function install_line_node() {
    if [ ! -d ${SUP}/log ]; then
        mkdir -p ${SUP}/log
    fi
    if [ ! -f ${SUP}/line_node.ini ]; then
        cp ./node/line_node.ini /virus/supervisor/
    else
        echo -e "\033[33m ${SUP}/line_node.ini had exist \033[0m"
    fi
    echo -e "\033[32m === line_node.ini load success === \033[0m"
}

function install_detect_node() {
    #supervisor安装、配置加载
    install_supervisor

    echo -e "\033[32m ################################# 安装检测节点 #################################\033[0m"
    case $1 in
    "acl")
        install_acl_node
    ;;
    "link")
        install_link_node
    ;;
    "line")
        install_line_node
    ;;
    "all")
        install_acl_node
        install_link_node
        install_line_node
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
    echo '#################### install node supervisor done #####################'
}

function install_acl_server() {
    if [ ! -d ${SUP}/log ]; then
        mkdir -p ${SUP}/log
    fi
    if [ ! -f ${SUP}/acl_server.ini ]; then
        cp ./server/acl_server.ini /virus/supervisor/
    else
        echo -e "\033[33m ${SUP}/acl_server.ini had exist \033[0m"
    fi
    echo -e "\033[32m === acl_server.ini load success === \033[0m"
}

function install_link_server() {
    if [ ! -d ${SUP}/log ]; then
        mkdir -p ${SUP}/log
    fi
    if [ ! -f ${SUP}/link_server.ini ]; then
        cp ./server/link_server.ini /virus/supervisor/
    else
        echo -e "\033[33m ${SUP}/link_server.ini had exist \033[0m"
    fi
    echo -e "\033[32m === link_server.ini load success === \033[0m"
}

function install_line_server() {
    if [ ! -d ${SUP}/log ]; then
        mkdir -p ${SUP}/log
    fi
    if [ ! -f ${SUP}/line_server.ini ]; then
        cp ./server/line_server.ini /virus/supervisor/
    else
        echo -e "\033[33m ${SUP}/line_server.ini had exist \033[0m"
    fi
    echo -e "\033[32m === line_server.ini load success === \033[0m"
}

function install_detect_server() {
    #supervisor安装、配置加载
    install_supervisor

    echo -e "\033[32m ################################# 安装检测服务 #################################\033[0m"
    case $1 in
    "acl")
        install_acl_server
    ;;
    "link")
        install_link_server
    ;;
    "line")
        install_line_server
    ;;
    "all")
        install_acl_server
        install_link_server
        install_line_server
    ;;
    esac

   #启动supervisor
    check_process supervisord
    if [ $? == 0 ]; then
        supervisord -c /etc/supervisor/supervisord.conf
    else
        echo -e "\033[33m supervisord had exist \033[0m"
    fi
    ps -ef|grep supervisord
    echo '#################### install supervisor done #####################'
}

function install_website_server() {
    #supervisor安装、配置加载
    install_supervisor
    if [ ! -d ${SUP}/log ]; then
        mkdir -p ${SUP}/log
    fi
    if [ ! -f ${SUP}/saasyd_website.ini ]; then
        cp ./server/saasyd_website.ini /virus/supervisor/
    else
        echo -e "\033[33m ${SUP}/saasyd_website.ini had exist \033[0m"
    fi
    echo -e "\033[32m === saasyd_website.ini load success === \033[0m"

    #启动supervisor
    check_process supervisord
    if [ $? == 0 ]; then
        supervisord -c /etc/supervisor/supervisord.conf
    else
        echo -e "\033[33m supervisord had exist \033[0m"
    fi
    ps -ef|grep supervisord
    echo '#################### install web server supervisor done #####################'
}

function print_help {
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
    echo -e "\033[35m #./install.sh {param} \033[0m"
    echo -e "\033[35m {param}: \033[0m"
    echo -e "\033[35m       acl_node        : make acl node \033[0m"
    echo -e "\033[35m       link_node       : make link node \033[0m"
    echo -e "\033[35m       line_node       : make line node \033[0m"
    echo -e "\033[35m       detect_node     : make all node \033[0m"
    echo -e "\033[35m       acl_server      : make acl server \033[0m"
    echo -e "\033[35m       line_server     : make link server \033[0m"
    echo -e "\033[35m       link_server     : make line server \033[0m"
    echo -e "\033[35m       detect_server   : make all server \033[0m"
    echo -e "\033[35m       website         : make website server \033[0m"
    echo -e "\033[35m       all             : make all node/server/website \033[0m"
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
    "acl_node")
        install_detect_node acl
    ;;
    "link_node")
        install_detect_node link
    ;;
    "line_node")
        install_detect_node line
    ;;
    "detect_node")
        install_detect_node all
    ;;
    "acl_server")
        install_detect_server acl
    ;;
    "link_server")
        install_detect_server link
    ;;
    "line_server")
        install_detect_server line
    ;;
    "detect_server")
        install_detect_server all
    ;;
    "all")
        install_detect_node all
        install_detect_server all
    ;;
    "website")
        install_website_server
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
