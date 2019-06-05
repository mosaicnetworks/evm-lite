#!/bin/bash


# 	This does assume that the docker user ID is 1000. On Linux this is 
#	a fair assumption. Long term that user should become configurable.


docker exec $1 cut -c2-53 /home/1000/.evm-lite/eth/keystore/key.json


