package lib

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//ITecent .
type ITecent interface {
	FindFace(string, string) (*Position, error)
}

// Tecenter ..
type Tecenter struct {
	AppID  string
	AppKey string
}

var tecent ITecent

// Tencet 获取单例
func Tencet() ITecent {
	conf := Conf()
	if tecent == nil {
		tecent = &Tecenter{
			AppID:  conf.GetConf("tecent", "appid"),
			AppKey: conf.GetConf("tecent", "appkey"),
		}
	}
	return tecent
}

//SetTencent 设置单例
func SetTencent(t ITecent) {
	tecent = t
}

// Position 位置
type Position struct {
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
	X2 float64 `json:"x2"`
	Y2 float64 `json:"y2"`
}

//
// FindFace 查找人脸 sourceImage和targetImage为图片的base64字符串
func (tecent *Tecenter) FindFace(sourceImage string, targetImage string) (position *Position, err error) {
	params := url.Values{}
	params.Set("app_id", tecent.AppID)
	params.Set("time_stamp", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("nonce_str", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("source_image", sourceImage)
	params.Set("target_image", targetImage)
	params.Set("sign", getSign(params))
	u := "https://api.ai.qq.com/fcgi-bin/face/face_detectcrossageface"
	str := params.Encode()
	header := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	myHTTP := HTTP()
	resp, err := myHTTP.Post(u, strings.NewReader(str), header)
	if err != nil {
		return
	}

	type Data struct {
		SourceFace Position `json:"source_face"`
		TargetFace Position `json:"target_face"`
		Score      float64  `json:"score"`
		FailFlag   int      `json:"fail_flag"`
	}
	type Result struct {
		Ret  int    `json:"ret"`
		Msg  string `json:"msg"`
		Data Data   `json:"data"`
	}
	var result Result
	json.Unmarshal(resp, &result)
	if result.Ret != 0 && result.Msg != "" {
		err = errors.New(string(resp))
		return
	}
	if result.Data.Score > 0.8 {
		return &result.Data.TargetFace, nil
	}
	return
}

// getSign 获取签名
func getSign(params url.Values) (sign string) {
	sign = params.Encode()
	t := tecent.(*Tecenter)
	sign = fmt.Sprintf("%s&%s=%s", sign, "app_key", t.AppKey)
	w := md5.New()
	io.WriteString(w, sign)
	sign = fmt.Sprintf("%x", w.Sum(nil))
	sign = strings.ToUpper(sign)
	return
}
