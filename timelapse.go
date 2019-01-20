/*
 * Copyright (C) 2018  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/sirupsen/logrus"
)

var dbx files.Client

func init() {
	token := MustGetenv("DBX_TOKEN")
	config := dropbox.Config{
		Token: token,
	}

	dbx = files.New(config)
}

func fu(e error) {
	if e != nil {
		logrus.Fatal(e)
	}
}

func MustGetenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		logrus.Fatalf("missing env %s", name)
	}
	return v
}

func takePic() (string, error) {
	name := "cam.jpg"
	cmd := exec.Command("/usr/bin/raspistill", "-rot", "180", "-vf", "-hf", "-q", "50", "-o", name)
	err := cmd.Run()
	return name, err
}

func makeTextLayer(text string) (string, error) {
	label := fmt.Sprintf("label:%s", text)
	name := "layer.png"
	cmd := exec.Command("/usr/bin/convert", "-background", "transparent", "-fill", "green", "-font", "Helvetica-Narrow", "-pointsize", "80", "-stroke", "white", "-strokewidth", "2", label, name)
	err := cmd.Run()
	return name, err
}

func composeLayers(cam, layer, gravity string) (string, error) {
	name := "latest.jpg"
	cmd := exec.Command("/usr/bin/convert", cam, layer, "-gravity", gravity, "-geometry", "+20+20", "-composite", name)
	err := cmd.Run()
	return name, err
}

func uploadPic(name, local, remote string) {
	f, err := os.Open(local)
	fu(err)

	p := fmt.Sprintf("/%s/%s", name, remote)
	ci := files.NewCommitInfo(p)
	ci.Mode.Tag = "overwrite"
	_, err = dbx.Upload(ci, f)
	fu(err)

	logrus.Infof("Uploaded %s", p)
}

func resizeLatest(cam, size string) (string, error) {
	name := "latest.jpg"
	cmd := exec.Command("/usr/bin/convert", cam, "-scale", size, name)
	err := cmd.Run()
	return name, err
}
func main() {
	name := MustGetenv("NAME")

	logrus.Info("Taking picture..")
	cam, err := takePic()
	fu(err)

	logrus.Info("Adding watermark..")
	t := time.Now()
	d := t.Format("2006/01/02 15:04")
	date, err := makeTextLayer(d)
	fu(err)
	local, err := composeLayers(cam, date, "southeast")
	fu(err)

	namel, err := makeTextLayer(name)
	fu(err)
	local, err = composeLayers(local, namel, "northwest")
	fu(err)

	local, err = composeLayers(local, "watermark-logo.png", "southwest")
	fu(err)

	logrus.Info("Uploading files")
	remote := fmt.Sprintf("%d.jpg", int32(time.Now().Unix()))
	uploadPic(name, local, remote)

	logrus.Info("Resizing latest")
	latest, err := resizeLatest(local, "33%")
	uploadPic(name, latest, "latest.jpg")
}
