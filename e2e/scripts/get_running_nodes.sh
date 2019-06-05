#!/bin/bash


mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

# This will become a command line parameter in production
 
IPSDAT="$mydir/../../deploy/terraform/local/ips.dat"


# Find nodes running babble
# The search parameters may become soft coded in a future revision
RUNNINGNODES=$(docker ps --format "{{.Names}}" | sort -u )

for i in $(echo $RUNNINGNODES)
do
  nodeline=$(grep $i "$IPSDAT")

  if [ -z "$nodeline" ] ; then 
	nodeline=$i
  fi 
  echo $nodeline
done
