package account

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

const (
	wechatMiniProgramCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?grant_type=authorization_code"
)

// DecodeWechatMiniProgramLoginRequest --
func DecodeWechatMiniProgramLoginRequest(r *http.Request) (*WechatMiniProgramLoginRequest, error) {
	decoder := json.NewDecoder(r.Body)
	request := &WechatMiniProgramLoginRequest{}
	if err := decoder.Decode(request); err != nil {
		return nil, err
	}
	if !request.Check() {
		return nil, errors.New("WechatMiniProgramLoginRequest check error!")
	}
	return request, nil
}

// DecodeWechatMiniProgramUserInfoRequest --
func DecodeWechatMiniProgramUserInfoRequest(r *http.Request) (*WechatMiniProgramUserInfoRequest, error) {
	decoder := json.NewDecoder(r.Body)
	request := &WechatMiniProgramUserInfoRequest{}
	if err := decoder.Decode(request); err != nil {
		return nil, err
	}
	if !request.Check("raw") && !request.Check("encrypted") {
		return nil, errors.New("WechatMiniProgramUserInfoRequest check error!")
	}
	return request, nil
}

// WechatMiniProgramCode2Session --
func WechatMiniProgramCode2Session(appID string, appSecret string, code string) (*WechatCode2SessionResponse, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: transport,
	}
	url := wechatMiniProgramCode2SessionURL + "&appid=" + appID + "&secret=" + appSecret + "&js_code=" + code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("code2session start %s\n", time.Now())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("code2session end   %s\n", time.Now())
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	response := &WechatCode2SessionResponse{}
	if err := decoder.Decode(response); err != nil {
		return nil, err
	}
	if response.ErrCode != 0 {
		return nil, errors.New(response.ErrMsg)
	}
	return response, nil
}
