package interfaces

type Store interface {
	// GetAccessToken 获取AccessToken
	GetAccessToken(key string) (string, error)
	SetAccessToken(key string, accessToken string, expires int64) error
}
