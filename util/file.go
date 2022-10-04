package util

import (
	"os"
)

var WorkingPath, _ = os.Getwd()

func MustReadFile(path string) []byte {
	return SelectAnyByteSlice(os.ReadFile(path))
}

func MustDeleteFile(path string) {
	Must(os.RemoveAll(path))
}

func MustWriteFile(path string, data []byte) {
	Must(os.WriteFile(path, data, 0666))
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
