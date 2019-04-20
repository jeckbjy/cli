#!/usr/bin/env bash
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}")" && pwd )
export GOPATH=$DIR/../../../../../
echo $GOPATH

APP=$1

if [[ $# -lt 1 ]]; then
    echo 'usage: build.sh $app'
    exit
fi

if [[ ! -d $DIR/$APP ]]; then
    echo 'cannot find'
    exit
fi

go build -o $DIR/output/$APP $DIR/$APP/main.go