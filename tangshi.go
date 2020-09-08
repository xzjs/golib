package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/jinzhu/gorm"
)

//ITangshi .
type ITangshi interface {
	Search(keytword string) (poetry *Poetry, err error)
}

//Tangshier .
type Tangshier struct {
}

var tangshi ITangshi

//Tangshi .
func Tangshi() ITangshi {
	if tangshi == nil {
		tangshi = &Tangshier{}
	}
	return tangshi
}

//SetTangshi .
func SetTangshi(t ITangshi) {
	tangshi = t
}

//Poetry .
type Poetry struct {
	gorm.Model
	Title        string `json:"title"`
	Type         string `json:"type"`
	Content      string `json:"content" gorm:"type:text"`
	Explanation  string `json:"explanation" gorm:"type:text"`
	Appreciation string `json:"appreciation" gorm:"type:text"`
	Author       string `json:"author"`
}

//SearchResult .
type SearchResult struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Result Result `json:"result"`
}

//Result .
type Result struct {
	List []Poetry `json:"list"`
}

//Search 查询古诗词
func (t *Tangshier) Search(keyword string) (poetry *Poetry, err error) {
	keyword = replace(keyword)
	u := "https://api.jisuapi.com/tangshi/search"
	values := url.Values{}
	conf := Conf()
	values.Set("appkey", conf.GetConf("tangshi", "appkey"))
	values.Set("pagesize", "1")
	values.Set("pagenum", "1")
	values.Set("keyword", keyword)
	u = fmt.Sprintf("%s?%s", u, values.Encode())
	// value := url.QueryEscape(keyword)
	myHTTP := HTTP()
	resp, err := myHTTP.Get(u, nil, nil)
	if err != nil {
		return nil, err
	}
	var response SearchResult
	json.Unmarshal(resp, &response)
	if response.Status != 0 || len(response.Result.List) == 0 {
		err = errors.New(string(resp))
		return
	}
	return &response.Result.List[0], nil
}

func replace(keyword string) string {
	pointDict := make(map[string]string, 0)
	pointDict[","] = "，"
	pointDict["."] = "。"
	for k, v := range pointDict {
		keyword = strings.Replace(keyword, k, v, -1)
	}
	return keyword
}
