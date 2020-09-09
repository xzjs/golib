package golib

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestWXAudit(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}
	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)

	wx := WeiXin()

	text := "圣人做庄警儒教，浏览器打开威尼斯警连五肖"
	err := wx.Audit(text)
	if err != nil {
		t.Error(err.Error())
		return
	}
}
