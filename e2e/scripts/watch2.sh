#!/bin/bash

PORT=${1:-8000}

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

# This will become a command line parameter in production
 


# Find nodes running babble
# The search parameters may become soft coded in a future revision
watch -n 1 $mydir/watchtick.sh
