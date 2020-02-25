package account

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

// LoginRequest 登录请求
type LoginRequest struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

// Check --
func (request LoginRequest) Check() bool {
	if request.Mobile == "" {
		return false
	}
	if request.Password == "" {
		return false
	}
	return true
}

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code string `json:"code"`
}

// Check --
func (request WechatLoginRequest) Check() bool {
	if request.Code == "" {
		return false
	}
	return true
}

// WechatUserInfoWatermark --
type WechatUserInfoWatermark struct {
	AppID     string `json:"appid"`
	Timestamp uint64 `json:"timestamp"`
}

// WechatUserInfo --
type WechatUserInfo struct {
	Nickname  string                  `json:"nickName"`
	AvatarURL string                  `json:"avatarUrl"`
	Gender    int8                    `json:"gender"`
	Province  string                  `json:"province"`
	City      string                  `json:"city"`
	Country   string                  `json:"country"`
	Language  string                  `json:"language"`
	OpenID    string                  `json:"openId"`
	UnionID   string                  `json:"unionId"`
	Watermark WechatUserInfoWatermark `json:"watermark"`
}

// WechatUserInfoRequest 微信用户信息请求
type WechatUserInfoRequest struct {
	UserInfo      WechatUserInfo `json:"userInfo"`
	RawData       string         `json:"rawData"`
	Signature     string         `json:"signature"`
	EncryptedData string         `json:"encryptedData"`
	IV            string         `json:"iv"`
}

// Check --
func (request WechatUserInfoRequest) Check(cmd string) bool {
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
	}
	return true
}

// CheckSignature --
func (request WechatUserInfoRequest) CheckSignature(sessionKey string) bool {
	signature := fmt.Sprintf("%x", sha1.Sum([]byte(request.RawData+sessionKey)))
	return bytes.Equal([]byte(signature), []byte(request.Signature))
}

// Decrypt 解密
func (request WechatUserInfoRequest) Decrypt(sessionKey string) (*WechatUserInfo, error) {
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

	wechatUserInfo := &WechatUserInfo{}
	decoder := json.NewDecoder(bytes.NewReader(originData))
	if err = decoder.Decode(wechatUserInfo); err != nil {
		return nil, err
	}

	return wechatUserInfo, nil
}

// LogoutRequest 退出请求
type LogoutRequest struct {
	Token string
}
