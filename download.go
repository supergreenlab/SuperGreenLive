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
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/sirupsen/logrus"
)

var (
	uploadname    string
	min_timestamp int64
	max_timestamp int64
	n_pics        int
)

func init() {
	flag.StringVar(&uploadname, "u", "SuperGreenKit", "Name of the timelapse")
	flag.Int64Var(&min_timestamp, "min", -1, "[optional] min timestamp")
	flag.Int64Var(&max_timestamp, "max", -1, "[optional] max timestamp")
	flag.IntVar(&n_pics, "n", 0, "[optional] number of pics")

	flag.Parse()
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

func download(dbx files.Client, dbxFile string) error {
	if _, err := os.Stat(dbxFile[1:]); os.IsExist(err) {
		logrus.Info("Already exists: %s", dbxFile)
		return nil
	}

	logrus.Info(dbxFile)
	da := files.NewDownloadArg(dbxFile)

	_, c, err := dbx.Download(da)
	if err != nil {
		return err
	}

	outFile, err := os.Create(dbxFile[1:])
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, c)
	logrus.Infof("Downloaded %s", dbxFile)
	return nil
}

func main() {
	token := MustGetenv("DBX_TOKEN")
	if uploadname == "" {
		logrus.Fatal("missing timelapse uploadname arg")
	}

	config := dropbox.Config{
		Token: token,
	}
	dbx := files.New(config)
	path := fmt.Sprintf("/%s", uploadname)
	args := files.NewListFolderArg(path)
	res, err := dbx.ListFolder(args)
	fu(err)

	fs := []*files.FileMetadata{}
	for {
		for _, m := range res.Entries {
			switch f := m.(type) {
			case *files.FileMetadata:
				if f.Size < 1800000 {
					continue
				}
				if min_timestamp != -1 && max_timestamp != -1 {
					t, _ := strconv.ParseInt(strings.Split(f.Name, ".")[0], 10, 64)
					if !(t > min_timestamp && t < max_timestamp) {
						continue
					}
				}
				fs = append(fs, f)
			}
		}
		if res.HasMore {
			logrus.Infof("Listed %d files", len(fs))
			res, err = dbx.ListFolderContinue(files.NewListFolderContinueArg(res.Cursor))
			fu(err)
		} else {
			break
		}
	}
	logrus.Infof("Listed %d files", len(fs))

	if n_pics != 0 {
		fs = fs[len(fs)-n_pics:]
	}

	var wg sync.WaitGroup
	for i, f := range fs {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			for i := 0; i < 4; i++ {
				if err := download(dbx, f); err != nil {
					logrus.Warning(err)
					continue
				}
				break
			}
		}(fmt.Sprintf("%s/%s", path, f.Name))
		if i%20 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}
