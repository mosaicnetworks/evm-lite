#!/bin/bash


docker exec $1 cut -c2-53 /home/1000/.evm-lite/eth/keystore/key.json


