#!/bin/bash

# This script adds Babble configuration to an evm-lite evml.toml file. 

set -e

N=${1:-4}
IPBASE=${2:-node}
IPADD=${3:-0}
DEST=${4:-"conf"}

l=$((N-1))

PFILE=$DEST/evml.toml
echo "[babble]" >> $PFILE 
echo "store = false" >> $PFILE
echo "heartbeat = \"50ms\"" >> $PFILE
echo "timeout = \"200ms\"" >> $PFILE
    
for i in $(seq 0 $l) 
do
	dest=$DEST/node$i
	cp $DEST/evml.toml $dest/evml.toml
	echo "listen = \"$IPBASE$(($IPADD +$i)):1337\"" >> $dest/evml.toml
done