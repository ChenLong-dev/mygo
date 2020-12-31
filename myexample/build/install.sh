#!/bin/bash
#@auth cl
#@time 20201214

TARGET=myexample

VERSION=`date +%s`
SHELLPATH=$(cd `dirname $0`; pwd)
cd ${SHELLPATH}
FILEPATH=${SHELLPATH}/${TARGET}
ROOTPATH=/data/sa/service
INSTALLPATH=${ROOTPATH}/${TARGET}
INSTALLFILE=${INSTALLPATH}/${TARGET}
SUPERPATH=${SHELLPATH}/supervisor
DSTSUPERPATH="/data/supervisor"

function install_usr() {
    USRPATH=${SHELLPATH}/usr
    if [ ! -f /usr/local/bin/supervisord ] && [ ! -f /usr/bin/supervisord ]; then
        cp ${USRPATH}/local/bin/supervisord /usr/local/bin
        md5sum /usr/local/bin/supervisord
        chmod +x /usr/local/bin/supervisord
    fi

    if [ ! -f /usr/local/bin/supervisorctl ] && [ ! -f /usr/bin/supervisorctl ]; then
        cp ${USRPATH}/local/bin/supervisorctl /usr/local/bin
        md5sum /usr/local/bin/supervisorctl
        chmod +x /usr/local/bin/supervisorctl
    fi
    ls -al /usr/local/bin/supervisor*
    echo -e "\033[32m === usr install success === \033[0m"
}

function install_tools() {
     if [ ! -f /usr/local/bin/mdbg ] && [ ! -f /usr/bin/mdbg ]; then
        cp ${USRPATH}/local/bin/supervisord /usr/local/bin
        md5sum /usr/local/bin/supervisord
        chmod +x /usr/local/bin/supervisord
    fi
}

function install_service() {
    if [ ! -d ${INSTALLPATH} ]; then
        mkdir -p ${INSTALLPATH}
    else
        if [ -f ${INSTALLFILE} ]; then
            mv ${INSTALLFILE} ${INSTALLFILE}_bak_${VERSION}
        fi
    fi
    chmod +x ${FILEPATH}/${TARGET}
    cp ${FILEPATH}/${TARGET} ${INSTALLPATH}
    if [ $? != 0 ]; then
        return 1
    fi
    md5sum ${INSTALLFILE}
    echo -e "\033[32m === ${INSTALLFILE} install success === \033[0m"
}

function install_config() {
    FILETYPE=$1
    SRCCONFPATH=${FILEPATH}/conf
    SRCCONFPATH=${SRCCONFPATH}/${FILETYPE}.toml
    DSTCONFPATH=${INSTALLPATH}/conf
    DSTCONFFILE=${DSTCONFPATH}/${FILETYPE}.toml
    if [ ! -d ${DSTCONFPATH} ]; then
        mkdir -p ${DSTCONFPATH}
    else
        if [ -f ${DSTCONFFILE} ]; then
            mv ${DSTCONFFILE} ${DSTCONFFILE}_bak_${VERSION}
        fi
    fi
    cp ${SRCCONFPATH} ${DSTCONFPATH}
    if [ $? != 0 ]; then
        return 1
    fi
    md5sum ${DSTCONFFILE}
    echo -e "\033[32m === ${DSTCONFFILE} install success === \033[0m"
}

function check_process() {
    count=`ps -ef |grep $1 |grep -v "grep" |wc -l`
    if [ 0 == $count ];then
        return 0
    fi
    return 1
}

function install_supervisor() {
    check_process supervisord
    if [ $? == 0 ]; then
        sh ${SUPERPATH}/install.sh
    fi
    INIPATH=${SUPERPATH}/server
    if [ ! -f ${DSTSUPERPATH}/${TARGET}.ini ]; then
        cp ${INIPATH}/${TARGET}.ini ${DSTSUPERPATH}
    else
        echo -e "\033[33m ${DSTSUPERPATH}/${TARGET}.ini had exist \033[0m"
    fi
    count=`supervisorctl status | grep ${TARGET} | wc -l`
    if [ 0 == $count ];then
        res1=`supervisorctl update`
        echo -e "\033[32m === res1 ${res1} update success === \033[0m"
        res2=`supervisorctl start ${TARGET}`
        echo -e "\033[32m === res2 ${res2} start success === \033[0m"
    else
        res3=`supervisorctl restart ${TARGET}`
        echo -e "\033[32m === res3 ${res3} restart success === \033[0m"
    fi
    status=`supervisorctl status | grep ${TARGET}`
    echo -e "\033[32m === ${status} === \033[0m"
}

function print_help {
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
    echo -e "\033[35m #sh ./install.sh {param} \033[0m"
    echo -e "\033[35m {param}: \033[0m"
    echo -e "\033[35m        -b         : 安装${TARGET}执行文件 \033[0m"
    echo -e "\033[35m        -d         : 安装${TARGET}执行文件和安装或替换dev配置文件 \033[0m"
    echo -e "\033[35m        -t         : 安装${TARGET}执行文件和安装或替换test配置文件 \033[0m"
    echo -e "\033[35m        -p         : 安装${TARGET}执行文件和安装或替换prod配置文件 \033[0m"
    echo -e "\033[35m        -a         : 安装${TARGET}执行文件和所有配置文件 \033[0m"
    echo -e "\033[35m        -s         : supervisor相关文件 \033[0m"
    echo -e "\033[35m        -u         : 安装usr相关文件 \033[0m"
    echo -e "\033[35m        -r         : 删除备份文件 \033[0m"
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
}

function main() {
    echo -e "\033[34m ######################### sh install.sh $1 ######################### \033[0m"
    case $1 in
        "-b")
            install_service ${TARGET}
        ;;
        "-d")
            install_service ${TARGET}
            install_config dev
        ;;
        "-t")
            install_service ${TARGET}
            install_config test
        ;;
        "-p")
            install_service ${TARGET}
            install_config prod
        ;;
        "-a")
            install_service ${TARGET}
            install_config dev
            install_config test
            install_config prod
        ;;
        "-u")
            install_usr
        ;;
        "-s")
            install_supervisor
        ;;
        "-r")
            rm -rf ${INSTALLFILE}_bak_*
            rm -rf ${INSTALLFILE}/conf/*_bak_*
        ;;
        "-h")
            print_help
        ;;
        *)
            print_help
        ;;
        esac
}

main $@