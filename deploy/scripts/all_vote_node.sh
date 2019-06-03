#!/bin/bash

NOMINEENODE=$1

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"


# Find nodes running babble
# The search parameters may become soft coded in a future revision
RUNNINGNODES=$(docker ps | grep "run babble" | cut -c134- | sort -u )
PASSWD=$(cat $mydir/../conf/eth/pwd.txt)
NOMADD=$($mydir/get_sleeping_node_address.sh $NOMINEENODE)

for i in $(echo $RUNNINGNODES)
do			
   $mydir/write_evmlc_config.sh $i
   FROMADD=$($mydir/get_sleeping_node_address.sh $i)
   evmlc poa vote --address $NOMADD --from $FROMADD --verdict true --pwd $PASSWD

done


