package notice

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"k8s.io/klog"
)

type Mp3 struct {
	Path string // mp3 path
}

func NewMp3(path string) NoticeInterface {
	return &Mp3{
		Path: path,
	}
}

func NewDefaultMp3() NoticeInterface {
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatal(err)
	}
	return NewMp3(fmt.Sprintf("%s%s", dir, "/music/everything_I_need.mp3"))
}

func (m *Mp3) Notify() error {
	audioFile, err := os.Open(m.Path)
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
