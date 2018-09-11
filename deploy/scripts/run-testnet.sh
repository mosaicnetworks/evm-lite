#!/bin/bash

set -eux

N=${1:-4}
MPWD=$(pwd)

docker network create \
  --driver=bridge \
  --subnet=172.77.0.0/16 \
  --ip-range=172.77.5.0/24 \
  --gateway=172.77.5.254 \
  monet

for i in $(seq 1 $N)
do
    docker create --name=node$i --net=monet --ip=172.77.5.$i mosaicnetworks/evm-lite:0.1.0 solo
    docker cp conf/node$i node$i:/.evm-lite
    docker start node$i
done