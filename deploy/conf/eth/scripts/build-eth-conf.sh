#!/bin/bash

# This script produces the Ethereum configuration for each node on the testnet. 
# It generates a new Ethereum key, controlled by the same password, for each 
# node, and aggretates the corresponding public keys in a genesis.json file, 
# which is used by evm-lite to initialize the accounts in the State. 
# Additionally, it creates the base for the evml.toml file used by evm-lite to 
# read configuration, on top of command-line flags. The configuration files are 
# placed in different folders (named node0...nodeN) which can be copied or 
# mounted directly in the root directory for evm-lite (controlled by 'datadir' 
# flag). The output of this script, executed with default parameters, will look 
# something like this:
#
#	conf/solo/conf
#	├── genesis.json
#	├── keystore
#	│   ├── node0-key.json
#	│   ├── node1-key.json
#	│   ├── node2-key.json
#	│   └── node3-key.json
#	├── node0
#	│   └── eth
#	│       ├── genesis.json
#	│       ├── keystore
#	│       │   └── key.json
#	│       └── pwd.txt
#	├── node1
#	│   └── eth
#	│       ├── genesis.json
#	│       ├── keystore
#	│       │   └── key.json
#	│       └── pwd.txt
#	├── node2
#	│   └── eth
#	│       ├── genesis.json
#	│       ├── keystore
#	│       │   └── key.json
#	│       └── pwd.txt
#	└── node3
#	    └── eth
#	        ├── genesis.json
#	        ├── keystore
#	        │   └── key.json
#	        └── pwd.txt




set -e

N=${1:-4} # number of nodes
DEST=${2:-"$(pwd)/../conf"} # output directory
PASS=${3:-"$(pwd)/../pwd.txt"} # password file for Ethereum accounts

l=$((N-1))

for i in $(seq 0 $l) 
do
	dest=$DEST/node$i/eth
	mkdir -p $dest
    # Use a Docker container to run the 'evml keys' command that creates 
	# accounts. This saves us the trouble of installing evml locally.
	# The file is written directly into the mounted directory.
    docker run --rm \
		-u $(id -u) \
		-v $dest:/datadir \
		-v $PASS:/pwd.txt \
		mosaicnetworks/evm-lite:latest keys --passfile=/pwd.txt generate /datadir/keystore/key.json  | \
    		awk '/Address/ {print $2}'  >> $dest/addr
done

# Generate the genesis file
GFILE=$DEST/genesis.json
echo "{" > $GFILE 
printf "\t\"alloc\": {\n" >> $GFILE
for i in $(seq 0 $l)
do
	com=","
	if [[ $i == $l ]]; then 
		com=""
	fi
	printf "\t\t\"$(cat $DEST/node$i/eth/addr)\": {\n" >> $GFILE
    printf "\t\t\t\"balance\": \"133700000000000000000$i\",\n" >> $GFILE
    printf "\t\t\t\"moniker\": \"node$i\"\n" >> $GFILE
    printf "\t\t}%s\n" $com >> $GFILE
done
printf "\t}\n" >> $GFILE
echo "}" >> $GFILE



# Generate the pregenesis file
GFILE=$DEST/pregenesis.json
echo "{" > $GFILE 

printf "\t\"precompiler\":{\n" >> $GFILE
printf "\t\t \"contracts\": [\n" >> $GFILE
printf "\t\t\t {\n" >> $GFILE
printf "\t\t\t\t \"address\": \"abbaabbaabbaabbaabbaabbaabbaabbaabbaabba\",\n" >> $GFILE
printf "\t\t\t\t \"filename\": \"genesis.sol\",\n" >> $GFILE
printf "\t\t\t\t \"authorising\": \"true\",\n" >> $GFILE
printf "\t\t\t\t \"contractname\": \"POA_Genesis\",\n" >> $GFILE
printf "\t\t\t\t \"balance\": \"1337000000000000000099\",\n" >> $GFILE
printf "\t\t\t\t \"preauthorised\": [\n" >> $GFILE


comma=""

if [ $l -lt 3 ] ; then
  preauthnum=0
else
  preauthnum=1
fi
for i in $(seq 0 $preauthnum)
do
   printf "\t\t\t\t\t $comma{ \"address\": \"$(cat $DEST/node$i/eth/addr)\", \"moniker\": \"User$i\"}\n" >> $GFILE
   comma=","
done


printf "\t\t\t\t ]\n" >> $GFILE
printf "\t\t\t }\n" >> $GFILE
printf "\t\t ]\n" >> $GFILE
printf "\t },\n" >> $GFILE



printf "\t\"alloc\": {\n" >> $GFILE
for i in $(seq 0 $l)
do
	com=","
	if [[ $i == $l ]]; then 
		com=""
	fi
	printf "\t\t\"$(cat $DEST/node$i/eth/addr)\": {\n" >> $GFILE
    printf "\t\t\t\"balance\": \"133700000000000000000$i\",\n" >> $GFILE
    printf "\t\t\t\"moniker\": \"node$i\"\n" >> $GFILE
    printf "\t\t}%s\n" $com >> $GFILE
done
printf "\t}\n" >> $GFILE
echo "}" >> $GFILE





gKeystore=$DEST/keystore
mkdir -p $gKeystore

# Copy files into each node's folder and cleanup
for i in $(seq 0 $l) 
do
	dest=$DEST/node$i
	cp $DEST/genesis.json $dest/eth
	cp $PASS $dest/eth
	cp -r $dest/eth/keystore/key.json $gKeystore/node$i-key.json
    rm $dest/eth/addr
done

