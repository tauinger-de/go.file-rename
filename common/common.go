package common

import (
	"fmt"
	"os"
	"time"
)

type ExifReader interface {
	Get(name string, file *os.File) (string, error)
	DateTimeOriginal(file *os.File) (time.Time, error)
}

func HandleFatal(action string, err error) {
	if err != nil {
		fmt.Printf("Error while %s: \"%s\"\n", action, err.Error())
		os.Exit(1)
	}
}

func HandleWarn(action string, err error) bool {
	if err != nil {
		fmt.Printf("Warning while %s: \"%s\"\n", action, err.Error())
		return true
	} else {
		return false
	}
}
