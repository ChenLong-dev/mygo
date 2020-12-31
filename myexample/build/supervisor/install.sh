#!/bin/bash
#@auth cl
#@time 20201214

TARGET=myexample

SUPERPATH="/data/supervisor"
ROOTPATH=/data/sa/service
INSTALLPATH=${ROOTPATH}/${TARGET}


function install_supervisor {
    echo -e "\033[32m ################################# 安装检测supervisor服务 #################################\033[0m"
    # 守护进程supervisor安装
    if [ ! -f /usr/local/bin/supervisord ] && [ ! -f /usr/bin/supervisord ]; then
        pip3 install supervisor
        #ln -s /usr/local/python3/bin/supervisord   /usr/bin/supervisord
        #ln -s /usr/local/python3/bin/supervisorctl /usr/bin/supervisorctl
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

function check_process() {
    count=`ps -ef |grep $1 |grep -v "grep" |wc -l`
    if [ 0 == $count ];then
        return 0
    fi
    return 1
}

function install_service() {
    if [ ! -d ${SUPERPATH}/log ]; then
        mkdir -p ${SUPERPATH}/log
    fi
    if [ ! -f ${SUPERPATH}/${TARGET}.ini ]; then
        cp ./server/${TARGET}.ini ${SUPERPATH}
    else
        echo -e "\033[33m ${SUPERPATH}/${TARGET}.ini had exist \033[0m"
    fi
    echo -e "\033[32m === ${TARGET}.ini load success === \033[0m"
}

function main() {
    #supervisor安装、配置加载
    install_supervisor

    echo -e "\033[32m ################################# 安装检测${TARGET}服务 #################################\033[0m"
    install_service
     #启动supervisor
    check_process supervisord
    if [ $? == 0 ]; then
        supervisord -c /etc/supervisor/supervisord.conf
    else
        echo -e "\033[33m supervisord had exist! \033[0m"
    fi
    ps -ef|grep supervisord
    echo -e "\033[32m #################### install supervisor done #####################\033[0m"
}

main $@