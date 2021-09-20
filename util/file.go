package util

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

var WorkingPath, _ = os.Getwd()

func ReadFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.WithError(err).WithField("util", "ReadFile").Errorf("unable to read '%s'", path)
		return nil
	}
	return bytes
}

func DeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logrus.WithError(err).WithField("util", "DeleteFile").Errorf("unable to delete '%s'", path)
	}
}

func WriteFile(path string, data []byte) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logrus.WithError(err).WithField("util", "WriteFile").Errorf("unable to open '%s'", path)
		return
	}
	_, err = io.WriteString(f, string(data))
	if err != nil {
		logrus.WithError(err).WithField("util", "WriteFile").Errorf("unable to write '%s'", path)
		return
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
