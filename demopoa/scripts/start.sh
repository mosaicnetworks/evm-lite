#!/bin/bash

# This script creates the docker network where all the nodes will be hosted and
# starts node1 

set -eux

MPWD=$(pwd)

docker network create \
  --driver=bridge \
  --subnet=172.77.0.0/16 \
  --ip-range=172.77.0.0/16 \
  --gateway=172.77.5.254 \
  babblenet


docker create \
    -u $(id -u) \
    --name=node1 \
    --net=babblenet \
    --ip=172.77.5.1 \
    mosaicnetworks/evm-lite:latest run babble

docker cp $MPWD/conf/node1 node1:"/.evm-lite"

docker start node1

