package account

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// AuthKey JWT索引
type AuthKey string

const (
	// AuthInfoKey JWT用户键
	AuthInfoKey AuthKey = "auth-info-key"
)

const (
	// AuthTypeAdministrator --
	AuthTypeAdministrator = "administrator"
	// AuthTypeUser --
	AuthTypeUser = "user"
)

// AuthInfo --
type AuthInfo struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

// Authenticationable --
type Authenticationable interface {
	AuthInfo() *AuthInfo
}

// GenerateAuthToken 根据载荷生成token
func GenerateAuthToken(user Authenticationable, expires int64, secret string) (string, *jwt.StandardClaims, error) {
	jwtClaims, err := generateJWTClaims(user.AuthInfo(), expires)
	if err != nil {
		return "", jwtClaims, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", jwtClaims, err
	}
	return tokenString, jwtClaims, nil
}

// ParseAuthToken 根据token解析出载荷
func ParseAuthToken(tokenString string, secret string) (*AuthInfo, *jwt.StandardClaims, error) {
	jwtClaims, err := parseJWTClaims(tokenString, secret)
	if err != nil {
		return nil, nil, err
	}
	authInfo := &AuthInfo{}
	err = json.Unmarshal([]byte(jwtClaims.Subject), authInfo)
	if err != nil {
		return nil, nil, errors.New("cannot convert subject to jwtSubject")
	}
	return authInfo, jwtClaims, nil
}

// GetAuthToken --
func GetAuthToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

// RequestWithAuthInfo --
func RequestWithAuthInfo(r *http.Request, authInfo *AuthInfo) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), AuthInfoKey, authInfo))
}

// GetAuthInfo --
func GetAuthInfo(r *http.Request) (*AuthInfo, error) {
	value := r.Context().Value(AuthInfoKey)
	if value == nil {
		return nil, errors.New("jwtsubject error")
	}
	authInfo, ok := value.(*AuthInfo)
	if !ok {
		return nil, errors.New("jwtsubject type error")
	}
	return authInfo, nil
}

// GetAuthType 取JWT用户
func GetAuthType(r *http.Request) (string, error) {
	value := r.Context().Value(AuthInfoKey)
	if value == nil {
		return "", errors.New("jwtsubject error")
	}
	authInfo, ok := value.(*AuthInfo)
	if !ok {
		return "", errors.New("jwtsubject type error")
	}
	return authInfo.Type, nil
}

// GetAuthID 取JWT用户
func GetAuthID(r *http.Request, checkType string) (int64, error) {
	value := r.Context().Value(AuthInfoKey)
	if value == nil {
		return 0, errors.New("jwtsubject error")
	}
	authInfo, ok := value.(*AuthInfo)
	if !ok {
		return 0, errors.New("jwtsubject type error")
	}
	if authInfo.Type != checkType {
		return 0, errors.New("jwtsubject type check fail")
	}
	return authInfo.ID, nil
}

/*
// GenerateAuthInfo --
func GenerateAuthInfo(user Authenticationable) (*AuthInfo, error) {
	return user.AuthInfo(), nil
}
*/
