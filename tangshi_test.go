package golib

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}

	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)
	tangshi := Tangshi()
	text := "白日依山尽，黄河入海流"
	p, err := tangshi.Search(text)
	if err != nil {
		t.Error(err)
	}
	if p == nil {
		t.Error("识别失败")
	}
	t.Log(p.Appreciation)
}
