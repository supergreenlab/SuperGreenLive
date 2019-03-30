![SuperGreenLab](assets/sgl.png?raw=true "SuperGreenLab")

Table of Contents
=================

   * [SuperGreenTimelapse](#supergreentimelapse)
   * [Features](#features)
   * [Hardware requirements](#hardware-requirements)
   * [Quickstart](#quickstart)
      * [Installation](#installation)
         * [Dropbox setup](#dropbox-setup)
         * [RaspberryPI easy setup](#raspberrypi-easy-setup)
            * [Wifi configuration](#wifi-configuration)
            * [Dropbox](#dropbox)
            * [Timelapse on-screen display configuration](#timelapse-on-screen-display-configuration)
            * [Add temperature/humidity graphics](#add-temperaturehumidity-graphics)
         * [RaspberryPI hand setup](#raspberrypi-hand-setup)
         * [Watermark](#watermark)
      * [Manage timelapses](#manage-timelapses)
      * [Create timelapse](#create-timelapse)
   * [View live over http](#view-live-over-http)
   * [One more thing:](#one-more-thing)
   * [Subscribe our youtube channel:)](#subscribe-our-youtube-channel)

# SuperGreenTimelapse

Build a remote camera with a raspberryPI and webcam, and always have an eye on your growth.

Example:
[![Example](assets/screenshot-live.png?raw=true "Example")](https://live.supergreenlab.com/)

See live action: [https://live.supergreenlab.com/](https://live.supergreenlab.com/) (auto updates every 10 minutes)

# Features

- Live(-ish) secured webcam
- Take a picture every X minutes
- Upload to dropbox app folder
- Produce a timelapse video with extra image interpolation for better smoothness

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

### RaspberryPI easy setup

The most straight forward way to setup everything up is by using our [custom raspbian image](https://github.com/supergreenlab/SuperGreenLive/releases/download/1.0/image_2019-03-29-SuperGreenLiveOS-full.zip).

We'd recommend using something like [Etcher](https://www.balena.io/etcher/) for that.

[This tutorial](https://www.raspberrypi.org/documentation/installation/installing-images/) might help if you've never done that.

Then you'll first want to put the sd card in your real computer (ie not the raspberrypi), and mount it like any usb key.

The directory you'll see contains a bunch of files of great interests:

```
- wpa_supplicant.conf
- timelapse_dropbox_token
- timelapse_uploadname
- timelapse_name
- timelapse_strain
- timelapse_controllerid
```

#### Wifi configuration

Edit the wpa_supplicant.conf file, you'll have to enter your wifi credentials there.

Copy/paste this (don't forget to replace the values between []):

```

ctrl_interface=/run/wpa_supplicant                                          
update_config=1                                                                 
                                                                                            
network={                                   
        ssid="[ SSID ]"
        psk="[ PASSPHRASE ]"
}

```

#### Dropbox

Edit the `timelapse_dropbox_token` file and put the token, be sure not to have any spaces or empty lines.

Edit the `timelapse_uploadname` and put the name of the folder you want all pictures to be stored in. For example `SpaceTomatoes`.

#### Timelapse on-screen display configuration

The content of `timelapse_name` will be placed at the top left of the pictures.
`timelapse_strain` right under, some sort of a subtitle.

#### Add temperature/humidity graphics

For now this only works with the [SuperGreenController](https://github.com/supergreenlab/SuperGreenController), let me know if you're interested, I'll put more here:)

### RaspberryPI hand setup

(this is not up-to-date)

First follow the [raspberryPI quickstart guide](https://www.raspberrypi.org/learning/software-guide/quickstart/) if you have never done that before.

Always good to upgrade after a fresh install:
```sh

sudo apt-get update && sudo apt-get upgrade -y

```

Install or upgrade to the latest binary with the following command:

```sh

sudo curl https://github.com/supergreenlab/SuperGreenTimelapse/releases/download/PreRelease/timelapse -o /usr/local/bin/timelapse
sudo curl https://raw.githubusercontent.com/supergreenlab/SuperGreenTimelapse/master/watermark-logo.png -o /home/pi/watermark-logo.png
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

Those are really live (in my office right now):

My bloom box:
![Live](https://timelapse.chronic-o-matic.com/SuperGreenOffice#3 "Live")

[10 days timelapse](https://youtu.be/GGo-XaIuKoU)

And the veg boxes (there's two stacked on top of each other):
![LiveVeg](https://timelapse.chronic-o-matic.com/SuperGreenOfficeVeg#3 "LiveVeg")
![LiveVeg](https://timelapse.chronic-o-matic.com/SuperGreenOfficeVeg2#3 "LiveVeg2")

Btw the these 3 boxes use only 72w dispatched on 6 led panels, 4 in bloom, 2 in veg, the whole setup is controlled by only one [controller](https://github.com/supergreenlab/SuperGreenDriver) :) Ventilations are independantly controlled too.

# Subscribe our youtube channel:)

[![Youtube video:)](https://img.youtube.com/vi/0vjswZQ0rk4/0.jpg)](https://www.youtube.com/watch?v=0vjswZQ0rk4)
