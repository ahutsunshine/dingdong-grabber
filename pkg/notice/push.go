package notice

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"

	"github.com/dingdong-grabber/pkg/constants"
	"github.com/dingdong-grabber/pkg/http"
	"k8s.io/klog"
)

type Push struct {
	token   string
	title   string
	content string
}

func NewPush(token, title, content string) NoticeInterface {
	return &Push{
		token:   token,
		title:   title,
		content: content,
	}
}

func (p *Push) Notify() error {
	marshal, err := json.Marshal(p)
	if err != nil {
		klog.Errorf("序列化推送内容失败，错误: %v", err)
		return err
	}

	client := http.NewClient(constants.Push)
	resp, err := client.RawPost(nil, nil, marshal)
	if err != nil {
		klog.Error(err)
		return err
	}

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusCreated {
		klog.Infof("推送返回不合法的状态值: %d", resp.StatusCode)
		return fmt.Errorf("%v", resp.StatusCode)
	}
	return nil
}
