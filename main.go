package main

import (
	"flag"
	"fmt"
	"gofire/common"
	"gofire/exifnative"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const dateTimePattern string = "%d-%02d-%02d %02d.%02d.%02d"

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
	common.HandleFatal("reading directory entries", err)
	fmt.Printf("Found %d directory entries\n", len(entries))

	count := 0
	var imgInfoArray []imgInfo
	exifReader := exifnative.NewReader()

	for _, v := range entries {
		// open file
		imgFile, err := os.Open(dir + v.Name())
		common.HandleFatal(fmt.Sprintf("opening file `%s`", v.Name()), err)
		fileInfo, err := v.Info()
		common.HandleFatal(fmt.Sprintf("retrieving file details for `%s`", v.Name()), err)

		// skip dirs
		if fileInfo.IsDir() {
			continue
		}

		// stuff
		var filenameDateTime, modifiedDateTime, exifDateTime *time.Time

		// get date-time from filename
		//regexp := regexp.MustCompile("(\\d{4})-(\\d{2})-(\\d{2})\\s(\\d{2})\\.(\\d{2})\\.(\\d{2}).*")
		var year, month, day, hour, min, sec int
		_, err = fmt.Sscanf(v.Name(), dateTimePattern, &year, &month, &day, &hour, &min, &sec)
		if err == nil {
			// alternativ time.Parse() nach regexp match auf erwartetes format und substring
			filenameDateTime = addressOfTime(
				time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local),
			)
		} else {
			filenameDateTime = nil
		}

		// get modification time
		modifiedDateTime = addressOfTime(fileInfo.ModTime())

		// movies dont have EXIF data so we just skip EXIF parsing if we get an error here
		dt, err := exifReader.DateTimeOriginal(imgFile)
		if !common.HandleWarn(fmt.Sprintf("getting 'DateTime' EXIF entry for file `%s`", imgFile.Name()), err) {
			exifDateTime = &dt
		}

		var finalDateTime *time.Time = nil
		if exifDateTime == nil {
			if filenameDateTime == nil {
				finalDateTime = modifiedDateTime
			} else {
				finalDateTime = filenameDateTime
			}
		} else {
			finalDateTime = exifDateTime
		}

		err = imgFile.Close()
		common.HandleWarn(fmt.Sprintf("closing file `%s`", imgFile.Name()), err)

		imgInfoArray = append(imgInfoArray, imgInfo{
			path:     dir + v.Name(),
			dateTime: *finalDateTime,
		})
	}

	sort.Sort(imgInfoList(imgInfoArray))

	for _, v := range imgInfoArray {
		// build new filename
		newFilename := fmt.Sprintf(dateTimePattern+" %s-%04d%s",
			v.dateTime.Year(), v.dateTime.Month(), v.dateTime.Day(),
			v.dateTime.Hour(), v.dateTime.Minute(), v.dateTime.Second(),
			*topic, count, filepath.Ext(v.path))
		count++
		newPath := dir + newFilename

		// skip if name doesn't change
		if v.path == newPath {
			continue
		}

		// check for existing file
		if _, err := os.Stat(newPath); err != nil {
			// rename it!
			err = os.Rename(v.path, newPath)
			common.HandleWarn(fmt.Sprintf("renaming `%s` to `%s`", v.path, newPath), err)
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

func addressOfTime(t time.Time) *time.Time {
	return &t
}
