#!/bin/bash

NODE=$1

STATUS=$(docker inspect -f '{{ .State.Status }}' $NODE)

if [ "$STATUS" != "running" ] ; then
  >&2 echo "Node $NODE is $STATUS. Aborting."
  exit 1
fi

exit 0
