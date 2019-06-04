#!/bin/bash

FROMNODE=$1
NOMINEENODE=$2
VERDICT=$3

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NOMADD=$(scripts/get_sleeping_node_address.sh $NOMINEENODE)
FROMADD=$(scripts/get_sleeping_node_address.sh $FROMNODE)
PASSWD=$(cat $mydir/../conf/eth/pwd.txt)

evmlc poa vote --address $NOMADD --from $FROMADD --verdict $VERDICT --pwd $PASSWD
