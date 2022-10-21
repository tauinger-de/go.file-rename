package exiftool

import (
	"errors"
	"fmt"
	"github.com/barasher/go-exiftool"
	"gofire/common"
	"os"
	"time"
)

func NewReader() (common.ExifReader, error) {
	newExiftool, err := exiftool.NewExiftool()
	if err != nil {
		return nil, err
	} else {
		r := ExifToolReader{
			exiftool: newExiftool,
		}
		return r, nil
	}
}

type ExifToolReader struct {
	exiftool *exiftool.Exiftool
}

func (r ExifToolReader) Get(name string, file *os.File) (string, error) {
	return "", nil
}

func (r ExifToolReader) DateTimeOriginal(file *os.File) (time.Time, error) {
	metadata, _ := r.getMetadata(file)       // todo error
	v, _ := metadata.GetString("CreateDate") // todo error
	dt, _ := parseDateTime(v)                // todo error
	return dt, nil
}

func (r ExifToolReader) getMetadata(file *os.File) (*exiftool.FileMetadata, error) {
	metadata := r.exiftool.ExtractMetadata(file.Name())
	if len(metadata) == 1 {
		return &metadata[0], nil
	} else {
		return nil, errors.New(fmt.Sprintf("Got %d instead of 1 result", len(metadata)))
	}
}

func parseDateTime(v string) (time.Time, error) {
	return time.Parse("2006:01:02 15:04:05", v)
}
