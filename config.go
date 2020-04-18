package account

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
	return true
}
