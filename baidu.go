package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// IBaidu baidu的接口
type IBaidu interface {
	GetAccessToken() (token string, err error)
	Audit(text string) (conclusionType int, err error)
	ASR(audio []byte, cuid string) (string, error)
	Translate(text string, IsChinese bool) (res string, err error)
	Chat(query string, userID string, sessionID string) (say string, backSessionID string, err error)
	ProduSpeech(text string, cuid string) ([]byte, error)
	OCR(image string) (result string, err error)
}

// Baiduer .
type Baiduer struct {
}

var baidu IBaidu

// Baidu 获取单例
func Baidu() IBaidu {
	if baidu == nil {
		baidu = &Baiduer{}
	}
	return baidu
}

// SetBaidu 设置单例
func SetBaidu(bd IBaidu) {
	baidu = bd
}

// GetAccessToken 获取accesstoken
func (bd *Baiduer) GetAccessToken() (token string, err error) {
	token, _ = Cache().Get("AccessToken")
	if token == "" {
		type response struct {
			AccessToken      string `json:"access_token"`
			ExpiresIN        int    `json:"expires_in"`
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		conf := Conf()
		var atr response
		values := url.Values{}
		values.Set("client_id", conf.GetConf("baidu", "key"))
		values.Set("client_secret", conf.GetConf("baidu", "secret"))
		values.Set("grant_type", "client_credentials")
		u := "https://aip.baidubce.com/oauth/2.0/token"
		myHTTP := HTTP()
		resp, err := myHTTP.Get(u, values, nil)
		if err != nil {
			return "", err
		}
		json.Unmarshal(resp, &atr)
		if atr.Error != "" || atr.AccessToken == "" {
			return "", errors.New(atr.Error + ":" + atr.ErrorDescription)
		}
		if err = Cache().Set("AccessToken", atr.AccessToken, atr.ExpiresIN); err != nil {
			return "", err
		}
		token = atr.AccessToken
	}
	return
}

// AuditResult .
type AuditResult struct {
	LogID          int64       `json:"log_id"`
	ConclusionType int         `json:"conclusionType"`
	Data           interface{} `json:"data"`
	ErrorCode      int         `json:"error_code"`
	ErrorMsg       string      `json:"error_msg"`
}

// Audit 文本审核，返回值为审核结果类型，可取值1.合规，2.不合规，3.疑似，4.审核失败
func (bd *Baiduer) Audit(text string) (conclusionType int, err error) {
	u := "https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined"
	values := url.Values{}
	accessToken, err := bd.GetAccessToken()
	if err != nil {
		return
	}
	values.Set("access_token", accessToken)
	u = fmt.Sprintf("%s?%s", u, values.Encode())
	header := make(map[string]string)
	header["contentType"] = "application/x-www-form-urlencoded"
	value := url.QueryEscape(text)
	myHTTP := HTTP()
	resp, err := myHTTP.Post(u, strings.NewReader("text="+value), header)
	if err != nil {
		return
	}
	var result AuditResult
	json.Unmarshal(resp, &result)
	if result.ErrorMsg != "" {
		return 4, errors.New(result.ErrorMsg)
	}
	return result.ConclusionType, nil
}

// ASR 语音识别
func (bd *Baiduer) ASR(audio []byte, cuid string) (string, error) {
	type ASRResponse struct {
		ErrorNo int      `json:"error_no"`
		ErrMsg  string   `json:"err_msg"`
		SN      string   `json:"sn"`
		Result  []string `json:"result"`
	}
	u := "http://vop.baidu.com/server_api"
	values := url.Values{}
	accessToken, err := bd.GetAccessToken()
	if err != nil {
		return "", err
	}
	values.Set("token", accessToken)
	values.Set("cuid", cuid)
	u = fmt.Sprintf("%s?%s", u, values.Encode())
	header := make(map[string]string)
	header["Content-Type"] = "audio/pcm;rate=16000"
	myHTTP := HTTP()
	resp, err := myHTTP.Post(u, bytes.NewReader(audio), header)
	if err != nil {
		return "", err
	}
	var response ASRResponse
	json.Unmarshal(resp, &response)
	if response.ErrorNo != 0 || len(response.Result) == 0 {
		return "", errors.New(string(resp))
	}
	return response.Result[0], nil
}

//Translate 语言翻译
func (bd *Baiduer) Translate(text string, IsChinese bool) (res string, err error) {
	var (
		formLan string = "zh"
		toLan   string = "en"
	)
	if !IsChinese {
		formLan = "en"
		toLan = "zh"
	}
	//调用百度API接口
	conf := Conf()
	appid := conf.GetConf("fanyi", "appid")
	secretKey := conf.GetConf("fanyi", "secret")
	salt := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000)
	signStr := appid + text + string(salt) + secretKey
	//md5加密签名
	w := md5.New()
	io.WriteString(w, signStr)
	sign := fmt.Sprintf("%x", w.Sum(nil))
	u := "http://api.fanyi.baidu.com/api/trans/vip/translate"
	v := url.Values{}
	v.Set("q", text)
	v.Set("from", formLan)
	v.Set("to", toLan)
	v.Set("appid", appid)
	v.Set("salt", string(salt))
	v.Set("sign", sign)
	myHTTP := HTTP()
	response, err := myHTTP.Get(u, v, nil)
	if err != nil {
		return "", err
	}
	//转化utf-8的编码
	respU, err := UnescapeUnicode(response)
	if err != nil {
		return "", errors.New("转化编码失败")
	}

	//TransResult 翻译结果
	type TransResultModel struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	}
	type Resp struct {
		TransResult []TransResultModel `json:"trans_result"`
		ErrorCode   int                `json:"error_code"`
	}
	var resp Resp
	err = json.Unmarshal(respU, &resp)
	if err != nil {
		fmt.Println("json转换错误")
		return
	}
	//如果发生错误，response会有error_msg
	if resp.ErrorCode != 0 {
		return "发生错误", errors.New(string(respU))
	}
	var dst string
	for _, tranResult := range resp.TransResult {
		dst = tranResult.Dst
	}
	return dst, nil
}

// Chat 聊天接口
func (bd *Baiduer) Chat(query string, userID string, sessionID string) (say string, backSessionID string, err error) {
	type ChatRequest struct {
		UserID string `json:"user_id"`
		Query  string `json:"query"`
	}

	type Chater struct {
		Version   string      `json:"version"`
		ServiceID string      `json:"service_id"`
		LogID     string      `json:"log_id"`
		SessionID string      `json:"session_id"`
		Request   ChatRequest `json:"request"`
	}

	type ChatAction struct {
		Say string `json:"say"`
	}

	type ChatResponse struct {
		Status     int          `json:"status"`
		Msg        string       `json:"msg"`
		ActionList []ChatAction `json:"action_list"`
	}

	type ChatResult struct {
		SessionID    string         `json:"session_id"`
		ResponseList []ChatResponse `json:"response_list"`
	}

	type ChatReturn struct {
		ErrorCode int        `json:"error_code"`
		ErrorMsg  string     `json:"error_msg"`
		Result    ChatResult `json:"result"`
	}
	u := "https://aip.baidubce.com/rpc/2.0/unit/service/chat"
	values := url.Values{}
	accessToken, err := bd.GetAccessToken()
	if err != nil {
		return
	}
	values.Set("access_token", accessToken)
	u = fmt.Sprintf("%s?%s", u, values.Encode())
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	request := ChatRequest{
		UserID: userID,
		Query:  query,
	}
	chat := Chater{
		Version:   "2.0",
		ServiceID: "S31887",
		LogID:     uuid.New().String(),
		SessionID: sessionID,
		Request:   request,
	}
	chatStr, _ := json.Marshal(chat)
	myHTTP := HTTP()
	resp, err := myHTTP.Post(u, bytes.NewReader(chatStr), header)
	if err != nil {
		return
	}
	var response ChatReturn
	json.Unmarshal(resp, &response)
	if response.ErrorCode != 0 || len(response.Result.ResponseList) == 0 {
		err = errors.New(string(resp))
		return
	}
	for _, r := range response.Result.ResponseList {
		if r.Status == 0 {
			for _, a := range r.ActionList {
				return a.Say, response.Result.SessionID, nil
			}
		}
	}
	err = errors.New("没有合适的对话")
	return
}

//ProduSpeech 由文本合成语音
func (bd *Baiduer) ProduSpeech(text string, cuid string) ([]byte, error) {
	if len(text) > 2048 {
		return nil, errors.New("字符串太长超过限制")
	}
	//调用百度api语音合成接口
	u := "http://tsn.baidu.com/text2audio"
	values := url.Values{}
	accessToken, err := bd.GetAccessToken()
	if err != nil {
		return nil, err
	}
	values.Set("tok", accessToken)
	values.Set("cuid", cuid)
	values.Set("ctp", "1")
	values.Set("lan", "zh")
	values.Set("per", "111")
	values.Set("tex", url.QueryEscape(text))
	myHTTP := HTTP()
	resp, header, err := myHTTP.PostNeedHeader(u, strings.NewReader(values.Encode()), nil)
	if err != nil {
		return nil, err
	}
	if header.Get("Content-Type") == "audio/mp3" {
		return resp, nil
	}
	return nil, errors.New(string(resp))
}

//OCR .
func (bd *Baiduer) OCR(image string) (result string, err error) {

	//Probability .
	type Probability struct {
		Variance float64 `json:"variance"`
		Average  float64 `json:"average"`
		Min      float64 `json:"min"`
	}
	//ORCWord .
	type ORCWord struct {
		Words       string      `json:"words"`
		Probability Probability `json:"Probability"`
	}
	//OCRReturn .
	type OCRReturn struct {
		ErrorCode  int       `json:"error_code"`
		ErrorMsg   string    `json:"error_msg"`
		LogID      int       `json:"log_id"`
		WordNUm    int       `json:"words_result_num"`
		WordResult []ORCWord `json:"words_result"`
		Direction  int       `json:"direction"`
	}
	u := "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic"
	values := url.Values{}
	accessToken, err := bd.GetAccessToken()
	if err != nil {
		return
	}
	values.Set("access_token", accessToken)
	u = fmt.Sprintf("%s?%s", u, values.Encode())
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	value := url.Values{}
	value.Set("probability", "true")
	value.Set("image", image)
	myHTTP := HTTP()
	resp, err := myHTTP.Post(u, strings.NewReader(value.Encode()), header)
	if err != nil {
		return "", err
	}
	var response OCRReturn
	json.Unmarshal(resp, &response)
	if response.ErrorCode != 0 || len(response.WordResult) == 0 {
		err = errors.New(string(resp))
		return
	}

	w := response.WordResult[0]
	for _, v := range response.WordResult {
		if v.Probability.Average > w.Probability.Average {
			w = v
		}
	}
	return w.Words, nil

}
