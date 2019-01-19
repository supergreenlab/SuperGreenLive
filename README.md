![SuperGreenLab](assets/sgl.png?raw=true "SuperGreenLab")

[![SuperGreenLab](assets/reddit-button.png?raw=true "SuperGreenLab")](https://www.reddit.com/r/SuperGreenLab)

Table of Contents
=================

   * [SuperGreenTimelapse](#supergreentimelapse)
   * [Features](#features)
      * [TODO](#todo)
   * [Hardware requirements](#hardware-requirements)
   * [Quickstart](#quickstart)
      * [Installation](#installation)
         * [Dropbox setup](#dropbox-setup)
         * [RaspberryPI setup](#raspberrypi-setup)
         * [Watermark](#watermark)
      * [Manage timelapses](#manage-timelapses)
      * [Create timelapse](#create-timelapse)
   * [View live over http](#view-live-over-http)
   * [One more thing:](#one-more-thing)

# SuperGreenTimelapse

A bunch of scripts/programs to produce dropbox-backed timelapses for raspberryPi.

![Example](assets/example.jpg?raw=true "Example")

# Features

- Live(-ish) secured webcam
- Take a picture every X minutes
- Upload to dropbox app folder
- Produce a timelapse video with extra image interpolation for better smoothness

## TODO

- Integrate https://github.com/gographics/imagick
- SSL ?

# Hardware requirements

- [RaspberryPI](https://www.raspberrypi.org/products/) + [Wifi (optional, most rpi have integrated wifi now)](https://www.raspberrypi.org/products/raspberry-pi-usb-wifi-dongle/)
- [Camera](https://www.raspberrypi.org/products/camera-module-v2/), I got [those](https://www.amazon.com/SainSmart-Fish-Eye-Camera-Raspberry-Arduino/dp/B00N1YJKFS) for the wide angle lens, but that's only for small spaces (this is the one used for the pic above).
- [Power supply](https://www.raspberrypi.org/products/raspberry-pi-universal-power-supply/)

# Quickstart

## Installation

### Dropbox setup

The problem that arises when you want to take timelapses is that taking pictures every 10 minutes takes a lot of space.
And having a raspberrypi running 24/24 and storing big amounts of data on an SD card is looking for trouble:P

So what seems to be a good solution is to upload everything to dropbox.

It also allows to view the latest pic online, which actually makes it some sort of cloud live camcorder. Good times.

There's a little setup to do on dropbox's side. For obvious security purpose you have to let dropbox know that he needs to create a space and access for our program.

Got to the [app creation page](https://www.dropbox.com/developers/apps/create), and choose: `Dropbox API` -> `App folder` -> `SuperGreenTimelapse`.

Now scroll to the `Generated access token` section, and click the `Generate` button below. Copy-paste the long id that looks like `vrB4PlxSQpsAAAAAAAC1SvJJbXi08sdjlkaWWfalk25iX4GAqsfk67rkM0sM0uyC`, we'll need that in the next step.

### RaspberryPI setup

First follow the [raspberryPI quickstart guide](https://www.raspberrypi.org/learning/software-guide/quickstart/) if you have never done that before.

Always good to upgrade after a fresh install:
```sh

sudo apt-get update && sudo apt-get upgrade -y

```

Install or upgrade to the latest binary with the following command:

```sh

sudo curl https://github.com/supergreenlab/SuperGreenTimelapse/releases/download/PreRelease/timelapse -o /usr/local/bin/timelapse
sudo curl https://github.com/supergreenlab/SuperGreenTimelapse/releases/download/PreRelease/watermark-logo.png -o /home/pi/watermark-logo.png
sudo chmod +x /usr/local/bin/timelapse

```

Now setup the [cron](https://en.wikipedia.org/wiki/Cron) job that will call our timelapse every 10 minutes:

```sh

echo "*/10 *  * * *   pi      DBX_TOKEN=[ Insert your dropbox token here ] NAME=[ Insert a name here ] /usr/local/bin/timelapse 2>&1" >> /etc/crontab

```

The `*/10` means "every 10 minutes".

To change the settings later, don't repeat the command, but open the file instead `nano /etc/crontab`, the line above should be the last in the file.

### Watermark

The watermark on the picture is located at `/home/pi/watermark-logo.png`, you can change to whatever you want. Keep it to support us :P

## Manage timelapses

All timelapses are stored on your dropbox's root in the `Apps/SuperGreenTimelapse/` directory.
The latest picture taken is named `latest.jpg`.

## Create timelapse

Creating the timelapse requires go to the `Apps/SuperGreenTimelapse/` directory in your Dropbox folder, then start the [create_timelapse.sh](https://raw.githubusercontent.com/supergreenlab/SuperGreenTimelapse/master/create_timelapse.sh) script.

```sh

curl -O https://raw.githubusercontent.com/supergreenlab/SuperGreenTimelapse/master/create_timelapse.sh
chmod +x create_timelapse.sh
./create_timelapse.sh [ The name of one of the timelapses ]

```

This will take a while to process. What it does is take each pics, create 4 versions to interpolate with `composite` then creates a video with all pics with `ffmpeg`.

The video will be written as `[ The name ].mp4`.

I remember some gotchas there, but I can't recall them, please post issues, or directly at [r/SuperGreenLab](https://www.reddit.com/r/SuperGreenLab).

# View live over http

First thing is to get a hosting solution.

Then [install docker](https://docs.docker.com/install/).

Then, run this command as root on your server:

```sh

docker run -d -p 80:80 -p 443:443 -e 'DBX_TOKEN=[ Insert your dropbox token here ]' --restart=always supergreenlab/supergreenlive

```

And now navigating to `http://[ your hosting IP or domain ]/[ The name you chose ]` will show the latest pic.

# One more thing:

This one's really live (in my office right now):

![Live](https://timelapse.chronic-o-matic.com/SuperGreenOffice#3 "Live")
