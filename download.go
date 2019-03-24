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
	"strconv"
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
	name := ""
	n_pics := 0
	if len(os.Args) >= 3 {
		n_pics, _ = strconv.Atoi(os.Args[2])
	}
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

	fs := []*files.FileMetadata{}
	for {
		for _, m := range res.Entries {
			switch f := m.(type) {
			case *files.FileMetadata:
				if f.Size < 1500000 {
					continue
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
}
