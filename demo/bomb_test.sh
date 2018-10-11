#!/bin/bash

set -u

FROM=$1
TO=${2:-"0xd95a4329822f2baa4cfbbb5e155b8a5bd92b3e39"}
COUNT=${3:-10}
NODE=${4:-"172.77.5.1"}


for i in `seq 1 $COUNT`; do
    curl -X POST http://$NODE:8080/tx -d '{"from":"'$FROM'","to":"'$TO'","value":6666}'
done; 