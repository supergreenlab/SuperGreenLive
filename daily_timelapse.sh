#!/bin/bash

set -e

#rm -rf $1
mkdir -p $1

./download $1 150
./create_timelapse.sh $1
