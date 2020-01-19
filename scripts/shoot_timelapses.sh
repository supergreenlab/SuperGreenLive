#!/bin/bash

PATH="$(pwd):$PATH"
TS="$(date +%s)"

function shoot_timelapse() {
  mkdir -p $1
  if [ "$4" = "r" ]; then
    ./download -u $1 -min $(date -d "-2day $2" +%s) -max $(date -d "-1day $3" +%s)
  else
    ./download -u $1 -min $(date -d "-1day $2" +%s) -max $(date -d "-1day $3" +%s)
  fi
  ./create_timelapse.sh $1
  cp $1.mp4 $1.$TS.mp4
  mv $1.mp4 /var/www/html/timelapses/
  upload $1.$TS.mp4
  rm $1.$TS.mp4
  rm -rf $1
}

#
# Uncomment line below and replace TIMELAPSE_UPLOAD_NAME
#

# shoot_timelapse TIMELAPSE_UPLOAD_NAME "1:00" "19:00"

cd /var/www/html/timelapses
./mosaic.sh
