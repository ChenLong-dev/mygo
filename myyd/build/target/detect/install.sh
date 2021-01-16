#!/bin/bash
#@auth alexchen
#@time 20200603

SHELLPATH=$(cd `dirname $0`; pwd)
cd ${SHELLPATH}

ACL="/virus/cloud_waf_detect/acl"
LINK="/virus/cloud_waf_detect/link"
LINE="/virus/cloud_waf_detect/line"
ACLSERVER="/virus/cloud_waf_detect/aclserver"
LINKSERVER="/virus/cloud_waf_detect/linkserver"
LINESERVER="/virus/cloud_waf_detect/lineserver"
WEBSITE="/virus/cloud_waf_detect/saas_services/saasyd_website"

function install_acl_node() {
    if [ ! -d ${ACL} ]; then
        mkdir -p ${ACL}
    else
        rm -rf ${ACL}/acl_node
    fi
    chmod +x node/acl_node
    cp node/acl_node ${ACL}
    md5sum ${ACL}/acl_node
    echo -e "\033[32m === ACL install success === \033[0m"
}

function install_link_node() {
    if [ ! -d ${LINK} ]; then
        mkdir -p ${LINK}
    else
        rm -rf ${LINK}/link_node
    fi
    chmod +x node/link_node
    cp node/link_node ${LINK}
    md5sum ${LINK}/link_node
    echo -e "\033[32m === LINK install success === \033[0m"
}

function install_line_node() {
    if [ ! -d ${LINE} ]; then
        mkdir -p ${LINE}
    else
        rm -rf ${LINE}/line_node
    fi
    chmod +x node/line_node
    cp node/line_node ${LINE}
    md5sum ${LINE}/line_node
    echo -e "\033[32m === LINE install success === \033[0m"
}

function install_detect_node() {
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
}

function install_acl_server() {
    if [ ! -d ${ACLSERVER} ]; then
        mkdir -p ${ACLSERVER}
    else
        rm -rf ${ACLSERVER}/acl_server
    fi
    chmod +x server/acl_server
    cp server/acl_server ${ACLSERVER}
    md5sum ${ACLSERVER}/acl_server
    echo -e "\033[32m === ACLSERVER install success === \033[0m"
}

function install_link_server() {
    if [ ! -d ${LINKSERVER} ]; then
        mkdir -p ${LINKSERVER}
    else
        rm -rf ${LINKSERVER}/link_server
    fi
    chmod +x server/link_server
    cp server/link_server ${LINKSERVER}
    md5sum ${LINKSERVER}/link_server
    echo -e "\033[32m === LINKSERVER install success === \033[0m"
}

function install_line_server() {
    if [ ! -d ${LINESERVER} ]; then
        mkdir -p ${LINESERVER}
    else
        rm -rf ${LINESERVER}/line_server
    fi
    chmod +x server/line_server
    cp server/line_server ${LINESERVER}
    md5sum ${LINESERVER}/line_server
    echo -e "\033[32m === ACLSERVER install success === \033[0m"
}

function install_detect_server() {
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
}

function install_website_server() {
    if [ ! -d ${WEBSITE} ]; then
        mkdir -p ${WEBSITE}
    else
        rm -rf ${WEBSITE}/saasyd_website
    fi
    chmod +x server/saasyd_website
    cp server/saasyd_website ${WEBSITE}
    md5sum ${WEBSITE}/saasyd_website
    echo -e "\033[32m === WEBSITE install success === \033[0m"
}

function print_help {
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
    echo -e "\033[35m #./install.sh {param} \033[0m"
    echo -e "\033[35m {param}: \033[0m"
    echo -e "\033[35m       acl             : make acl node \033[0m"
    echo -e "\033[35m       link            : make link node \033[0m"
    echo -e "\033[35m       line            : make line node \033[0m"
    echo -e "\033[35m       detect_node     : make all node \033[0m"
    echo -e "\033[35m       aclserver       : make acl server \033[0m"
    echo -e "\033[35m       lineserver      : make link server \033[0m"
    echo -e "\033[35m       linkserver      : make line server \033[0m"
    echo -e "\033[35m       detect_server   : make all server \033[0m"
    echo -e "\033[35m       website         : make website server \033[0m"
    echo -e "\033[35m       all             : make all node/server/website \033[0m"
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
}

function main() {
    case $1 in
    "acl")
        install_detect_node acl
    ;;
    "link")
        install_detect_node link
    ;;
    "line")
        install_detect_node line
    ;;
    "detect_node")
        install_detect_node all
    ;;
    "aclserver")
        install_detect_server acl
    ;;
    "linkserver")
        install_detect_server link
    ;;
    "lineserver")
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