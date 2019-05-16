#!/bin/bash

set -e

PREBUILT=${1:-" "}



if [ ! -d "templates/$PREBUILT" ] ; then
	echo "Could not find prebuilt configuration $PREBUILT."
	echo "Aborting."

	echo "Current prebuilt configurations are:"
	echo $( for i in templates/* 
		do
			echo $(basename $i)
		done
		)		

	exit 1
fi


cp -p -r templates/$PREBUILT/* ../

if [ -f "templates/$PREBUILT/.message" ] ; then 
	 echo ""
	 cat "templates/$PREBUILT/.message"
	 echo ""
fi

