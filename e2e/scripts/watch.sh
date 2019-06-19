#!/bin/bash
set -eux
PORT=${1:-8000}

watch -t -n 1 '
docker ps --filter name=node --format "{{.Names}}" | sort -u | \
xargs docker inspect -f "{{.NetworkSettings.Networks.monet.IPAddress}}" | \
xargs -I % curl -s -m 1 http://%:'${PORT}'/stats | \
tr -d "{}\"" | \
awk -F "," '"'"'{gsub (/[,]/," "); print;}'"'"'
'