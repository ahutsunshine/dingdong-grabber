package notice

import (
	"testing"
)

func TestPlay(t *testing.T) {
	go func() {
		mp3 := &Mp3{}
		if err := mp3.Play("../../music/everything_I_need.mp3"); err != nil {
			panic(err)
		}
	}()
	t.Log("播放成功")
	// 取消select{}注释测试播放
	//select {}
}
