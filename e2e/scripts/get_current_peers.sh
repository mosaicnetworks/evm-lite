#!/bin/bash

PEERSOUT=$1
PORT=8080

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NODEIP=$(bash $mydir/get_running_nodes.sh | shuf -n 1 | awk '{print $2}')

wget -q -O $PEERSOUT $NODEIP:$PORT/peers