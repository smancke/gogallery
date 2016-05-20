#!/bin/bash

dir=/go/src/github.com/smancke/gogallery

docker run --rm -v "$PWD":$dir -w $dir golang:1.6 bash \
       -c 'go get -v ./... &&  go build -o gogallery .'

exit $?
