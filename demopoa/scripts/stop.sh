#!/bin/bash

docker ps -f name=node -aq | xargs docker rm -f 
docker network rm babblenet