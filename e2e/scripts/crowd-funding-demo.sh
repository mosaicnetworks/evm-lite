#!/bin/bash

set -eu

KEY_DIR=${1:-"../deploy/conf/babble/conf/keystore/"}
PWD_FILE=${2:-"../deploy/conf/eth/pwd.txt"}
PORT=${3:-8080}
SOL_FILE=${4:-"smart-contracts/CrowdFunding.sol"}

ips=($(bash scripts/get_running_nodes.sh | grep 'node0\|node1\|node2\|node3' | awk '{print $2}' | paste -sd "," -))

set -x
node crowd-funding/demo.js --ips=$ips \
    --port=$PORT \
    --contract=$SOL_FILE \
    --keystore=$KEY_DIR \
    --pwd=$PWD_FILE