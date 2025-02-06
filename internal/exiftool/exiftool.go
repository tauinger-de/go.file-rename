package exiftool

import (
	"errors"
	"fmt"
	"github.com/barasher/go-exiftool"
	"gofire/internal/common"
	"os"
	"strings"
	"time"
)

func NewReader() (common.ExifReader, error) {
	newExiftool, err := exiftool.NewExiftool(exiftool.Charset("filename=utf8"))
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

func (r ExifToolReader) Get(_ string, _ *os.File) (string, error) {
	return "", nil
}

func (r ExifToolReader) DateTimeOriginal(file *os.File) (time.Time, error) {
	metadata, err := r.getMetadata(file)
	if err == nil {
		// .MOV files have "CreationDate" with timezone and "CreateDate" without!
		v, err := metadata.GetString("CreationDate")
		if err == nil {
			return parseDateTime(v)
		} else {
			fields := []string{"DateTimeOriginal", "FileModifyDate"}
			for _, f := range fields {
				v, err := metadata.GetString(f)
				if err == nil {
					return parseDateTime(v)
				}
			}
			err = fmt.Errorf("EXIF data has none of these keys: %s",
				strings.Join(fields, ","))
		}
	}
	return time.Time{}, err
}

func (r ExifToolReader) getMetadata(file *os.File) (*exiftool.FileMetadata, error) {
	metadata := r.exiftool.ExtractMetadata(file.Name())
	if len(metadata) == 1 {
		if metadata[0].Err != nil {
			return nil, metadata[0].Err
		} else {
			return &metadata[0], nil
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Got %d instead of 1 result", len(metadata)))
	}
}

func parseDateTime(v string) (time.Time, error) {
	formats := []string{"2006:01:02 15:04:05", "2006:01:02 15:04:05-07:00"}
	for _, f := range formats {
		t, err := time.Parse(f, v)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("none of the attempted datetime formats matched")
}
