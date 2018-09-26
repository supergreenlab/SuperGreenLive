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

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/sirupsen/logrus"
)

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

func main() {
	token := MustGetenv("DBX_TOKEN")
	name := ""
	if len(os.Args) >= 2 {
		name = os.Args[1]
	} else {
		logrus.Fatal("missing arg name")
	}

	config := dropbox.Config{
		Token: token,
	}
	dbx := files.New(config)

	local := "image_00001.jpg"
	remote := "lol.jpg"

	f, err := os.Open(local)
	fu(err)

	ci := files.NewCommitInfo(fmt.Sprintf("/%s/%s", name, remote))
	res, err := dbx.Upload(ci, f)
	fu(err)

	logrus.Info(res)
}
