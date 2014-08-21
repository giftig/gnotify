#!/bin/bash

export GOPATH=`pwd`

# Install requirements from requirements.txt
for requirement in $(cat $(dirname $0)/../requirements.txt)
do
  go get $requirement
done

