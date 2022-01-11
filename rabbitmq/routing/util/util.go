package util

import (
	"log"
	"os"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func IsFile(path string) bool {
	if f, err := os.Stat(path); err == nil {
		return f.Mode().IsRegular()
	}
	return false
}
