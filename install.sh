#!/bin/bash

export GOPATH=`pwd`

# Install requirements from requirements.txt
for requirement in `cat requirements.txt`
do
  go get $requirement
done

