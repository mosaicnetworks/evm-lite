#!/bin/bash



usage()
{
    echo "usage: nominate_node.sh  [ --from NominatingNode --nominee NominatedNode | --help ]"
}



while [ "$1" != "" ]; do
    case $1 in
        -f | --from )           shift
                                FROMNODE=$1
                                ;;
        -n | --nominee )        shift
                                NOMINEENODE=$1
                                ;;
        -h | --help )           usage
                                exit
                                ;;
        * )                     usage
                                exit 1
    esac
    shift
done


if [ -z "$FROMNODE" ] ; then
	>&2 echo "You must specify the nominating node with the --from parameter. Aborting."
	exit 1	
fi

if [ -z "$NOMINEENODE" ] ; then
	>&2 echo "You must specify the nominated node with the --nominee parameter. Aborting."
	exit 1	
fi


mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"


$mydir/check_node_running.sh "$FROMNODE"
if  [ $? -gt 0 ] ; then
	>&2 echo "Aborting."
        exit 1
fi

$mydir/check_node_exists.sh "$NOMINEENODE"
if  [ $? -gt 0 ] ; then
	>&2 echo "Aborting."
        exit 1
fi

# TODO check that $NOMINEENODE is exited and that the node has not been 
# nominated already and that the $FROMNODE is authorised. 

NOMADD=$($mydir/get_node_address.sh $NOMINEENODE)
FROMADD=$($mydir/get_node_address.sh $FROMNODE)
PASSWD=$(cat $mydir/../../deploy/conf/eth/pwd.txt)

evmlc poa nominate --nominee 0x$NOMADD --from 0x$FROMADD --moniker $NOMINEENODE --pwd $PASSWD
