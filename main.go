package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func main() {
	fmt.Println("Welcome to the Go File Renamer")

	//dir := "C:\\Users\\Thomas\\Documents\\Foto-Sync\\2022\\09 - Höhlenfahrt\\"
	dir := "./img/"
	entries, err := os.ReadDir(dir)
	handleFatal("reading directory entries", err)

	count := 0
	imgInfoArray := []imgInfo{}

	for _, v := range entries {
		imgFile, err := os.Open(dir + v.Name())
		handleFatal(fmt.Sprintf("opening file `%s`", v.Name()), err)

		var fileModDateTime, exifDateTime *time.Time
		fileInfo, err := v.Info()
		fileModDateTime = addressOfTime(fileInfo.ModTime())

		// movies dont have EXIF data so we just skip EXIF parsing if we get an error here
		metaData, err := exif.Decode(imgFile)
		if err == nil || !exif.IsCriticalError(err) {
			jsonByte, err := metaData.MarshalJSON()
			handleFatal(fmt.Sprintf("marshaling EXIF data as json for file `%s`", v.Name()), err)

			jsonString := string(jsonByte)
			exifDateTime = addressOfTime(parseJsonDateTime(jsonString, "DateTime"))
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
		var topic string = "Höhlenfahrt"
		newFilename := fmt.Sprintf("%d-%02d-%02d_%02d%02d_%s_%04d%s",
			v.dateTime.Year(), v.dateTime.Month(), v.dateTime.Day(),
			v.dateTime.Hour(), v.dateTime.Minute(),
			topic, count, filepath.Ext(v.path))
		count++
		newPath := dir + newFilename

		if _, err := os.Stat(newPath); err != nil {
			err = os.Rename(v.path, newPath)
			handleWarn(fmt.Sprintf("renaming `%s` to `%s`", v.path, newPath), err)
		} else {
			fmt.Printf("File with new target filename `%s` already exists -- ABORTING. Rerun with different topic string.\n", newPath)
			os.Exit(1)
		}
	}

	fmt.Printf("Renamed %d files\n", count)
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

func handleWarn(action string, err error) {
	if err != nil {
		fmt.Printf("Warning while %s: \"%s\"\n", action, err.Error())
	}
}

func parseJsonDateTime(jsonString, key string) time.Time {
	jsonValue := gjson.Get(jsonString, key)
	time, err := time.Parse("2006:01:02 15:04:05", jsonValue.Str)
	handleFatal(fmt.Sprintf("parsing date `%s` from EXIF attribute `%s`", jsonValue, key), err)
	return time
}

func addressOfTime(time time.Time) *time.Time {
	return &time
}
