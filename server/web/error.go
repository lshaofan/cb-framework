package web

import "net/http"

var (
	ServerError = NewErrorModel(500, "服务器错误", nil, http.StatusInternalServerError)
	// UNAUTHORIZED 未登录
	UNAUTHORIZED = NewErrorModel(10000, "未登录", nil, http.StatusUnauthorized)
	// InvalidToken 非法Token
	InvalidToken = NewErrorModel(10001, "非法Token", nil, http.StatusUnauthorized)
	// TokenExpired Token过期
	TokenExpired = NewErrorModel(10002, "Token过期", nil, http.StatusUnauthorized)
	// UsernameOrPasswordError 用户名或密码错误
	UsernameOrPasswordError = NewErrorModel(10003, "用户名或密码错误", nil, http.StatusUnauthorized)
	// PlatformNotExist 平台不存在
	PlatformNotExist = NewErrorModel(10004, "平台不存在", nil, http.StatusPreconditionFailed)
	// PlatformIdCanNotEmpty 平台id不能为空
	PlatformIdCanNotEmpty = NewErrorModel(10005, "平台id不能为空", nil, http.StatusPreconditionFailed)
)

// ErrorModel 错误模型
type ErrorModel struct {
	Code       int         `json:"code" `
	Message    string      `json:"message" `
	Result     interface{} `json:"result"`
	HttpStatus int         `json:"httpStatus" swaggerignore:"true"`
}

func NewErrorModel(code int, message string, result interface{}, httpStatus int) *ErrorModel {
	return &ErrorModel{Code: code, Message: message, Result: result, HttpStatus: httpStatus}
}

func (e *ErrorModel) Error() string {
	return e.Message
}
