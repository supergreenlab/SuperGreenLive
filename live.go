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
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var dbx files.Client

func init() {
	token := os.Getenv("DBX_TOKEN")
	if token == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter dropbox token: ")
		token, _ = reader.ReadString('\n')
		token = token[:len(token)-1]
	}
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

func GetFileReader(name string) (io.ReadCloser, error) {
	dbxFile := fmt.Sprintf("/%s/latest.jpg", name)

	da := files.NewDownloadArg(dbxFile)

	_, c, err := dbx.Download(da)

	logrus.Infof("Downloaded %s", dbxFile)
	return c, err
}

func serve(c *gin.Context) {
	name := c.Param("name")

	content, err := GetFileReader(name)
	if err != nil {
		logrus.Warning(err)
		c.Status(http.StatusNotFound)
		return
	}

	c.Header("Content-Type", "image/jpeg")
	c.Status(http.StatusOK)
	io.Copy(c.Writer, content)
}

func main() {
	r := gin.Default()
	r.GET("/:name", serve)

	certFile := "certs/supergreenlive.com.crt"
	keyFile := "certs/supergreenlive.com.key"
	if _, err := os.Stat(certFile); err == nil {
		logrus.Info("exists1")
		if _, err := os.Stat(keyFile); err == nil {
			logrus.Info("exists2")
			go r.RunTLS(":443", certFile, keyFile)
		}
	}
	if len(os.Args) == 2 {
		r.Run(fmt.Sprintf(":%s", os.Args[1]))
	} else {
		r.Run(":8080")
	}
}
