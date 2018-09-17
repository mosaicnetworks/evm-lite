#!/bin/bash

# This script creates the configuration for a Raft testnet with a variable  
# number of nodes. It will assemble a peers.json file in the format used by 
# Raft.  

set -e

N=${1:-4}
DEST=${2:-"conf"}
PORT=${4:-1337}

l=$((N-1))

mkdir -p $DEST

PFILE=$DEST/peers.json
echo "[" > $PFILE 
for i in $(seq 0 $l)
do
	com=","
	if [[ $i == $l ]]; then 
		com=""
	fi
	
	printf "\t{\n" >> $PFILE
	printf "\t\t\"address\":\"node$i:$PORT\",\n" >> $PFILE
	printf "\t\t\"id\":\"node$i\"\n" >> $PFILE
	printf "\t}%s\n"  $com >> $PFILE

done
echo "]" >> $PFILE

for i in $(seq 0 $l) 
do
	dest=$DEST/node$i/raft
	mkdir -p $dest
    cp $DEST/peers.json $dest/
done

