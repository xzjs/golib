package golib

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"gopkg.in/ini.v1"
)

func TestEncode(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}

	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)
	text := "今天天气真好啊"
	a := Encode(text)
	t.Log(a)
}

func TestDecode(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}
	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)
	a := "zHLeFsS/lFYp5mwYYu9Ac1kopz3FbkwhMBBL0+IkbNk="
	b := Decode(a)
	t.Log(b)
	if b != "今天天气真好啊" {
		t.Error("解密失败")
	}
}

func TestRand(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}
	a := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000000000)
	key := strconv.FormatInt(a, 10)
	if len(key) != 16 {
		t.Error("faile")
	}
	t.Log(key)
}
