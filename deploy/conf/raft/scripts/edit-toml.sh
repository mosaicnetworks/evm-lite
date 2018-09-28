#!/bin/bash

# This script adds Raft configuration to an evm-lite config.toml file. 

set -e

N=${1:-4}
IPBASE=${2:-node}
IPADD=${3:-0}
DEST=${4:-"conf"}

l=$((N-1))

PFILE=$DEST/config.toml
echo "[raft]" >> $PFILE  
for i in $(seq 0 $l) 
do
	dest=$DEST/node$i
	cp $DEST/config.toml $dest/config.toml
	echo "node_addr = \"$IPBASE$(($IPADD + $i)):1337\"" >> $dest/config.toml
	echo "server_id = \"node$i\"" >> $dest/config.toml
done