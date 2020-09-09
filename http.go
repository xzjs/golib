package golib

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

//IMyHTTP .
type IMyHTTP interface {
	Get(u string, values url.Values, header map[string]string) ([]byte, error)
	PostNeedHeader(u string, value io.Reader, header map[string]string) ([]byte, http.Header, error)
	Post(u string, value io.Reader, header map[string]string) ([]byte, error)
}

// MyHTTP .
type MyHTTP struct {
}

var myHTTP IMyHTTP

//HTTP .
func HTTP() IMyHTTP {
	if myHTTP == nil {
		myHTTP = &MyHTTP{}
	}
	return myHTTP
}

//SetMyHTTP 设置单例
func SetMyHTTP(http IMyHTTP) {
	myHTTP = http
}

// Get .
func (h *MyHTTP) Get(u string, values url.Values, header map[string]string) ([]byte, error) {
	u = fmt.Sprintf("%s?%s", u, values.Encode())
	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// PostNeedHeader .
func (h *MyHTTP) PostNeedHeader(u string, value io.Reader, header map[string]string) ([]byte, http.Header, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", u, value)
	if err != nil {
		return nil, nil, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return body, resp.Header, nil
}

// Post .
func (h *MyHTTP) Post(u string, value io.Reader, header map[string]string) ([]byte, error) {
	body, _, err := h.PostNeedHeader(u, value, header)
	return body, err
}
