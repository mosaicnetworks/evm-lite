#!/bin/bash

DOCKER=$(docker ps --all)

DOCKERLINECNT=$(echo "$DOCKER" | wc -l )


if [[ "$DOCKERLINECNT" -lt 2 ]] ; then
	echo "No Docker Containers found."
	exit 0
fi

echo "$DOCKER"

>&2 echo "Docker containers found. Aborting."
exit 1

