#!/bin/bash

# This script creates the configuration for a Babble testnet with a variable  
# number of nodes. It will generate crytographic key pairs and assemble a 
# peers.json file in the format used by Babble. The files are copied into 
# individual folders for each node which can be used as the datadir that Babble 
# reads configuration from. 

set -e

N=${1:-4}
IPBASE=${2:-node}
IPADD=${3:-0}
DEST=${4:-"$PWD/conf"}
PORT=${5:-1337}


l=$((N-1))

for i in $(seq 0 $l) 
do
	dest=$DEST/node$i/babble
	mkdir -p $dest
	echo "Generating key pair for node$i"
	docker run \
		-v $dest:/.babble \
		--rm mosaicnetworks/babble:0.4.0 keygen
	echo "$IPBASE$(($IPADD + $i)):$PORT" > $dest/addr
done

PFILE=$DEST/peers.json
echo "[" > $PFILE 
for i in $(seq 0 $l)
do
	dest=$DEST/node$i/babble
	
	com=","
	if [[ $i == $l ]]; then 
		com=""
	fi
	
	printf "\t{\n" >> $PFILE
	printf "\t\t\"NetAddr\":\"$(cat $dest/addr)\",\n" >> $PFILE
	printf "\t\t\"PubKeyHex\":\"$(cat $dest/key.pub)\"\n" >> $PFILE
	printf "\t}%s\n"  $com >> $PFILE

done
echo "]" >> $PFILE

for i in $(seq 0 $l) 
do
	dest=$DEST/node$i/babble
	cp $DEST/peers.json $dest/
done

