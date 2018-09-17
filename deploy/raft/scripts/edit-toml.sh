#!/bin/bash

# This script adds Raft configuration to an evm-lite config.toml file. 

set -e

N=${1:-4}
DEST=${2:-"conf"}

l=$((N-1))

PFILE=$DEST/config.toml
echo "[raft]" >> $PFILE 
echo "snapshot_dir = \"/raft.snapshots\"" >> $PFILE
    
for i in $(seq 0 $l) 
do
	dest=$DEST/node$i
	cp $DEST/config.toml $dest/config.toml
	echo "node_addr = \"node$i:1337\"" >> $dest/config.toml
	echo "server_id = \"node$i\"" >> $dest/config.toml
done