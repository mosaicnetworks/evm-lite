#!/bin/bash



mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

#TODO for the moment the pregenesis file is a fixed name and location. 
#     It could be parameterised later, but there is no pressing need.

PREGENESIS=$(readlink -f "$mydir/../conf/poa/pregenesis.json")
TMPOUT=$(readlink -f "$mydir/../conf/poa")
if [ ! -d "$TMPOUT" ] ; then
	mkdir "$TMPOUT"
fi

echo "Pregenesis file: $PREGENESIS"

if [ ! -f "$PREGENESIS" ] ; then
	echo "No Pregenesis file found. Aborting."
	exit 1
else
    echo "$PREGENESIS"
    ls -ld "$PREGENESIS"
fi


#   .option('-p, --precompfile <file>', 'pregenesis.json file')
#   .option('-o, --output-dir <directory>', 'output directory')

echo node $mydir/src/compilegenesis.js --precompfile $PREGENESIS  --output-dir $TMPOUT
node $mydir/src/compilegenesis.js --precompfile $PREGENESIS  --output-dir $TMPOUT


cp $TMPOUT/genesis.json $TMPOUT/../genesis.json

for i in $mydir/../conf/node*/eth
do
  echo "Writing $i"
  cp $mydir/../conf/genesis.json $i/genesis.json
done




