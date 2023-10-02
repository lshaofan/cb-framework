package interfaces

type HttpClient interface {
	// Get get request 想要将返回值解析出来 就需要传入obj的指针 obj 中必须要有response.CommonError{}
	Get(uri string) ([]byte, error)

	// Post post request 想要将返回值解析出来 就需要传入obj的指针 中必须要有response.CommonError{}
	Post(uri string, data []byte, header map[string]string) ([]byte, error)
	PostJSON(uri string, params interface{}) ([]byte, error)
}
