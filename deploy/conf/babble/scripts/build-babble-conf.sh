#!/bin/bash

# This script creates the configuration for a Babble testnet with a variable  
# number of nodes.  
#
# The Babble validator keys are the same as the prefunded genesis-account keys.
# This script reuses the keys in conf/eth which were generated in a previous 
# step by build-eth-conf.sh. 
#
# This script also generates a peers.json file in the format used by Babble. The 
# files are copied into individual folders which can be used as the datadir for 
# each Babble instance.  

set -e

N=${1:-4}
IPBASE=${2:-node}
IPADD=${3:-0}
DEST=${4:-"$PWD/conf"}
PORT=${5:-1337}


l=$((N-1))

# use 'evml keys inspect' to extract the raw keys from the Ethereum-style 
# keystore files.
for i in $(seq 0 $l) 
do
	babble_dest=$DEST/node$i/babble
	eth_source=$DEST/node$i/eth

	mkdir -p $babble_dest
	echo "Generating key pair for node$i"
	
	docker run \
		-u $(id -u) \
		-v $eth_source:/.evm-lite \
		--rm \
		mosaicnetworks/evm-lite:latest keys --json --passfile /.evm-lite/pwd.txt inspect --private /.evm-lite/keystore/key.json  \
		> $babble_dest/key_info
	
	awk '/PublicKey/ { gsub("[\",]", ""); print $2 }' $babble_dest/key_info >> $babble_dest/key.pub
	awk '/PrivateKey/ { gsub("[\",]", ""); print $2 }' $babble_dest/key_info >> $babble_dest/priv_key
	
	echo "$IPBASE$(($IPADD + $i)):$PORT" >> $babble_dest/addr
done

# create the peers.json file
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
	printf "\t\t\"PubKeyHex\":\"0x$(cat $dest/key.pub)\"\n" >> $PFILE
	printf "\t}%s\n"  $com >> $PFILE

done
echo "]" >> $PFILE

for i in $(seq 0 $l) 
do
	dest=$DEST/node$i/babble
	cp $DEST/peers.json $dest/
	cp $DEST/peers.json $dest/peers.genesis.json
done

