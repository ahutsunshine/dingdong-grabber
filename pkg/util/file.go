package util

import (
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
