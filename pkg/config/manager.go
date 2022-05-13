package config

import (
	"errors"
	"fmt"

	"github.com/dingdong-grabber/pkg/util"
	"gopkg.in/yaml.v3"
	"k8s.io/klog"
)

type Manager struct {
	Conf *Conf
}

func (m *Manager) LoadConfig() error {
	dir, err := util.GetRootDir()
	if err != nil {
		return err
	}

	dir = fmt.Sprintf("%s/%s", dir, "config.yaml")
	data, err := util.ReadFile(dir)
	if err != nil {
		return err
	}

	if err = m.Decode(data); err != nil {
		return err
	}
	return m.Validate()
}

func (m *Manager) Decode(data []byte) error {
	if err := yaml.Unmarshal(data, &m.Conf); err != nil {
		klog.Error(err)
		return err
	}
	return nil
}

func (m *Manager) Validate() error {
	if m.Conf == nil {
		return errors.New("根目录config.yaml必须配置参数")
	}
	if m.Conf.Config.Cookie == "" && m.Conf.Config.Device != IosType && m.Conf.Config.Strategy != 3 {
		return errors.New("非测试模式(strategy!=3)或者使用android或default device时，请求头cookie为必填项")
	}
	return nil
}
