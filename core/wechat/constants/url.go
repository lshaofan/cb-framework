package constants

// Work 企业微信
const (
	// WorkAccessTokenURL 获取access_token的url
	WorkAccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)

// MiniProgram 小程序
const (
	// MiniProgramAccessTokenURL  获取access_token的url
	MiniProgramAccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	// MiniProgramCode2SessionURL 登录凭证校验的url
	MiniProgramCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)
