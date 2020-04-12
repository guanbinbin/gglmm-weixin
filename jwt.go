package account

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// JWTExpires JWT失效时间
	JWTExpires int64 = 24 * 60 * 60
)

func generateJWTClaims(subject interface{}, expires int64) (*jwt.StandardClaims, error) {
	subjectBytes, err := json.Marshal(subject)
	if err != nil {
		return nil, err
	}
	jwtClaims := &jwt.StandardClaims{}
	now := time.Now().Unix()
	jwtClaims.IssuedAt = now
	jwtClaims.NotBefore = now
	jwtClaims.ExpiresAt = now + expires
	jwtClaims.Subject = string(subjectBytes)
	return jwtClaims, nil
}

func parseJWTClaims(tokenString string, secret string) (*jwt.StandardClaims, error) {
	jwtClaims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("cannot convert claim to jwt.StandardClaims")
	}
	return jwtClaims, nil
}
