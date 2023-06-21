package web

type Action interface {
	// Success 成功并返回数据
	Success(data interface{})
	Error(err interface{})
	ThrowError(err *ErrorModel)
	ThrowValidateError(err error)
	BindParam(param interface{}) error
	CreateOK()
	UpdateOK()
	DeleteOK()
}
