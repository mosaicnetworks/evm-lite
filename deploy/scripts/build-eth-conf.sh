#!/bin/bash

# This script produces the Ethereum configuration for each node on the testnet. 
# It generates a new Ethereum key, controlled by the same password, for each 
# node, and aggretates the corresponding public keys in a genesis.json file, 
# which is used by evm-lite to initialize the accounts in the State. 
# Additionally, it creates the base for the config.toml file used by evm-lite to 
# read configuration, on top of command-line flags. The configuration files are 
# placed in different folders (named node0...nodeN) which can be copied or 
# mounted directly in the root directory for evm-lite (controlled by 'datadir' 
# flag). The output of this script, executed with default parameters, will look 
# something like this:
#
#	conf/
#	├── config.toml
#	├── genesis.json
#	├── keystore
#	│   ├── UTC--2018-09-15T13-58-07.652863115Z--664f52f5866d0bea946fcb5cec67f18b93b574c0
#	│   ├── UTC--2018-09-15T13-58-09.640035569Z--c1a67fac13e90b93f28fce79a34eea63a9cfebfc
#	│   ├── UTC--2018-09-15T13-58-11.693951535Z--5917d40005da07924796a396d2e522da49490afd
#	│   └── UTC--2018-09-15T13-58-13.749396512Z--98a6b9400d294bb5a787583affd72b28faf2273f
#	├── node0
#	│   ├── config.toml
#	│   └── eth
#	│       ├── genesis.json
#	│       ├── keystore
#	│       │   └── UTC--2018-09-15T13-58-07.652863115Z--664f52f5866d0bea946fcb5cec67f18b93b574c0
#	│       └── pwd.txt
#	├── node1
#	│   ├── config.toml
#	│   └── eth
#	│       ├── genesis.json
#	│       ├── keystore
#	│       │   └── UTC--2018-09-15T13-58-09.640035569Z--c1a67fac13e90b93f28fce79a34eea63a9cfebfc
#	│       └── pwd.txt
#	├── node2
#	│   ├── config.toml
#	│   └── eth
#	│       ├── genesis.json
#	│       ├── keystore
#	│       │   └── UTC--2018-09-15T13-58-11.693951535Z--5917d40005da07924796a396d2e522da49490afd
#	│       └── pwd.txt
#	└── node3
#	    ├── config.toml
#	    └── eth
#	        ├── genesis.json
#	        ├── keystore
#	        │   └── UTC--2018-09-15T13-58-13.749396512Z--98a6b9400d294bb5a787583affd72b28faf2273f
#	        └── pwd.txt



set -e

N=${1:-4} # number of nodes
DEST=${2:-"$(pwd)/conf"} # output directory
PASS=${3:-"$(pwd)/pwd.txt"} # password file for Ethereum accounts

l=$((N-1))

for i in $(seq 0 $l) 
do
	dest=$DEST/node$i/eth
	mkdir -p $dest
    # use a Docker container to run the geth command that creates accounts. This
	# saves us the trouble of installing geth locally
    docker run --rm \
		-u `id -u $USER` \
		-v $dest:/datadir \
		-v $PASS:/pwd.txt \
		ethereum/client-go -verbosity=1 --datadir=/datadir --password=/pwd.txt account new  | \
    		awk '{gsub("[{}]", "\""); print $2}'  >> $dest/addr
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
	printf "\t\t$(cat $DEST/node$i/eth/addr): {\n" >> $GFILE
    printf "\t\t\t\"balance\": \"1337000000000000000000\"\n" >> $GFILE
    printf "\t\t}%s\n" $com >> $GFILE
done
printf "\t}\n" >> $GFILE
echo "}" >> $GFILE

# Generate the evm-lite config file
CFILE=$DEST/config.toml
echo "[eth]" > $CFILE 
echo "db = \"/eth.db\"" >> $CFILE

gKeystore=$DEST/keystore
mkdir -p $gKeystore

# Copy files into each node's folder and cleanup
for i in $(seq 0 $l) 
do
	dest=$DEST/node$i
	cp $DEST/config.toml $dest/config.toml
	cp $DEST/genesis.json $dest/eth
	cp $PASS $dest/eth
	cp -r $dest/eth/keystore/* $gKeystore
    rm $dest/eth/addr
done

