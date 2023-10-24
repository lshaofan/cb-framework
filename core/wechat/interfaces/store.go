package interfaces

type Store interface {
	// GetAccessToken 获取AccessToken
	GetAccessToken(key string) (string, error)
	SetAccessToken(key string, accessToken string, expires int64) error

	// GetJsapiTicket 获取JsapiTicket
	GetJsapiTicket(key string) (string, error)
	SetJsapiTicket(key string, jsapiTicket string, expires int64) error
}
