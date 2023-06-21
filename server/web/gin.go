package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type GinActionImpl struct {
	c   *gin.Context
	res *Response
	req *Request
}

/** =================================response================================= */

func (g *GinActionImpl) returnJsonWithStatusOK() {
	g.c.AbortWithStatusJSON(http.StatusOK, g.res)
}

func (g *GinActionImpl) returnJsonWithStatusBadRequest() {
	g.c.AbortWithStatusJSON(http.StatusBadRequest, g.res)
}

// ThrowError 抛出错误
func (g *GinActionImpl) ThrowError(err *ErrorModel) {

	g.c.AbortWithStatusJSON(err.HttpStatus, NewResponse(
		err.Code,
		err.Message,
		err.Result,
	))
}

// Success 成功
func (g *GinActionImpl) Success(data interface{}) {
	g.res = NewResponse(SUCCESS, Succeed, data)
	g.returnJsonWithStatusOK()
}

// Error 失败
func (g *GinActionImpl) Error(err interface{}) {
	g.res = NewResponse(ERROR, "", nil)
	// 判断err 类型
	switch err.(type) {
	case *ErrorModel:
		g.ThrowError(err.(*ErrorModel))
		return
	case string:
		g.res.Message = err.(string)
	case error:
		g.res.Message = err.(error).Error()

	default:
		g.res.Message = "未知错误"

	}
	g.returnJsonWithStatusBadRequest()
}

// ThrowValidateError 参数验证错误抛出异常
func (g *GinActionImpl) ThrowValidateError(err error) {
	//	判断是否为ErrorModel
	if errModel, ok := err.(*ErrorModel); ok {
		g.ThrowError(errModel)
	} else {
		g.Error(err.Error())
	}

}

// CreateOK 创建成功
func (g *GinActionImpl) CreateOK() {
	g.res = NewResponse(SUCCESS, CreateSuccess, nil)
	g.returnJsonWithStatusOK()
}

// UpdateOK 更新成功
func (g *GinActionImpl) UpdateOK() {
	g.res = NewResponse(SUCCESS, UpdateSuccess, nil)
	g.returnJsonWithStatusOK()
}

// DeleteOK 删除成功
func (g *GinActionImpl) DeleteOK() {
	g.res = NewResponse(SUCCESS, DeleteSuccess, nil)
	g.returnJsonWithStatusOK()
}

/** =================================request================================= */

// BindParam 绑定参数
func (g *GinActionImpl) BindParam(param interface{}) error {
	//	 判断入参是否为指针是否为空
	if param == nil {
		panic("绑定参数不能为空")
	}
	//	 是否是指针
	if reflect.TypeOf(param).Kind() != reflect.Ptr {
		panic("绑定参数必须为指针")
	}
	//	 绑定参数
	err := g.c.ShouldBind(param)
	if err != nil {
		return g.req.GetValidateErr(err, param)
	}

	return nil
}

func NewGinActionImpl(c *gin.Context) *GinActionImpl {
	return &GinActionImpl{
		c:   c,
		req: NewRequest(),
	}
}
