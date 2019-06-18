#!/bin/bash

# Return the "[name] [ip]" list of running evm-lite nodes

docker ps --filter name=node --format "{{.Names}}" | \
sort -u | \
xargs docker inspect -f "{{.Config.Hostname}} {{.NetworkSettings.Networks.monet.IPAddress}}"