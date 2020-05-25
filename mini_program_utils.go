package weixin

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const (
	miniProgramCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?grant_type=authorization_code"
)

// DecodeMiniProgramLoginRequest --
func DecodeMiniProgramLoginRequest(r *http.Request) (*MiniProgramLoginRequest, error) {
	decoder := json.NewDecoder(r.Body)
	request := &MiniProgramLoginRequest{}
	if err := decoder.Decode(request); err != nil {
		return nil, err
	}
	return request, nil
}

// DecodeMiniProgramUserInfoRequest --
func DecodeMiniProgramUserInfoRequest(r *http.Request) (*MiniProgramUserInfoRequest, error) {
	decoder := json.NewDecoder(r.Body)
	request := &MiniProgramUserInfoRequest{}
	if err := decoder.Decode(request); err != nil {
		return nil, err
	}
	return request, nil
}

// MiniProgramCode2Session --
func MiniProgramCode2Session(appID string, appSecret string, code string) (*MiniProgramCode2SessionResponse, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: transport,
	}
	url := miniProgramCode2SessionURL + "&appid=" + appID + "&secret=" + appSecret + "&js_code=" + code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	response := &MiniProgramCode2SessionResponse{}
	if err := decoder.Decode(response); err != nil {
		return nil, err
	}
	if response.ErrCode != 0 {
		return nil, errors.New(response.ErrMsg)
	}
	return response, nil
}
