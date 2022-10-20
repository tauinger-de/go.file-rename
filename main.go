package main

import (
	"flag"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	fmt.Println("Welcome to the Go File Renamer")

	sourceDir := flag.String("source", ".", "Specifies the folder of the images to rename. Default is current directory.")
	topic := flag.String("topic", "media", "Defines the topic string to be included in the filename.")
	flag.Parse()

	dir := filepath.ToSlash(*sourceDir)
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	entries, err := os.ReadDir(dir)
	handleFatal("reading directory entries", err)

	count := 0
	var imgInfoArray []imgInfo

	for _, v := range entries {
		// open file
		imgFile, err := os.Open(dir + v.Name())
		handleFatal(fmt.Sprintf("opening file `%s`", v.Name()), err)

		// get modification time
		var fileModDateTime, exifDateTime *time.Time
		fileInfo, err := v.Info()
		fileModDateTime = addressOfTime(fileInfo.ModTime())

		// skip dirs
		if fileInfo.IsDir() {
			continue
		}

		// movies dont have EXIF data so we just skip EXIF parsing if we get an error here
		metaData, err := exif.Decode(imgFile)
		if err == nil || !exif.IsCriticalError(err) {
			dt, err := metaData.DateTime()
			if !handleWarn("getting 'DateTime' EXIF entry", err) {
				exifDateTime = &dt
			}
		} else {
			handleWarn(fmt.Sprintf("decoding EXIF for `%s` -- no EXIF available", v.Name()), err)
		}

		var finalDateTime *time.Time = nil
		if exifDateTime == nil {
			finalDateTime = fileModDateTime
		} else {
			finalDateTime = exifDateTime
		}

		err = imgFile.Close()
		handleWarn(fmt.Sprintf("closing file `%s`", imgFile.Name()), err)

		imgInfoArray = append(imgInfoArray, imgInfo{
			path:     dir + v.Name(),
			dateTime: *finalDateTime,
		})
	}

	sort.Sort(imgInfoList(imgInfoArray))

	for _, v := range imgInfoArray {
		newFilename := fmt.Sprintf("%d-%02d-%02d_%02d%02d_%s_%04d%s",
			v.dateTime.Year(), v.dateTime.Month(), v.dateTime.Day(),
			v.dateTime.Hour(), v.dateTime.Minute(),
			*topic, count, filepath.Ext(v.path))
		count++
		newPath := dir + newFilename

		// skip if name doesn't change
		if v.path == newPath {
			continue
		}

		// check for existing file
		if _, err := os.Stat(newPath); err != nil {
			err = os.Rename(v.path, newPath)
			handleWarn(fmt.Sprintf("renaming `%s` to `%s`", v.path, newPath), err)
		} else {
			fmt.Printf("File with new target filename `%s` already exists -- ABORTING. Rerun with different topic string.\n", newPath)
			os.Exit(1)
		}
	}

	fmt.Printf("Processed %d files\n", count)
}

type imgInfo struct {
	path     string
	dateTime time.Time
}

type imgInfoList []imgInfo

func (l imgInfoList) Len() int {
	return len(l)
}

func (l imgInfoList) Less(i, j int) bool {
	return l[i].dateTime.Before(l[j].dateTime)
}

func (l imgInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func handleFatal(action string, err error) {
	if err != nil {
		fmt.Printf("Error while %s: \"%s\"\n", action, err.Error())
		os.Exit(1)
	}
}

func handleWarn(action string, err error) bool {
	if err != nil {
		fmt.Printf("Warning while %s: \"%s\"\n", action, err.Error())
		return true
	} else {
		return false
	}
}

func addressOfTime(t time.Time) *time.Time {
	return &t
}
