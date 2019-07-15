#!/bin/bash

# randomly select a node from the list of running nodes, and query its /peers
# api endpoint

PEERSOUT=$1
PORT=8000 # babble service port

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NODEIP=$($mydir/get_running_nodes.sh | shuf -n 1 | awk '{print $2}')

wget -q -O $PEERSOUT $NODEIP:$PORT/peers