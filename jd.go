package golib

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
)

//IJingDong .
type IJingDong interface {
	VPR(source []byte, userID string, status string) (string, error)
}

//JingDonger .
type JingDonger struct {
	Appkey    string
	Secretkey string
}

var jingdong IJingDong

//JinggDong .
func JinggDong() IJingDong {
	if jingdong == nil {
		jingdong = &JingDonger{
			Appkey:    conf.GetConf("jd", "appkey"),
			Secretkey: conf.GetConf("jd", "secretkey"),
		}
	}
	return jingdong
}

//SetJingDong .
func SetJingDong(jd IJingDong) {
	jingdong = jd
}

//VPR 京东声纹识别按钮
func (jd *JingDonger) VPR(source []byte, userID string, status string) (string, error) {
	//header的字段
	type Encode struct {
		Channel    int    `json:"channel"`
		Format     string `json:"format"`
		SampleRate int    `json:"sample_rate"`
	}

	type Property struct {
		Autoend  bool   `json:"autoend"`
		Platform string `json:"platform"`
		Version  string `json:"version"`
		VprMode  string `json:"vpr_mode"`
		Encode   Encode `json:"encode"`
	}

	secretkey := jd.Secretkey
	//毫秒时间戳
	//golang没有直接生成毫秒时间戳的函数
	timestamp := time.Now().UnixNano() / 1e6
	time := strconv.FormatInt(timestamp, 10)
	//md5签名
	w := md5.New()
	io.WriteString(w, secretkey+time)
	sign := fmt.Sprintf("%x", w.Sum(nil))
	u := "https://aiapi.jd.com/jdai/vpr"
	values := url.Values{}
	values.Set("appkey", jd.Appkey)
	values.Set("timestamp", time)
	values.Set("sign", sign)
	u = fmt.Sprintf("%s?%s", u, values.Encode())
	header := make(map[string]string)
	header["Content-Type"] = "application/octet-stream"
	header["Request-Id"] = time
	header["User-Id"] = userID
	header["Sequence-Id"] = "-1"
	header["Server-Protocol"] = "-1"
	header["Net-State"] = "2"
	header["Applicator"] = "1"
	encode := Encode{
		Channel:    1,
		Format:     "wav",
		SampleRate: 16000,
	}
	property := Property{
		Autoend:  false,
		Platform: "Linux&Centos&7.3",
		Version:  "0.0.0.1",
		VprMode:  status,
		Encode:   encode,
	}
	b, err := json.Marshal(property)
	if err != nil {
		return "", errors.New("json转化失败")
	}
	header["Property"] = string(b)
	myHTTP := HTTP()
	body, err := myHTTP.Post(u, bytes.NewReader(source), header)
	if err != nil {
		return "", err
	}
	fmt.Println("调用", string(body))
	type Result struct {
		Text string `json:"text"`
	}
	type VPResult struct {
		Status int      `json:"status"`
		Result []Result `json:"result"`
	}
	type VPRReturn struct {
		Code   string   `json:"code"`
		Msg    string   `json:"msg"`
		Result VPResult `json:"result"`
	}
	var vpr VPRReturn
	json.Unmarshal(body, &vpr)
	if vpr.Code == "10000" && len(vpr.Result.Result) != 0 && vpr.Result.Status == 0 {
		result := vpr.Result.Result[0]
		fmt.Println(result.Text)
		return result.Text, nil
	}
	return "", errors.New(vpr.Msg)
}
