package account

import "golang.org/x/crypto/bcrypt"

// BcryptGenerateFromPassword 加密密码
func BcryptGenerateFromPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// BcryptCompareHashAndPassword 比较加密密码和明文密码
func BcryptCompareHashAndPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
