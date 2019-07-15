#!/bin/bash


BUILDNAME=${1:-"newtemplate"}
CONSENSUS=${2:-"babble"}
TERRAFORM=${3:-"local"}

# Huge overkill to change directory to the deploy directory

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"

cd $DIR/../..

if [ -d "prebuilt/templates/$BUILDNAME" ] ; then
	echo "prebuilt/templates/$BUILDNAME already exists. Exiting."
 	exit 1
fi

BASE="prebuilt/templates/$BUILDNAME"

mkdir "$BASE"
cd "$BASE"
mkdir conf terraform

cp -pr  "../../../conf/$CONSENSUS" conf/
cp -pr  "../../../conf/eth" conf/
cp -pr  "../../../terraform/$TERRAFORM" terraform/

rm -rf "terraform/$TERRAFORM/.terraform"
rm -f terraform/$TERRAFORM/*tfstate*
rm -f terraform/$TERRAFORM/*.tf
rm -f "terraform/$TERRAFORM/makefile"

find . -name .gitignore -exec rm -f {} \;
find . -name makefile -exec rm -f {} \;
find . -name scripts -type d -exec rm -rf {} \;


pwd



