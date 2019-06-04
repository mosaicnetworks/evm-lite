#!/bin/bash

FROMNODE=$1
NOMINEENODE=$2

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NOMADD=$(scripts/get_sleeping_node_address.sh $NOMINEENODE)
FROMADD=$(scripts/get_sleeping_node_address.sh $FROMNODE)
PASSWD=$(cat $mydir/../../deploy/conf/eth/pwd.txt)

evmlc poa nominate --nominee $NOMADD --from $FROMADD --moniker $NOMINEENODE --pwd $PASSWD