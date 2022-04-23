package sign

import (
	"encoding/json"
	"io/ioutil"

	"github.com/robertkrimen/otto"
	"k8s.io/klog"
)

type JsSign struct {
	file string
}

func NewDefaultJsSign() SignInterface {
	return NewSign("../../sign.js")
}

func NewSign(file string) SignInterface {
	return &JsSign{
		file: file,
	}
}

func (s *JsSign) Sign(data interface{}) (map[string]string, error) {
	bytes, err := ioutil.ReadFile(s.file)
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	vm := otto.New()
	if _, err = vm.Run(string(bytes)); err != nil {
		klog.Error(err)
		return nil, err
	}
	bytes, _ = json.Marshal(data)
	value, err := vm.Call("sign", nil, string(bytes))
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	var signs map[string]string
	if err = json.Unmarshal([]byte(value.String()), &signs); err != nil {
		klog.Errorf("解析签名结果出错，错误: %v", err)
		return nil, err
	}
	return signs, nil
}
