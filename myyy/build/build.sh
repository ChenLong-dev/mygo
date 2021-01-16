#!/bin/sh

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
TARGETDIR=${BUILDDIR}/target
DETECTDIR=${TARGETDIR}/detect
MD5FILE=${DETECTDIR}/md5sum.info
SHOWDIR_SH=${BUILDDIR}/showdir.sh

# 编译源码文件
function make_src_code {
    dos2unix make.sh
    ./make.sh all
    if [ $? != 0 ]; then
        echo -e "\033[31m [make_src_code] xxx make all is failed! \033[0m"
        return 1
    fi
     echo -e "\033[32m [make_src_code] make src code is success! \033[0m"
     return 0
}

# 生成MD5值文件
function build_md5sum {
    echo -e "\033[34m ######################### MD5 ######################### \033[0m"
    cd ${DETECTDIR}
    if [ ! -f ${MD5FILE} ]; then
        touch ${MD5FILE}
    else
        > ${MD5FILE}
    fi

    md5sum ${DETECTDIR}/node/* ${DETECTDIR}/server/* > ${MD5FILE}
    if [ $? != 0 ]; then
        echo -e "\033[31m [build_md5sum] xxx md5sum is failed! \033[0m"
        return 1
    fi
    echo -e "\033[32m [build_md5sum] md5sum is success! \033[0m"
    return 0
}

# 生成压缩包
function tar_package {
    cd
}

# 打成压缩包
function build_package {
    echo -e "\033[34m ######################### 打包 ######################### \033[0m"
    cd ${BUILDDIR}
    if [ ! -d ${PACKAGEDIR} ]; then
        mkdir -p ${PACKAGEDIR}
    fi

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

function print_help {
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
    echo -e "\033[35m #./build.sh {param} \033[0m"
    echo -e "\033[35m {param}: \033[0m"
    echo -e "\033[35m       make           : make src node \033[0m"
    echo -e "\033[35m       build          : md5sum and build package \033[0m"
    echo -e "\033[35m       all            : make src code and build package  \033[0m"
    echo -e "\033[35m ######################### 帮助 ######################### \033[0m"
}

function main {
    case $1 in
        "make")
            make_src_code
        ;;
        "build")
            build_md5sum
            build_package
        ;;
        "all")
            make_src_code
            build_md5sum
            build_package
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



