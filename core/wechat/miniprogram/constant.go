package miniprogram

const (
	// AccessToken40001 40001错误
	AccessToken40001 = 40001
	// CacheKeyPrefix 缓存key的前缀
	CacheKeyPrefix = "wechat:miniprogram:"
	// AccessTokenURL  获取access_token的url
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	// Code2SessionURL 登录凭证校验的url
	Code2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)
