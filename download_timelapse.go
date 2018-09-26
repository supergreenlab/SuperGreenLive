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
	"os"
	"sync"

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

func download(dbx files.Client, dbxFile string) {
	logrus.Info(dbxFile)
	da := files.NewDownloadArg(dbxFile)

	_, c, err := dbx.Download(da)
	fu(err)

	outFile, err := os.Create(dbxFile[1:])
	fu(err)
	defer outFile.Close()

	_, err = io.Copy(outFile, c)
	logrus.Infof("Downloaded %s", dbxFile)
}

func main() {
	token := MustGetenv("DBX_TOKEN")
	name := ""
	if len(os.Args) >= 2 {
		name = os.Args[1]
	} else {
		logrus.Fatal("missing timelapse name arg")
	}

	config := dropbox.Config{
		Token: token,
	}
	dbx := files.New(config)
	path := fmt.Sprintf("/%s", name)
	args := files.NewListFolderArg(path)
	res, err := dbx.ListFolder(args)
	fu(err)

	var wg sync.WaitGroup
	i := 0
	for {
		for _, m := range res.Entries {
			switch f := m.(type) {
			case *files.FileMetadata:
				logrus.Infof("f %s", f.PathLower)
				wg.Add(1)
				go func(f string) {
					defer wg.Done()
					download(dbx, f)
				}(fmt.Sprintf("%s/%s", path, f.Name))
			}
			i++
			if i%20 == 0 {
				wg.Wait()
			}
		}
		if res.HasMore {
			res, err = dbx.ListFolderContinue(files.NewListFolderContinueArg(res.Cursor))
			fu(err)
		} else {
			break
		}
	}
}
