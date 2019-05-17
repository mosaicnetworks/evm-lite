#!/bin/bash

set -eu

NODENO=${1:-0}
IPS=${2:-"../deploy/terraform/local/ips.dat"}
KEY_DIR=${3:-"../deploy/conf/babble/conf/keystore/"}
PWD_FILE=${4:-"../deploy/conf/eth/pwd.txt"}
PORT=${5:-8080}
SOL_FILE=${6:-"smart-contracts/genesis.sol"}


# IP addresses are sorted so in node0, node1 etc order.
ips=($(sort ${IPS} | awk '{ print $2 }' | paste -sd "," -))

node src/genesiscontract.js --ips=$ips \
    --port=$PORT \
    --contract=$SOL_FILE \
    --keystore=$KEY_DIR \
    --pwd=$PWD_FILE \
    --nodeno=$NODENO 
