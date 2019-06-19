#!/bin/bash

docker inspect -f '{{ .NetworkSettings.Networks.monet.IPAddress }}' $1


