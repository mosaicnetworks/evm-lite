#!/bin/bash

FROMNODE=$1

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

FROMADD=$(scripts/get_sleeping_node_address.sh $FROMNODE)
PASSWD=$(cat $mydir/../../deploy/conf/eth/pwd.txt)

evmlc poa init --from $FROMADD --pwd $PASSWD
