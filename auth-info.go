package account

// AuthenticationInfo --
type AuthenticationInfo struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarURL"`
}
