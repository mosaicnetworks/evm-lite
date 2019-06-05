#!/bin/bash

NOMINEENODE=$1

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

# Find nodes running babble
# The search parameters may become soft coded in a future revision
RUNNINGNODES=$(docker ps --format "{{.Names}}" | sort -u )
PASSWD=$(cat $mydir/../../deploy/conf/eth/pwd.txt)
NOMADD=$($mydir/get_sleeping_node_address.sh $NOMINEENODE)

echo "NOMADD: " $NOMADD

for i in $(echo $RUNNINGNODES)
do			
   $mydir/write_evmlc_config.sh $i
   FROMADD=$($mydir/get_sleeping_node_address.sh $i)
   evmlc poa vote --address $NOMADD --from $FROMADD --verdict true --pwd $PASSWD
done


