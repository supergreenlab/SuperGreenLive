#!/bin/bash

PATH="$(pwd):$PATH"
TS="$(date +%s)"

function shoot_timelapse() {
  rm -rf $1
  mkdir -p $1
  ./download -u $1 -min $(date -d "-1day $2" +%s) -max $(date -d "-1day $3" +%s)
  ./create_timelapse.sh $1
  cp $1.mp4 $1.$TS.mp4
  mv $1.mp4 /var/www/html/timelapses/
  upload $1.$TS.mp4
  rm $1.$TS.mp4
}

shoot_timelapse SuperGreenOffice "7:00" "19:00"
shoot_timelapse SuperGreenOfficeVeg "5:00" "19:00"
shoot_timelapse SuperGreenOfficeVeg2 "5:00" "19:00"
shoot_timelapse SuperGreenHouse "4:30" "22:30"
shoot_timelapse SuperGreenCardboard "6:00" "20:00"
