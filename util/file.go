package util

import (
	"io/ioutil"
	"os"
)

var WorkingPath, _ = os.Getwd()

func MustReadFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}

func MustDeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func MustWriteFile(path string, data []byte) {
	if err := os.WriteFile(path, data, 0666); err != nil {
		panic(err)
	}
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
