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
	"io"
	"net/http"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var dbx files.Client

var (
	dbxtoken = pflag.String("dbxtoken", "", "Dropbox token")
)

func init() {
	viper.SetDefault("DBXToken", "")
}

func getFileReader(name string) (io.ReadCloser, error) {
	dbxFile := fmt.Sprintf("/%s/latest.jpg", name)

	da := files.NewDownloadArg(dbxFile)

	_, c, err := dbx.Download(da)

	logrus.Infof("Downloaded %s", dbxFile)
	return c, err
}

func serve(c *gin.Context) {
	name := c.Param("name")

	content, err := getFileReader(name)
	if err != nil {
		logrus.Warning(err)
		c.Status(http.StatusNotFound)
		return
	}

	c.Header("Content-Type", "image/jpeg")
	c.Header("Cache-Control", "no-cache")
	c.Header("ETag", fmt.Sprintf("%d", int64(time.Now().Unix())/(60*10)))
	c.Status(http.StatusOK)
	io.Copy(c.Writer, content)
}

func main() {
	viper.SetConfigName("livecam")
	viper.AddConfigPath("/etc/livecam")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	viper.SetEnvPrefix("LIVECAM")
	viper.AutomaticEnv()

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	token := viper.GetString("DBXToken")
	config := dropbox.Config{
		Token: token,
	}
	dbx = files.New(config)

	r := gin.Default()
	r.GET("/:name", serve)
	r.Run(":3000")
}
