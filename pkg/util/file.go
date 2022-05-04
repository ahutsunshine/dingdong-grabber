package util

import (
	"io/ioutil"
	"os"

	"k8s.io/klog"
)

func GetRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		klog.Error(err)
		return "", err
	}
	return dir, nil
}

func ReadFile(file string) ([]byte, error) {
	if data, err := ioutil.ReadFile(file); err != nil {
		klog.Error(err)
		return nil, err
	} else {
		return data, nil
	}
}
