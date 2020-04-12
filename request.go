package account

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
