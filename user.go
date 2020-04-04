package account

import (
	"github.com/dgrijalva/jwt-go"
)

// User 用户模型
type User struct {
	Account
	Mobile *string `json:"mobile"`
}

// ResponseKey --
func (user User) ResponseKey() [2]string {
	return [...]string{"user", "users"}
}

// GenerateJWT --
func (user User) GenerateJWT(expires int64, jwtSecret string) (string, jwt.StandardClaims, error) {
	return user.generateJWT(UserTypeNormal, expires, jwtSecret)
}
