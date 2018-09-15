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
#   conf
#   ├── config.toml
#   ├── genesis.json
#   ├── node0
#   │   ├── config.toml
#   │   └── eth
#   │       ├── genesis.json
#   │       ├── keystore
#   │       │   └── UTC--2018-09-15T10-32-13.572377320Z--14d0b6ede43e996b89905899c77918e2eabad992
#   │       └── pwd.txt
#   ├── node1
#   │   ├── config.toml
#   │   └── eth
#   │       ├── genesis.json
#   │       ├── keystore
#   │       │   └── UTC--2018-09-15T10-32-15.684974010Z--bc6fdf999a9d857a9d4ca690efcf889c6ed11d77
#   │       └── pwd.txt
#   ├── node2
#   │   ├── config.toml
#   │   └── eth
#   │       ├── genesis.json
#   │       ├── keystore
#   │       │   └── UTC--2018-09-15T10-32-17.693146314Z--9b2e13f23730621c10549817c0bda63c4ade822a
#   │       └── pwd.txt
#   └── node3
#       ├── config.toml
#       └── eth
#           ├── genesis.json
#           ├── keystore
#           │   └── UTC--2018-09-15T10-32-19.693701497Z--6c81fa82dfa54e6cf8073bfe746a0826dc982050
#           └── pwd.txt

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

# Copy files into each node's folder and cleanup
for i in $(seq 0 $l) 
do
	dest=$DEST/node$i
	cp $DEST/config.toml $dest/config.toml
	cp $DEST/genesis.json $dest/eth
	cp $PASS $dest/eth
    rm $dest/eth/addr
done

