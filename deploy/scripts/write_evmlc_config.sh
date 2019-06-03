#!/bin/bash


NODEID=$1
# NODEID=node1

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"



HOSTIP=$(bash $mydir/get_ip_of_container.sh $NODEID  )
DEFAULTADDRESS=$(bash $mydir/get_node_address.sh $NODEID | sed -e's/"address"://g;s/"//g')
KEYSTORE=$(readlink -f $mydir/../conf/babble/conf/$NODEID/eth/keystore)

{
echo "[connection]"
echo "host = \"$HOSTIP\""
echo "port = 8080.0"
echo ""
echo "[defaults]"
echo "from = \"$DEFAULTADDRESS\""
echo "gas = 10000000.0"
echo "gasPrice = 0.0"
echo ""
echo "[storage]"
echo "keystore = \"$KEYSTORE\""
} > ~/.evmlc/config.toml
