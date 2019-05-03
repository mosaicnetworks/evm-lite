#!/bin/bash
set -eux
IPS=${1:-"../deploy/terraform/local/ips.dat"}
PORT=${2:-8000}

watch -n 1 '
cat '${IPS}'  | \
awk '"'"'{print $2}'"'"' | \
xargs -I % curl -s -m 1 http://%:'${PORT}'/stats |\
tr -d "{}\"" | \
awk -F "," '"'"'{gsub (/[,]/," "); print;}'"'"'
'