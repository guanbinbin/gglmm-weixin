package account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTKey JWT索引
type JWTKey string

const (
	// JWTClaimsKey JWT用户键
	JWTClaimsKey JWTKey = "jwt-claims-key"
)

const (
	// JWTExpires JWT失效时间
	JWTExpires int64 = 24 * 60 * 60
)

// JWTSubject --
type JWTSubject interface {
	ID() int64
	ToString() (string, error)
	FromString(value string) error
}

const (
	//UserTypeAdministrator --
	UserTypeAdministrator = "administrator"
	// UserTypeNormal --
	UserTypeNormal = "normal"
)

// JWTUser --
type JWTUser struct {
	UserType      string
	UserID        int64
	UserName      string
	UserAvatarURL string
}

// ID --
func (jwtUser *JWTUser) ID() int64 {
	return jwtUser.UserID
}

// ToString --
func (jwtUser *JWTUser) ToString() (string, error) {
	result, err := json.Marshal(jwtUser)
	if err != nil {
		return "", err
	}
	return string(result), err
}

// FromString --
func (jwtUser *JWTUser) FromString(value string) error {
	return json.Unmarshal([]byte(value), jwtUser)
}

func generateJWTClaims(subject JWTSubject, expires int64) (jwt.StandardClaims, error) {
	jwtClaims := jwt.StandardClaims{}
	subjectString, err := subject.ToString()
	if err != nil {
		return jwtClaims, err
	}
	now := time.Now().Unix()
	jwtClaims.IssuedAt = now
	jwtClaims.NotBefore = now
	jwtClaims.ExpiresAt = now + expires
	jwtClaims.Id = fmt.Sprintf("%d", subject.ID())
	jwtClaims.Subject = subjectString
	return jwtClaims, nil
}

// GenerateJWTToken 根据载荷生成token
func GenerateJWTToken(subject JWTSubject, expires int64, sceret string) (string, jwt.StandardClaims, error) {
	jwtClaims, err := generateJWTClaims(subject, expires)
	if err != nil {
		return "", jwtClaims, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(sceret))
	if err != nil {
		return "", jwtClaims, err
	}
	return tokenString, jwtClaims, nil
}

func keyFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}

// ParseJWTToken 根据token解析出载荷
func ParseJWTToken(tokenString string, secret string) (jwt.StandardClaims, error) {
	jwtClaims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims, keyFunc(secret))
	if err != nil {
		return jwtClaims, err
	}
	if !token.Valid {
		return jwtClaims, errors.New("cannot convert claim to jwt.StandardClaims")
	}
	return jwtClaims, nil
}

func parseJWTClaims(r *http.Request, secret string) (jwt.StandardClaims, error) {
	tokenString := getTokenString(r)
	return ParseJWTToken(tokenString, secret)
}

func getTokenString(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

func setJWTClaimsToRequest(r *http.Request, jwtClaims jwt.StandardClaims) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), JWTClaimsKey, jwtClaims))
}

// GetJWTClaimsSubjectFromRequest 取JWT用户
func GetJWTClaimsSubjectFromRequest(r *http.Request, subject JWTSubject) error {
	claims := r.Context().Value(JWTClaimsKey).(jwt.StandardClaims)
	return subject.FromString(claims.Subject)
}
