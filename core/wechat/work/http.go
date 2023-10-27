package work

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lshaofan/cb-framework/server/web"
	"io"
	"net/http"
)

var (
	// ObjDecodeErr 返回值对象解析错误
	ObjDecodeErr = web.NewErrorModel(-1, "json解析失败", nil, http.StatusBadRequest)
	// StructNotHasCommonError 请求需要解析的对象错误  没有commonError
	StructNotHasCommonError = web.NewErrorModel(-1, "需要解析的结构体错误", nil, http.StatusBadRequest)
)

type HttpCliOptions func(*HttpClient)

func WithHttpClientCtx(c context.Context) HttpCliOptions {
	return func(client *HttpClient) {
		client.ctx = c
	}
}

func WithHttpClientDebug(debug bool) HttpCliOptions {
	return func(client *HttpClient) {
		client.debug = debug
	}
}

type HttpClient struct {
	ctx   context.Context
	debug bool
}

func (h *HttpClient) PostJSON(uri string, params interface{}) ([]byte, error) {
	jsonBuf := new(bytes.Buffer)
	enc := json.NewEncoder(jsonBuf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(params)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(uri, "application/json;charset=utf-8", jsonBuf)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return io.ReadAll(response.Body)
}

func (h *HttpClient) Post(uri string, data []byte, header map[string]string) ([]byte, error) {
	body := bytes.NewBuffer(data)
	request, err := http.NewRequestWithContext(h.ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}
	for key, value := range header {
		request.Header.Set(key, value)
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, web.NewErrorModel(
			-1,
			fmt.Sprintf("微信服务器get请求出错: uri=%v , statusCode=%v", uri, res.StatusCode),
			res,
			res.StatusCode,
		)
	}
	return io.ReadAll(res.Body)
}

func (h *HttpClient) Get(uri string) ([]byte, error) {
	// debug
	// debug
	if h.debug {
		fmt.Println("get请求的url:", uri)
	}
	request, err := http.NewRequestWithContext(h.ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, web.NewErrorModel(
			-1,
			fmt.Sprintf("微信服务器get请求出错: uri=%v , statusCode=%v", uri, res.StatusCode),
			res,
			res.StatusCode,
		)
	}
	resp, err := io.ReadAll(res.Body)
	// debug
	if h.debug {
		fmt.Println("get请求的resp:", string(resp))
	}
	return resp, err
}

func NewHttpClient(opts ...HttpCliOptions) *HttpClient {
	h := new(HttpClient)
	for _, o := range opts {
		o(h)

	}
	if h.ctx == nil {
		h.ctx = context.Background()
	}
	return h

}
