#!/bin/bash


usage()
{
    echo "usage: compile_solidity.sh  [ --output-dir dir --contract contract | --help ]"
}



while [ "$1" != "" ]; do
    case $1 in
        -o | --output-dir )           shift
                                OUTPUTDIR=$1
                                ;;
        -c | --contract )        shift
                                CONTRACT=$1
                                ;;
        -h | --help )           usage
                                exit
                                ;;
        * )                     usage
                                exit 1
    esac
    shift
done



if [ -z "$OUTPUTDIR" ] ; then
	>&2 echo "You must specify the output directory with the --output-dir parameter. Aborting."
	exit 1	
fi

if [ -z "$CONTRACT" ] ; then
	>&2 echo "You must specify the contract file name with the --contract parameter. Aborting."
	exit 1	
fi


echo $*

command -v solc > /dev/null 2>&1 ||  { echo >&2 solc is not installed and is required. Aborting. ; exit 1; }

SOLCVERSION=$(solc --version)

echo Compiling using $SOLCVERSION


FILESTUB=$(basename $CONTRACT .sol)

ls -ls $CONTRACT

echo    solc --bin-runtime  --evm-version constantinople --overwrite --optimize --optimize-runs=100000 --abi -o $OUTPUTDIR $CONTRACT
   solc --bin-runtime  --evm-version constantinople --overwrite --optimize --optimize-runs=100000 --abi -o $OUTPUTDIR $CONTRACT


echo '{"version": "'$SOLCVERSION'"}' > ${OUTPUTDIR%/}/$FILESTUB.out 


