package lib

import (
	"io/ioutil"
	"testing"

	"gopkg.in/ini.v1"
)

func TestVPR(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}
	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)

	audio, err := ioutil.ReadFile("../tmp/case/result.wav")
	if err != nil {
		t.Error(err.Error())
		return
	}
	jingdong := JinggDong()
	r, err := jingdong.VPR(audio, "1", "test")
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}
