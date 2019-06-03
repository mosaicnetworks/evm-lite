#!/bin/bash

N=${1:-2}
DEST=${2:-"$PWD/conf"}

dest=$DEST/node$N

# get genesis.peers.json
echo "Fetching peers.genesis.json from node1"
curl -s http://172.77.5.1:8080/genesispeers > $dest/babble/peers.genesis.json

# get up-to-date peers.json
echo "Fetching peers.json from node1"
curl -s http://172.77.5.1:8080/peers > $dest/babble/peers.json

# start the new node
docker create \
    -u $(id -u) \
    --name=node$N \
    --net=babblenet \
    --ip=172.77.5.$N \
    mosaicnetworks/evm-lite:latest run babble

docker cp $dest node$N:"/.evm-lite"

docker start node$N

