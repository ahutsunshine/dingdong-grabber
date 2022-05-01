package util

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dingdong-grabber/pkg/constants"
	"k8s.io/klog"
)

func SignFilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		klog.Error(err)
		return "", err
	}

	return fmt.Sprintf("%s/%s", dir, constants.SignFile), nil
}

func SignConfigFilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		klog.Error(err)
		return "", err
	}

	return fmt.Sprintf("%s/%s", dir, constants.SignConfigFile), nil
}

func SignFile() (string, error) {
	file, err := SignFilePath()
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ClearSignConfigFile() {
	file, err := SignConfigFilePath()
	if err != nil {
		return
	}
	if err := os.Remove(file); err != nil {
		klog.Error(err)
	}
}
