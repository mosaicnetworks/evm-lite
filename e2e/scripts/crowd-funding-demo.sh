#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

KEY_DIR=${1:-"$mydir/../../deploy/conf/babble/conf/keystore/"}
PWD_FILE=${2:-"$mydir/../../deploy/conf/eth/pwd.txt"}
PORT=${3:-8080}
SOL_FILE=${4:-"$mydir/../smart-contracts/CrowdFunding.sol"}

ips=($($mydir/get_running_nodes.sh | awk '{print $2}' | paste -sd "," -))

set -x
node crowd-funding/demo.js --ips=$ips \
    --port=$PORT \
    --contract=$SOL_FILE \
    --keystore=$KEY_DIR \
    --pwd=$PWD_FILE