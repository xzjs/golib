package lib

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"testing"

	"gopkg.in/ini.v1"
)

//
func TestFindFace(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}

	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)

	data, err := ioutil.ReadFile("../tmp/case/source.png")
	if err != nil {
		t.Error(err.Error())
		return
	}
	source := base64.StdEncoding.EncodeToString(data)
	data, err = ioutil.ReadFile("../tmp/case/target.png")
	if err != nil {
		t.Error(err.Error())
		return
	}
	target := base64.StdEncoding.EncodeToString(data)

	position, err := Tencet().FindFace(source, target)
	bytea, _ := json.Marshal(position)
	t.Log(string(bytea))
	if err != nil {
		t.Error(err.Error())
		return
	}
}
