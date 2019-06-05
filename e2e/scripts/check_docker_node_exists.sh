#!/bin/bash



NODE=$1


STATUS=$(docker inspect -f '{{ .State.Status }}' $NODE 2> /dev/null)

if [ "$?" -gt 0 ] ; then
  >&2 echo "Node $NODE does not exist. Aborting."
  exit 1
fi

exit 0

