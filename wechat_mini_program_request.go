package wechat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
)

// WechatMiniProgramLoginRequest 微信小程序登录请求
type WechatMiniProgramLoginRequest struct {
	Code string `json:"code"`
}

// Check --
func (request WechatMiniProgramLoginRequest) Check() bool {
	if request.Code == "" {
		return false
	}
	return true
}

// WechatMiniProgramUserInfoWatermark 微信小程序用户信息水印
type WechatMiniProgramUserInfoWatermark struct {
	AppID     string `json:"appid"`
	Timestamp uint64 `json:"timestamp"`
}

// WechatMiniProgramUserInfo 微信小程序用户信息
type WechatMiniProgramUserInfo struct {
	Nickname  string                             `json:"nickName"`
	AvatarURL string                             `json:"avatarUrl"`
	Gender    int8                               `json:"gender"`
	Province  string                             `json:"province"`
	City      string                             `json:"city"`
	Country   string                             `json:"country"`
	Language  string                             `json:"language"`
	OpenID    string                             `json:"openId"`
	UnionID   string                             `json:"unionId"`
	Watermark WechatMiniProgramUserInfoWatermark `json:"watermark"`
}

// WechatMiniProgramUserInfoRequest 微信小程序用户信息请求
type WechatMiniProgramUserInfoRequest struct {
	UserInfo      WechatMiniProgramUserInfo `json:"userInfo"`
	RawData       string                    `json:"rawData"`
	Signature     string                    `json:"signature"`
	EncryptedData string                    `json:"encryptedData"`
	IV            string                    `json:"iv"`
}

// Check --
func (request WechatMiniProgramUserInfoRequest) Check(cmd string) bool {
	if cmd == "raw" {
		if request.UserInfo.Nickname == "" || request.UserInfo.AvatarURL == "" {
			return false
		}
		if request.Signature == "" {
			return false
		}
	} else if cmd == "encrypted" {
		if request.EncryptedData == "" {
			return false
		}
		if request.IV == "" {
			return false
		}
	} else {
		return false
	}
	return true
}

// CheckSignature --
func (request WechatMiniProgramUserInfoRequest) CheckSignature(sessionKey string) bool {
	signature := fmt.Sprintf("%x", sha1.Sum([]byte(request.RawData+sessionKey)))
	return bytes.Equal([]byte(signature), []byte(request.Signature))
}

// Decrypt 解密
func (request WechatMiniProgramUserInfoRequest) Decrypt(sessionKey string) (*WechatMiniProgramUserInfo, error) {
	decodedData, err := base64.StdEncoding.DecodeString(request.EncryptedData)
	if err != nil {
		return nil, err
	}
	log.Println(request.EncryptedData)
	log.Println(decodedData)

	decodedIV, err := base64.StdEncoding.DecodeString(request.IV)
	if err != nil {
		return nil, err
	}
	log.Println(request.IV)
	log.Println(decodedIV)

	decodedSessionKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	log.Println(sessionKey)
	log.Println(decodedSessionKey)

	block, err := aes.NewCipher(decodedSessionKey)
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCBCDecrypter(block, decodedIV)
	originData := make([]byte, len(decodedData))
	decrypter.CryptBlocks(originData, decodedData)

	log.Println(string(originData))

	wechatMiniProgramUserInfo := &WechatMiniProgramUserInfo{}
	decoder := json.NewDecoder(bytes.NewReader(originData))
	if err = decoder.Decode(wechatMiniProgramUserInfo); err != nil {
		return nil, err
	}

	return wechatMiniProgramUserInfo, nil
}
