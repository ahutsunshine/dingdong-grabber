package notice

import (
	"testing"
)

func TestPlay(t *testing.T) {
	go func() {
		if err := NewDefaultMp3(); err != nil {
			panic(err)
		}
	}()
	t.Log("播放成功")
	// 取消select{}注释测试播放
	//select {}
}
