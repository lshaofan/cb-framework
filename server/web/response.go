package web

// Response  返回数据
type Response struct {
	Code    int         `json:"code" `
	Result  interface{} `json:"result"`
	Message string      `json:"message" `
}

func NewResponse(code int, message string, result interface{}) *Response {
	return &Response{Code: code, Result: result, Message: message}
}

// PageList  分页数据
type PageList[T interface{}] struct {
	Total    int64 `json:"total" `
	Data     []T   `json:"data" `
	Page     int   `json:"page" `
	PageSize int   `json:"page_size" `
}

func NewPageList[T interface{}]() *PageList[T] {
	return &PageList[T]{}
}
