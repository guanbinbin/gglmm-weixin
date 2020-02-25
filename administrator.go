package account

import (
	"github.com/dgrijalva/jwt-go"
)

// Administrator 用户模型
type Administrator struct {
	Account
	Mobile string `json:"mobile"`
}

// ResponseKey --
func (administrator Administrator) ResponseKey() [2]string {
	return [...]string{"administrator", "administrators"}
}

// GenerateJWT --
func (administrator Administrator) GenerateJWT(expires int64, jwtSecret string) (string, jwt.StandardClaims, error) {
	return administrator.generateJWT(UserTypeAdministrator, expires, jwtSecret)
}
