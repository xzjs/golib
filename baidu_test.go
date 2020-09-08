package lib

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/ini.v1"
)

func TestAudit(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}
	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)

	baidu := Baidu()

	text := "这是我的3D打印机"
	result, err := baidu.Audit(text)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if result != 1 {
		t.Error("审查结果错误")
	}
}

func TestASR(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}

	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)

	audio, err := ioutil.ReadFile("../public/16k.pcm")
	if err != nil {
		t.Error(err.Error())
		return
	}

	baidu := Baidu()
	text, err := baidu.ASR(audio, "1")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if text != "北京科技馆。" {
		t.Log(text)
		t.Error("识别失败")
	}
}

func TestTranslate(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}

	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)

	baidu = Baidu()
	text := "审查结果错误"
	res, err := baidu.Translate(text, true)
	t.Log(res)
	if err != nil {
		t.Error(err)
	}
	text = "Review result error"
	res, err = baidu.Translate(text, false)
	t.Log(res)
	if err != nil {
		t.Error(err)
	}
}

func TestProduSpeech(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}

	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)

	baidu := Baidu()
	text := "糟了！东西方向的绿灯时间太短，道路变得更拥挤了，还有两辆车发生了追尾！再重新分配一下绿灯时间吧"
	resp, err := baidu.ProduSpeech(text, "1")
	if err != nil {
		t.Error(err)
	}
	name := "../public/trafficWrong.mp3"
	if err := ioutil.WriteFile(name, resp, 0644); err != nil {
		t.Error(err)
	}
	fmt.Println("生成音频文件成功:" + name)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		t.Error("文件不存在")
	} else {
		if err := os.Remove(name); err != nil {
			t.Error("删除文件失败")
		}
	}
}
func TestChat(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}

	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)

	query := "你叫什么名字"
	baidu := Baidu()
	say, sessionID, err := baidu.Chat(query, "1", "")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if say == "" || sessionID == "" {
		t.Error("返回对话为空")
		return
	}
}

func TestOCR(t *testing.T) {
	if testing.Short() {
		t.Skip("单测环境下跳过")
	}

	cfg, _ := ini.Load("../conf.ini")
	conf := &Confer{
		Cfg: cfg,
	}
	SetConf(conf)
	data, err := ioutil.ReadFile("../tmp/case/test.png")
	if err != nil {
		t.Error(err.Error())
		return
	}
	image := base64.StdEncoding.EncodeToString(data)
	baidu := Baidu()
	result, err := baidu.OCR(image)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if result == "" {
		t.Error("未识别成功")
		return
	}
	t.Log(result)
}
