#!/usr/bin/env bash

#cpath=`pwd`
#PROJECT_PATH=${cpath%src*}
#echo $PROJECT_PATH
#export GOPATH=$GOPATH:${PROJECT_PATH}

SOURCE_FILE_NAME=main
TARGET_FILE_NAME=reskd

rm -fr ${TARGET_FILE_NAME}*

build(){
    echo $GOOS $GOARCH
    tname=${TARGET_FILE_NAME}_${GOOS}_${GOARCH}${EXT}
    # 注意这里必须将变量的值设置到环境变量中，但是不需要设置为全局的环境变量
    env GOOS=$GOOS GOARCH=$GOARCH \
    go build -o ${tname} -v ${SOURCE_FILE_NAME}.go
    chmod +x ${tname}
    
    # 打包
     mv ${tname} ${TARGET_FILE_NAME}${EXT}
     if [[ ${GOOS} == "windows" ]];then
        zip ${tname}.zip ${TARGET_FILE_NAME}${EXT} boot.ini        
     else
        tar --exclude=*.bz2 --exclude=*.zip --exclude=*.git -cvjf ${tname}.tar.bz2  ${TARGET_FILE_NAME}${EXT} boot.ini run.sh  -C ./
        cp install.sh ${tname}.run
        cat ${tname}.tar.bz2 >> ${tname}.run
        chmod +x ${tname}.run
        rm ${tname}.tar.bz2
     fi 
     mv ${TARGET_FILE_NAME}${EXT} ${tname}    
}

CGO_ENABLED=0

#查看go支持的平台和cpu架构
#go tool dist list

echo "-----------mac os 64------------"
#mac os 64
GOOS=darwin
GOARCH=amd64
EXT=
#build

echo "-----------linux 64------------"
#linux 64
GOOS=linux
GOARCH=amd64
EXT=
build

echo "-----------windows 64------------"
#windows
#64
GOOS=windows
GOARCH=amd64
EXT='.exe'
#build

echo "-----------windows 32------------"
#32
GOARCH=386
#build
