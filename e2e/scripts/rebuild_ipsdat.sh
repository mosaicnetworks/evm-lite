#!/bin/bash

echo "Rebuilding ips.dat"
mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"


for node in $(docker ps --format "{{.Names}}" | sort -u)
do
  IP=$(bash $mydir/get_ip_of_container.sh  $node)	
  echo "$node $IP"
done > $mydir/../../deploy/terraform/local/ips.dat

echo "Rebuilt ips.dat"
