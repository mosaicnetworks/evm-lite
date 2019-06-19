#!/bin/bash

FROMNODE=$1
NOMINEENODE=$2
VERDICT=$3

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NOMADD=$($mydir/get_node_address.sh $NOMINEENODE)
FROMADD=$($mydir/get_node_address.sh $FROMNODE)
PASSWD=$(cat $mydir/../../deploy/conf/eth/pwd.txt)

evmlc poa vote --address $NOMADD --from $FROMADD --verdict $VERDICT --pwd $PASSWD
