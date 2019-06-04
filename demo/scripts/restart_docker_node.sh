#!/bin/bash

for i in "$@"
do
  docker restart $i 
done
