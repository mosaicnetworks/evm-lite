#!/bin/bash

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"


for node in $(docker ps | cut -c134- | grep -v NAMES | sort -u)
do
  IP=$(bash $mydir/get_ip_of_container.sh  $node)	
  echo "$node $IP"
done > $mydir/../terraform/local/ips.dat

