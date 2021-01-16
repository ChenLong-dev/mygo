#!/bin/bash
#@auth alexchen
#@time 20200603

SHELLPATH=$(cd `dirname $0`; pwd)
cd ${SHELLPATH}

EYE="/virus/cloud_waf_detect/eye"
EYESERVER="/virus/cloud_waf_detect/eyeserver"

function install_eye_node() {
    if [ ! -d ${EYE} ]; then
        mkdir -p ${EYE}
    else
        rm -rf ${EYE}/eye_node
    fi
    chmod +x node/eye_node
    cp node/eye_node ${EYE}
    md5sum ${EYE}/eye_node
    echo -e "\033[32m === EYE node install success === \033[0m"
}

function install_detect_node() {
    echo -e "\033[32m ################################# 安装检测节点 #################################\033[0m"
    case $1 in
    "eye")
        install_eye_node
    ;;
    "all")
        install_eye_node
    ;;
    esac
}

function install_eye_server() {
    if [ ! -d ${EYESERVER} ]; then
        mkdir -p ${EYESERVER}
    else
        rm -rf ${EYESERVER}/acl_server
    fi
    chmod +x server/eye_server
    cp server/eye_server ${EYESERVER}
    md5sum ${EYESERVER}/eye_server
    echo -e "\033[32m === EYESERVER install success === \033[0m"
}

function install_detect_server() {
    echo -e "\033[32m ################################# 安装检测服务 #################################\033[0m"
    case $1 in
    "eye")
        install_eye_server
    ;;
    "all")
        install_eye_server
    ;;
    esac
}

function print_help {
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
    echo -e "\033[35m #./install.sh {param} \033[0m"
    echo -e "\033[35m {param}: \033[0m"
    echo -e "\033[35m       eye_node        : make eye node \033[0m"
    echo -e "\033[35m       eye_server      : make eye server \033[0m"
    echo -e "\033[35m       all             : make all node/server \033[0m"
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
}

function main() {
    case $1 in
    "eye_node")
        install_detect_node eye
    ;;
    "eye_server")
        install_detect_server eye
    ;;
    "all")
        install_detect_node all
        install_detect_server all
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
