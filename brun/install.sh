#!/usr/bin/env bash

EXTRACT_DIR="./.resk"
EXEC="./run.sh"

function main() {
    NUM=`awk '/^__SHELL_END__/ {print NR+1;exit 0}' $0`
    rm -fr ${EXTRACT_DIR} && mkdir ${EXTRACT_DIR}
    echo ${EXTRACT_DIR} ${NUM}
    tail -n+${NUM} $0 | tar -xvj -C ${EXTRACT_DIR}
    cd ${EXTRACT_DIR}
    chmod +x ${EXEC}
    ${EXEC} $*
    return 0
}
main $@
exit $?
# shell 结束标识
# 注意最后一行，必须有且只有一个换行，而且没有任何字符
__SHELL_END__
