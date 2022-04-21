package notice

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"k8s.io/klog"
)

type Mp3 struct {
}

func (m *Mp3) Play(path string) error {
	audioFile, err := os.Open(path)
	if err != nil {
		klog.Error(err)
		return err
	}
	defer audioFile.Close()

	// 对文件进行解码
	audioStreamer, format, err := mp3.Decode(audioFile)
	if err != nil {
		klog.Error(err)
		return err
	}

	defer audioStreamer.Close()
	done := make(chan bool)
	// 这里播放音乐
	// SampleRate is the number of samples per second. 采样率
	_ = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(audioStreamer, beep.Callback(func() {
		done <- true
	})))
	<-done
	return nil
}
