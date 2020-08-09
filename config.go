package weixin

import "log"

// ConfigMiniProgram --
type ConfigMiniProgram struct {
	AppID     string
	AppSecret string
}

// Check --
func (config ConfigMiniProgram) Check() bool {
	if config.AppID == "" {
		return false
	}
	if config.AppSecret == "" {
		return false
	}
	log.Println("ConfigMiniProgram check pass")
	return true
}
