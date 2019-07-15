#!/bin/bash

for i in "$@"
do
  docker exec $i kill 1
done
