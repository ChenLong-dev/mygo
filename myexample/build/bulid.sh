#!/bin/sh
#@auth cl
#@time 20201214

TARGET=myexample

PACKAGENAME="target"
VERSION=$2
if [ "$VERSION" == "" ]; then
	VERSION=`date +%s`
fi

#切换到脚本自己的目录
SHELLPATH=$(cd `dirname $0`; pwd)
cd ${SHELLPATH}
BUILDDIR=${SHELLPATH}
PACKAGEDIR=${BUILDDIR}/package
TARGETDIR=${BUILDDIR}/target/${TARGET}
MD5FILE=${BUILDDIR}/target/md5sum.info

# 编译源码文件
function make_src_code {
    dos2unix make.sh
    sh make.sh
    if [ $? != 0 ]; then
        echo -e "\033[31m [make_src_code] make all is failed! \033[0m"
        return 1
    fi
     echo -e "\033[32m [make_src_code] make src code is success! \033[0m"
     return 0
}

# 生成MD5值文件
function build_md5sum {
    echo -e "\033[34m ######################### MD5 ######################### \033[0m"
    cd ${TARGETDIR}
    if [ ! -f ${MD5FILE} ]; then
        touch ${MD5FILE}
    else
        > ${MD5FILE}
    fi

    md5sum ${TARGETDIR}/${TARGET} ${TARGETDIR}/conf/* > ${MD5FILE}
    if [ $? != 0 ]; then
        echo -e "\033[31m [build_md5sum] xxx md5sum is failed! \033[0m"
        return 1
    fi
    echo -e "\033[32m [build_md5sum] md5sum is success! \033[0m"
    return 0
}

# 打成压缩包
function build_package {
    echo -e "\033[34m ######################### 打包 ######################### \033[0m"
    cd ${BUILDDIR}
    if [ ! -d ${PACKAGEDIR} ]; then
        mkdir -p ${PACKAGEDIR}
    fi
    rm -rf ${PACKAGEDIR}/*
    tar --exclude=.git -czf /tmp/${PACKAGENAME}_${VERSION}.tar.gz ./target
    if [ $? != 0 ]; then
        echo -e "\033[31m [build_package] xxx1 build package is failed! \033[0m"
        return 1
    fi
    cp /tmp/${PACKAGENAME}_${VERSION}.tar.gz ${PACKAGEDIR}/
    md5sum ${PACKAGEDIR}/${PACKAGENAME}_${VERSION}.tar.gz
    if [ $? != 0 ]; then
        echo -e "\033[31m [build_package] xxx2 build package is failed! \033[0m"
        return 1
    fi
    echo -e "\033[32m [build_package] build package is success! \033[0m"
    rm -rf /tmp/${PACKAGENAME}_${VERSION}.tar.gz
    return 0
}

function main {
    make_src_code
    build_md5sum
    build_package
}

main $@
