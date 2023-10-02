package response

import "fmt"

// CommonError 微信返回的通用错误json
type CommonError struct {
	ApiName string
	ErrMsg  string `json:"errmsg"`
	ErrCode int64  `json:"errcode"`
}

func (e *CommonError) Error() string {
	return fmt.Sprintf("请求出错: errcode=%d ; errmsg=%s ; api=%s", e.ErrCode, e.ErrMsg, e.ApiName)
}

func (e *CommonError) GetErrCode() int64 {
	return e.ErrCode
}

func (e *CommonError) GetErrMsg() string {
	return e.ErrMsg
}

func (e *CommonError) GetApiName() string {
	return e.ApiName
}
