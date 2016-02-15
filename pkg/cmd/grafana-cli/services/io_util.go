package services

import (
	m "github.com/grafana/grafana/pkg/cmd/grafana-cli/models"
	"io/ioutil"
	"os"
)

var IoUtil m.IoUtil = IoUtilImp{}

type IoUtilImp struct {
}

func (i IoUtilImp) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (i IoUtilImp) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (i IoUtilImp) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}

func (i IoUtilImp) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
