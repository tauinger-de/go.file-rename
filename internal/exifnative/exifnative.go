package exifnative

import (
	"github.com/rwcarlsen/goexif/exif"
	"gofire/internal/common"
	"os"
	"time"
)

func NewReader() common.ExifReader {
	return NativeExifReader{}
}

type NativeExifReader struct{}

func (r NativeExifReader) Get(name string, file *os.File) (string, error) {
	exifData, err := exif.Decode(file)
	if err == nil || !exif.IsCriticalError(err) {
		tag, err := exifData.Get(exif.FieldName(name))
		if err == nil {
			return tag.String(), nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func (r NativeExifReader) DateTimeOriginal(file *os.File) (time.Time, error) {
	exifData, err := exif.Decode(file)
	if err == nil || !exif.IsCriticalError(err) {
		dt, err := exifData.DateTime()
		if err == nil {
			return dt, nil
		} else {
			return time.Time{}, err
		}
	} else {
		return time.Time{}, err
	}
}
