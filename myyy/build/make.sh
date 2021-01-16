#!/bin/sh

WORKSPACEPATH=$(cd ../`dirname $0`; pwd)
cd ${WORKSPACEPATH}
NODEDIR=${WORKSPACEPATH}/src/availd
SERVERDIR=${WORKSPACEPATH}/src/availserverd
WEBSITEDIR=${WORKSPACEPATH}/src/saasydwebsite
TARGETPATH=${WORKSPACEPATH}/build/target
TARGETNODEPATH=${TARGETPATH}/detect/node
TARGETSERVERPATH=${TARGETPATH}/detect/server

function make {
    SRC_PATH=$1
    TARGET=$2
    DST_PATH=$3
    if [ ! ${TARGET} ]; then
        echo -e "\033[31m [make] xxx make ${TARGET} is nill! \033[0m"
        return 1
    fi
    if [ ! -d ${SRC_PATH} ]; then
        echo -e "\033[31m [make] xxx path ${SRC_PATH} is exits! \033[0m"
        return 1
    fi
    cd ${SRC_PATH}
    TARGET_FILE=${SRC_PATH}/${TARGET}
    if [ -f ${TARGET_FILE} ]; then
        rm -rf ${TARGET_FILE}
    fi
    echo -e "\033[33m === start make ${TARGET_FILE} ... \033[0m"
    go build -o ${TARGET}
    if [ $? != 0 ]; then
        echo -e "\033[31m [make_node] xxx make ${TARGET} is failed! \033[0m"
        return 1
    fi
    if [ ! -f ${TARGET_FILE} ]; then
        echo -e "\033[31m [make] xxx make file ${TARGET_FILE} is not exits! \033[0m"
        return 1
    fi
    copy_to_target ${TARGET_FILE} ${DST_PATH}
    if [ $? != 0 ]; then
        echo -e "\033[31m [make_node] xxx copy ${TARGET} is failed! \033[0m"
        return 1
    fi
    md5sum ${DST_PATH}/${TARGET}
    echo -e "\033[32m +++ [make] make and copy ${TARGET} is success! \033[0m"
    return 0
}

function copy_to_target {
    SRC_FILE=$1
    DST_PATH=$2
    if [ ! -d ${DST_PATH} ]; then
        mkdir -p ${DST_PATH}
    else
        if [ -f ${SRC_FILE} ]; then
            rm -rf SRC_FILE
        fi
    fi
    cp ${SRC_FILE} ${DST_PATH}
    return 0
}

function make_node {
     TARGET=$1
     cd ${NODEDIR}
     if [ ${TARGET} ]; then
         make ${NODEDIR} ${TARGET} ${TARGETNODEPATH}
         if [ $? != 0 ]; then
             echo -e "\033[31m [make_node] xxx make ${TARGET} is failed! \033[0m"
             return 1
         fi
         return 0
     else
         NODELIST=("eye_node")
         for target in ${NODELIST[@]}
         do
             make ${NODEDIR} ${target} ${TARGETNODEPATH}
             if [ $? != 0 ]; then
                 echo -e "\033[31m [make_node] xxx make ${TARGET} is failed! \033[0m"
                 return 1
             fi
         done
         return 0
     fi
 }

 function make_file {
     TYPE=$1
     TARGET=$2
     SRC_PATH=""
     DST_PATH=""
     TARGET_LIST=()
     if [ "node" == "${TYPE}" ]; then
        SRC_PATH=${NODEDIR}
        DST_PATH=${TARGETNODEPATH}
        TARGET_LIST=("eye_node")
     elif [ "server" == "${TYPE}" ] ; then
        SRC_PATH=${SERVERDIR}
        DST_PATH=${TARGETSERVERPATH}
        TARGET_LIST=("eye_node")
     elif [ "website" == "${TYPE}" ] ; then
        SRC_PATH=${WEBSITEDIR}
        DST_PATH=${TARGETSERVERPATH}
     else
        return 1
     fi

     cd ${SRC_PATH}
     if [ ${TARGET} ]; then
         make ${SRC_PATH} ${TARGET} ${DST_PATH}
         if [ $? != 0 ]; then
             echo -e "\033[31m [make_node] xxx make ${TARGET} is failed! \033[0m"
             return 1
         fi
         return 0
     else
         for target in ${TARGET_LIST[@]}
         do
             make ${SRC_PATH} ${target} ${DST_PATH}
             if [ $? != 0 ]; then
                 echo -e "\033[31m [make_node] xxx make ${TARGET} is failed! \033[0m"
                 return 1
             fi
         done
         return 0
     fi
 }

function print_help {
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
    echo -e "\033[35m #./make.sh {param} \033[0m"
    echo -e "\033[35m {param}: \033[0m"
    echo -e "\033[35m       eye_node        : make eye node \033[0m"
    echo -e "\033[35m       eye_server      : make eye server \033[0m"
    echo -e "\033[35m       all             : make all node/server/website \033[0m"
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
}

function main {
    echo -e "\033[34m ######################### 编译 ######################### \033[0m"
    case $1 in
        "eye_node")
            make_file node eye_node
        ;;
        "eye_server")
            make_file server eye_server
        ;;
        "all")
            make_file node eye_node
            make_file server eye_server
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