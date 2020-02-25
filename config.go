package account

import "log"

// ConfigJWT --
type ConfigJWT struct {
	Expires int64
	Secret  string
}

// Check --
func (config *ConfigJWT) Check(cmd string) bool {
	if cmd == "all" || cmd == "write" {
		if config.Expires <= 0 {
			return false
		}
	}
	if cmd == "all" || cmd == "read" {
		if config.Secret == "" {
			return false
		}
	}
	log.Println("config jwt check valid")
	return true
}

// ConfigMiniProgram --
type ConfigMiniProgram struct {
	AppID     string
	AppSecret string
}

// Check --
func (config *ConfigMiniProgram) Check() bool {
	if config.AppID == "" {
		return false
	}
	if config.AppSecret == "" {
		return false
	}
	log.Println("config miniprogram check valid")
	return true
}
