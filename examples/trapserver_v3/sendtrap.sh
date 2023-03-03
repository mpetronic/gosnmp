#!/bin/bash

LOOPS=${1:-1}
n=1

while true
do
    echo "Sending trap $n"
    sudo snmptrap -v3 -e 0x1122334455 -l authPriv -u testuser -a SHA -A testauth -x AES -X testpriv localhost:9162 "" .1.2.1 .1.2.2 s "HELLO_$n"
    [[ $n -ge $LOOPS ]] && exit 0
    n=$(($n + 1))
done
