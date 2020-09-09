package golib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
)

//IWeChat .
type IWeChat interface {
	GetAccessToken() (token string, err error)
	Audit(text string) error
}

//WeChater .
type WeChater struct {
}

var weixin IWeChat

//WeiXin .
func WeiXin() IWeChat {
	if weixin == nil {
		weixin = &WeChater{}
	}
	return weixin
}

//SetWeiXin .
func SetWeiXin(wx IWeChat) {
	weixin = wx
}

// GetAccessToken 获取accesstoken
func (wx *WeChater) GetAccessToken() (token string, err error) {
	token, _ = Cache().Get("WechatToken")
	if token == "" {
		type response struct {
			AccessToken string `json:"access_token"`
			ExpiresIN   int    `json:"expires_in"`
			ErrorCode   int    `json:"errcode"`
			ErrorMsg    string `json:"errmsg"`
		}
		conf := Conf()
		var atr response
		values := url.Values{}
		values.Set("grant_type", "client_credential")
		values.Set("appid", conf.GetConf("wechat", "appid"))
		values.Set("secret", conf.GetConf("wechat", "appsecret"))
		u := "https://api.weixin.qq.com/cgi-bin/token"
		myHTTP := HTTP()
		resp, err := myHTTP.Get(u, values, nil)
		if err != nil {
			return "", err
		}
		json.Unmarshal(resp, &atr)
		if atr.ErrorCode != 0 || atr.AccessToken == "" {
			return "", errors.New(atr.ErrorMsg)
		}
		if err = Cache().Set("WechatToken", atr.AccessToken, atr.ExpiresIN); err != nil {
			return "", err
		}
		token = atr.AccessToken
	}
	return
}

// AuditReturn .
type AuditReturn struct {
	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

// Audit 文本审核，返回值为审核结果类型，可取值1.合规，2.不合规，3.疑似，4.审核失败
func (wx *WeChater) Audit(text string) error {
	if len(text) > 255 {
		return errors.New("字数超过限制，最大为255")
	}
	u := "https://api.weixin.qq.com/wxa/msg_sec_check"
	values := url.Values{}
	accessToken, err := wx.GetAccessToken()
	if err != nil {
		return errors.New("token获取失败")
	}
	values.Set("access_token", accessToken)
	u = fmt.Sprintf("%s?%s", u, values.Encode())
	body := make(map[string]string)
	header := make(map[string]string)
	header["contentType"] = "application/json;charset=utf-8"
	body["content"] = text
	bodyByte, _ := json.Marshal(body)
	myHTTP := HTTP()
	resp, err := myHTTP.Post(u, bytes.NewReader(bodyByte), nil)
	if err != nil {
		return errors.New("调用微信审核接口失败")
	}
	log.Println(string(resp))
	var result AuditReturn
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return err
	}
	if result.ErrorCode != 0 {
		return errors.New(result.ErrorMsg)
	}
	return nil
}
