package web

import "github.com/gin-gonic/gin/binding"

type BindOption func(obj any) error

type Action interface {
	// Success 成功并返回数据
	Success(data interface{})
	Error(err interface{})
	ThrowError(err *ErrorModel)
	ThrowValidateError(err error)
	Bind(param any, opts ...BindOption) error
	BindParam(param interface{}) error
	BindUriParam(param interface{}) error
	ShouldBindBodyWith(param any, bb binding.BindingBody) error
	ShouldBindWith(param any, bb binding.Binding) error
	CreateOK()
	UpdateOK()
	DeleteOK()
	SuccessWithMessage(message string, data interface{})
	CreateOkWithMessage(message string)
	UpdateOkWithMessage(message string)
	DeleteOkWithMessage(message string)
}
