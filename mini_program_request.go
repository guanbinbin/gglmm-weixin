package weixin

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// MiniProgramLoginRequest 微信小程序登录请求
type MiniProgramLoginRequest struct {
	Code string `json:"code"`
}

// Check --
func (request MiniProgramLoginRequest) Check() bool {
	if request.Code == "" {
		return false
	}
	return true
}

// MiniProgramUserInfo 微信小程序用户信息
type MiniProgramUserInfo struct {
	Nickname  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
	Gender    int8   `json:"gender"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Language  string `json:"language"`
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
}

// MiniProgramUserInfoWatermark 微信小程序用户信息水印
type MiniProgramUserInfoWatermark struct {
	AppID     string `json:"appid"`
	Timestamp uint64 `json:"timestamp"`
}

// MiniProgramUserInfoRequest 微信小程序用户信息请求
type MiniProgramUserInfoRequest struct {
	UserInfo struct {
		MiniProgramUserInfo
		Watermark MiniProgramUserInfoWatermark
	} `json:"userInfo"`
	RawData       string `json:"rawData"`
	Signature     string `json:"signature"`
	EncryptedData string `json:"encryptedData"`
	IV            string `json:"iv"`
}

// Check --
func (request MiniProgramUserInfoRequest) Check(cmd string) bool {
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
func (request MiniProgramUserInfoRequest) CheckSignature(sessionKey string) bool {
	signature := fmt.Sprintf("%x", sha1.Sum([]byte(request.RawData+sessionKey)))
	return bytes.Equal([]byte(signature), []byte(request.Signature))
}

// Decrypt 解密
func (request MiniProgramUserInfoRequest) Decrypt(sessionKey string) (*MiniProgramUserInfo, error) {
	decodedData, err := base64.StdEncoding.DecodeString(request.EncryptedData)
	if err != nil {
		return nil, err
	}
	decodedIV, err := base64.StdEncoding.DecodeString(request.IV)
	if err != nil {
		return nil, err
	}
	decodedSessionKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(decodedSessionKey)
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCBCDecrypter(block, decodedIV)
	originData := make([]byte, len(decodedData))
	decrypter.CryptBlocks(originData, decodedData)
	miniProgramUserInfo := &MiniProgramUserInfo{}
	decoder := json.NewDecoder(bytes.NewReader(originData))
	if err = decoder.Decode(miniProgramUserInfo); err != nil {
		return nil, err
	}
	return miniProgramUserInfo, nil
}
