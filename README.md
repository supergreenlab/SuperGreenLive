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
         * [View latest](#view-latest)
      * [Manage timelapses](#manage-timelapses)
         * [List](#list)
         * [Download](#download)
         * [Delete](#delete)
      * [Create timelapse](#create-timelapse)

# SuperGreenTimelapse

A bunch of scripts/programs to produce dropbox-backed timelapses for raspberryPi.

![Example](assets/example.jpg?raw=true "Example")

# Features

- Live(-ish) secured webcam
- Take a picture every X minutes
- Upload to dropbox hidden directory
- Produce a timelapse video with extra image interpolation for better smoothness
- Dropbox hidden folder management

## TODO

- Integrate https://github.com/gographics/imagick

# Hardware requirements

- [RaspberryPI](https://www.raspberrypi.org/products/) + [Wifi (optional, most rpi have integrated wifi now)](https://www.raspberrypi.org/products/raspberry-pi-usb-wifi-dongle/)
- [Camera](https://www.raspberrypi.org/products/camera-module-v2/)
- [Power supply](https://www.raspberrypi.org/products/raspberry-pi-universal-power-supply/)

# Quickstart

## Installation

### Dropbox setup

The problem that arises when you want to take timelapses is that taking pictures every 10 minutes takes a lot of space.
And having a raspberrypi running 24/24 and storing big amounts of data on an SD card is looking for trouble:P

So what seems to be a good solution is to upload everything to dropbox.

It also allows to view the latest pic online, which actually makes it some sort of cloud live camcorder. Good times.

There's a little setup to do on dropbox's side. For obvious security purpose you have to let dropbox know that he needs to create a space and access for our program.

Got to the [app creation page](https://www.dropbox.com/developers/apps/create), and choose: `Dropbox API` -> `App folder` -> `SuperGreenTimelapse` (or whatever on that last one :P).

Now scroll to the `Generated access token` section, and click the `Generate` button below. Copy-paste the long id that looks like `vrB4PlxSQpsAAAAAAAC1SvJJbXi08sdjlkaWWfalk25iX4GAqsfk67rkM0sM0uyC`, we'll need that in the next step.

### RaspberryPI setup

First follow the [raspberryPI quickstart guide](https://www.raspberrypi.org/learning/software-guide/quickstart/) if you have never done that before.

Install or upgrade to the latest binary with the following command:

```sh

sudo curl https://github.com/supergreenlab/SuperGreenTimelapse/releases/download/PreRelease/timelapse -o /usr/local/bin/timelapse
sudo curl https://github.com/supergreenlab/SuperGreenTimelapse/releases/download/PreRelease/watermark-logo.png -o /home/pi/watermark-logo.png
sudo chmod +x /usr/local/bin/timelapse

```

Now setup the `cron` jog that will call our timelapse every 10 minutes:

```sh

echo "*/10 *  * * *   pi      DBX_TOKEN=[ Insert your dropbox token here ] NAME=[ Insert a name here ] /usr/local/bin/timelapse 2>&1" >> /etc/crontab

```

The `*/10` means "every 10 minutes".

To change the settings later, don't repeat the command, but open the file instead `nano /etc/crontab`, the line above should be the last in the file.

### Watermark

The watermark on the picture is located at `/home/pi/watermark-logo.png`, you can change to whatever you want. Keep it to support us :P

### View latest

You can view the latest pic taken, by downloading the [latest](https://github.com/supergreenlab/SuperGreenTimelapse/releases/tag/PreRelease) binary (pick the right one for your OS.
This binary is meant to be running on a server, but can still be used locally. It opens a webserver and serves the latest pictures of the given timelapse.

#### Locally

Just run the executable by double clicking it.

Then use your browser to go to `http://localhost:8080/[ Insert the name you chose ]`.

#### Online

First thing is to get a hosting solution.

Then [install docker]().

Then, run this command as root on your server:

```sh

docker run -d -p 80:80 -p 443:443 -e 'DBX_TOKEN=[ Insert your dropbox token here ]' --restart=always supergreenlab/SuperGreenTimelapse

```

## Manage timelapses

### List

### Download

### Delete

## Create timelapse

