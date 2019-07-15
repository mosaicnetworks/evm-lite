#!/bin/bash

FROMNODE=$1

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

FROMADD=$(scripts/get_node_address.sh $FROMNODE)
PASSWD=$mydir/../../deploy/conf/eth/pwd.txt)

evmlc poa init --from $FROMADD --pwd $PASSWD
