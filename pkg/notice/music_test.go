package notice

import "testing"

func TestPlay(t *testing.T) {
	go func() {
		mp3 := &Mp3{}
		if err := mp3.Play("../../music/everything_I_need.mp3"); err != nil {
			t.Fatal(err)
		}
	}()
	t.Log("播放成功")
	// 测试是否 真正播放成功需要取消select{}注释
	//select {}
}
