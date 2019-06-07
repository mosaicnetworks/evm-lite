#!/bin/bash



mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

#TODO for the moment the pregenesis file is a fixed name and location. 
#     It could be parameterised later, but there is no pressing need.

PREGENESIS=$(readlink -f "$mydir/../conf/pregenesis.json")
TMPOUT=$(readlink -f "$mydir/../conf")
if [ ! -d "$TMPOUT" ] ; then
	mkdir "$TMPOUT"
fi

# echo $PREGENESIS

if [ ! -f "$PREGENESIS" ] ; then
	>&2 echo "No Pregenesis file found. Aborting."
  echo "$PREGENESIS"
	exit 1
else
    echo "$PREGENESIS"
    ls -ld "$PREGENESIS"
fi


#   .option('-p, --precompfile <file>', 'pregenesis.json file')
#   .option('-o, --output-dir <directory>', 'output directory')

echo node $mydir/src/compilegenesis.js --precompfile $PREGENESIS  --output-dir $TMPOUT
node $mydir/src/compilegenesis.js --precompfile $PREGENESIS  --output-dir $TMPOUT

for i in $mydir/../conf/node*/eth
do
  echo "Writing $i"
  cp $mydir/../conf/genesis.json $i/genesis.json
done




