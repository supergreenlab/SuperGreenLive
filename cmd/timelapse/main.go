/*
 * Copyright (C) 2019  SuperGreenLab <towelie@supergreenlab.com>
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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/sirupsen/logrus"
	"gopkg.in/gographics/imagick.v2/imagick"
)

var dbx files.Client

var (
	uploadname      string
	boxname         string
	strain          string
	graphcontroller string
	graphbox        int
	uploadpath      string
	rotate          bool
)

func init() {
	flag.StringVar(&uploadname, "u", "SuperGreenKit", "Name for the box (used to upload)")
	flag.StringVar(&boxname, "n", "SuperGreenKit - bloom", "Name for the box (written on top of pic)")
	flag.StringVar(&strain, "s", "Bagseed", "Strain name")
	flag.StringVar(&graphcontroller, "c", "", "Graph's controller id")
	flag.IntVar(&graphbox, "b", 0, "Graph's controller box id")
	flag.BoolVar(&rotate, "r", false, "")

	flag.Parse()

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
	if rotate {
		cmd := exec.Command("/usr/bin/raspistill", "-vf", "-hf", "-q", "50", "-o", name)
		err := cmd.Run()
		return name, err
	}
	cmd := exec.Command("/usr/bin/raspistill", "-q", "50", "-o", name)
	err := cmd.Run()
	return name, err
}

func addText(mw *imagick.MagickWand, text, color string, size, stroke, x, y float64) {
	pw := imagick.NewPixelWand()
	defer pw.Destroy()

	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	dw.SetFont("plume.otf")
	dw.SetFontSize(size)

	pw.SetColor("white")
	dw.SetStrokeColor(pw)
	dw.SetStrokeWidth(stroke)

	pw.SetColor(color)
	dw.SetFillColor(pw)
	dw.Annotation(x, y, text)

	mw.DrawImage(dw)
}

func addPic(mw *imagick.MagickWand, file string, x, y float64) {
	pic := imagick.NewMagickWand()
	defer pic.Destroy()

	pic.ReadImage(file)

	dw := imagick.NewDrawingWand()
	dw.Composite(imagick.COMPOSITE_OP_ATOP, x, y, float64(pic.GetImageWidth()), float64(pic.GetImageHeight()), pic)

	mw.DrawImage(dw)
}

type MetricValue [][]float64

type Metrics struct {
	Metrics MetricValue
}

func (mv MetricValue) minMax() (float64, float64) {
	min := math.MaxFloat64
	max := math.SmallestNonzeroFloat32

	for _, v := range mv {
		min = math.Min(min, v[1])
		max = math.Max(max, v[1])
	}

	return min, max
}

func (mv MetricValue) current() float64 {
	if len(mv) < 1 {
		return 0
	}
	return mv[len(mv)-1][1]
}

func loadGraphValue(controller, metric string) Metrics {
	m := Metrics{}

	url := fmt.Sprintf("https://api2.supergreenlab.com/?cid=%s&q=%s&t=24&n=50", controller, metric)
	r, err := http.Get(url)
	if err != nil {
		return m
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&m)
	return m
}

func drawGraphLine(mw *imagick.MagickWand, pts []imagick.PointInfo, color string) {
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	cw := imagick.NewPixelWand()
	defer cw.Destroy()

	dw.SetStrokeAntialias(true)
	dw.SetStrokeWidth(2)
	dw.SetStrokeLineCap(imagick.LINE_CAP_ROUND)
	dw.SetStrokeLineJoin(imagick.LINE_JOIN_ROUND)

	cw.SetColor(color)
	dw.SetStrokeColor(cw)

	cw.SetColor("none")
	dw.SetFillColor(cw)

	dw.Polyline(pts)

	mw.DrawImage(dw)
}

func drawGraphBackground(mw *imagick.MagickWand, pts []imagick.PointInfo, color string) {
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	cw := imagick.NewPixelWand()
	defer cw.Destroy()

	dw.SetStrokeAntialias(true)
	dw.SetStrokeWidth(2)
	dw.SetStrokeLineCap(imagick.LINE_CAP_ROUND)
	dw.SetStrokeLineJoin(imagick.LINE_JOIN_ROUND)

	cw.SetColor("none")
	dw.SetStrokeColor(cw)

	cw.SetColor(color)
	cw.SetOpacity(0.4)
	dw.SetFillColor(cw)

	dw.Polygon(pts)

	mw.DrawImage(dw)
}

func addGraph(mw *imagick.MagickWand, x, y, width, height, min, max float64, mv MetricValue, color string) {
	var (
		spanX = width / float64(len(mv)-1)
	)

	pts := make([]imagick.PointInfo, 0, len(mv)+2)
	pts = append(pts, imagick.PointInfo{
		X: x, Y: y,
	})
	for i, v := range mv {
		pts = append(pts, imagick.PointInfo{
			X: x + float64(i)*spanX,
			Y: y - ((v[1] - min) * (height - 60) / (max - min)),
		})
	}
	pts = append(pts, imagick.PointInfo{
		X: x + width, Y: y,
	})

	drawGraphBackground(mw, []imagick.PointInfo{
		{x, y}, {x + width, y}, {x + width, y - height}, {x, y - height},
	}, "white")
	drawGraphLine(mw, pts[1:len(pts)-1], color)
	drawGraphBackground(mw, pts, color)

	cw := imagick.NewPixelWand()
	defer cw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	cw.SetColor("white")
	dw.SetStrokeColor(cw)
	dw.SetStrokeWidth(3)
	dw.Line(x, y, x, y-height)
	dw.Line(x, y, x+width, y)

	mw.DrawImage(dw)
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
	cmd := exec.Command("/usr/bin/mogrify", cam, "-quality", size, name)
	err := cmd.Run()
	return name, err
}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	logrus.Info("Taking picture..")
	cam, err := takePic()
	if err != nil {
		log.Println(err)
	}

	logrus.Info("Uploading raw file")
	remote := fmt.Sprintf("%d.jpg", int32(time.Now().Unix()))
	uploadPic(uploadname+"_raw", cam, remote)

	mw.ReadImage(cam)

	addText(mw, boxname, "#3BB30B", 120, 5, 25, 120)
	addText(mw, strain, "#FF4B4B", 80, 3, 25, 220)

	if graphcontroller != "" {
		t := loadGraphValue(graphcontroller, fmt.Sprintf("BOX_%d_TEMP", graphbox))
		h := loadGraphValue(graphcontroller, fmt.Sprintf("BOX_%d_HUMI", graphbox))
		var (
			x = float64(25)
			y = float64(mw.GetImageHeight() - 25)
		)
		addGraph(mw, x, y, 350, 200, 16, 40, t.Metrics, "#3BB30B")
		addText(mw, fmt.Sprintf("%d°", int(t.Metrics.current())), "#3BB30B", 150, 7, x+65, y-110)

		addGraph(mw, x+365, y, 400, 200, 20, 80, h.Metrics, "#0B81B3")
		addText(mw, fmt.Sprintf("%d%%", int(h.Metrics.current())), "#0B81B3", 150, 7, x+390, y-110)
	}

	t := time.Now()
	d := t.Format("2006/01/02 15:04")
	addText(mw, d, "#3BB30B", 90, 4, float64(mw.GetImageWidth()-900), float64(mw.GetImageHeight()-60))

	addPic(mw, "watermark-logo.png", float64(mw.GetImageWidth()-330), 10)

	mw.WriteImage("latest.jpg")

	logrus.Info("Uploading files")
	uploadPic(uploadname, "latest.jpg", remote)

	logrus.Info("Resizing latest")
	latest, _ := resizeLatest("latest.jpg", "50%")
	uploadPic(uploadname, latest, "latest.jpg")
}
