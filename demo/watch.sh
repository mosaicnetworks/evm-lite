#!/bin/bash
set -eux
IPS=${1:-"ips.dat"}
PORT=${2:-8000}

watch -n 1 '
cat '${IPS}'  | \
awk '"'"'{print $2}'"'"' | \
xargs -I % curl -s http://%:'${PORT}'/stats |\
tr -d "{}\"" | \
awk -F "," '"'"'{gsub (/[,]/," "); print;}'"'"'
'