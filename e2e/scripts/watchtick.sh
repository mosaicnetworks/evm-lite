#!/bin/bash

PORT=${1:-8000}

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

# This will become a command line parameter in production
 


# Find nodes running babble
# The search parameters may become soft coded in a future revision
   for i in $(docker ps --format "{{.Names}}"| sort -u )
   do
     ip=$(docker inspect -f '{{ .NetworkSettings.Networks.monet.IPAddress }}' $i)   
     echo -n  "$i $ip "
     curl -s -m 1 http://$ip:${PORT}/stats | tr -d "{}\""   | \
     awk -F "," '{gsub (/[,]/," "); print;}'   	
   done
