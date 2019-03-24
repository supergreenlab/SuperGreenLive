#!/bin/bash

PATH="$(pwd):$PATH"
TS="$(date +%s)"
DIR="daily_timelapses/$TS"
mkdir -p "$DIR"
pushd "$DIR"

daily_timelapse.sh SuperGreenOffice
mv SuperGreenOffice.mp4 SuperGreenOffice.$TS.mp4
upload SuperGreenOffice.$TS.mp4

daily_timelapse.sh SuperGreenOfficeVeg
mv SuperGreenOfficeVeg.mp4 SuperGreenOfficeVeg.$TS.mp4
upload SuperGreenOfficeVeg.$TS.mp4

daily_timelapse.sh SuperGreenOfficeVeg2
mv SuperGreenOfficeVeg2.mp4 SuperGreenOfficeVeg2.$TS.mp4
upload SuperGreenOfficeVeg2.$TS.mp4

popd
