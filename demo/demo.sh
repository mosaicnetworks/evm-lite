#!/bin/bash

set -eu

PORT=${2:-8080}
SOL_FILE=${3:-"nodejs/crowd-funding.sol"}
KEY_DIR=${4:-"conf/keystore"}
PWD_FILE=${5:-"conf/pwd.txt"}

ips="localhost,localhost,localhost,localhost"

node nodejs/demo.js --ips=$ips \
    --port=$PORT \
    --contract=$SOL_FILE \
    --keystore=$KEY_DIR \
    --pwd=$PWD_FILE
