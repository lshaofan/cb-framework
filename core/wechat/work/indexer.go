package work

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lshaofan/cb-framework/core/wechat/constants"
	"github.com/lshaofan/cb-framework/core/wechat/interfaces"
	"github.com/lshaofan/cb-framework/core/wechat/response"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strings"
	"sync"
)

type Options func(*Client)

func WithHttpClient(cli interfaces.HttpClient) Options {
	if cli == nil {
		panic("http client is nil")
	}
	return func(c *Client) {
		c.HttpClient = cli
	}
}

func WithRedisStore(cli *redis.Client) Options {
	if cli == nil {
		panic("redis client is nil")
	}
	return func(c *Client) {
		c.Store = NewRedisStore(cli)
	}
}

func WithAccessToken(accessToken string) Options {
	return func(c *Client) {
		c.AccessToken = accessToken
	}
}

// WithCachePrefix 设置缓存前缀，必填
func WithCachePrefix(prefix string) Options {
	return func(c *Client) {
		c.CachePrefix = prefix
	}
}

func WithAppidAndSecret(CorpID string, CorpSecret string) Options {
	// 判断appid和app_secret是否为空
	if CorpID == "" || CorpSecret == "" {
		panic("appid or app_secret is nil")
	}

	return func(c *Client) {
		c.CorpID = CorpID
		c.CorpSecret = CorpSecret
	}
}

// WithDebug 是否开启debug模式
func WithDebug(debug bool) Options {
	return func(c *Client) {
		c.Debug = debug
	}
}

type AccessTokenType string

const (
	// ContactAccessTokenType 通讯录
	ContactAccessTokenType AccessTokenType = "contact_secret"
	// ExternalContactAccessTokenType 外部联系人
	ExternalContactAccessTokenType AccessTokenType = "external_contact_secret"
)

type Client struct {
	ctx               *gin.Context
	CorpID            string `json:"corp_id"`
	CorpSecret        string `json:"corp_secret"`
	CachePrefix       string `json:"cache_prefix" yaml:"cache_prefix"`
	AccessToken       string `json:"access_token" yaml:"access_token"`
	AccessTokenLock   *sync.Mutex
	refreshTokenCount int
	HttpClient        interfaces.HttpClient
	Store             interfaces.Store
	Debug             bool
}

func NewClient(o ...Options) *Client {
	c := &Client{
		AccessTokenLock: new(sync.Mutex),
	}
	for _, opt := range o {
		opt(c)
	}
	// 设置缓存前缀
	if c.CachePrefix == "" {
		panic("cache_prefix 缓存前缀必填")
	}
	c.CachePrefix = fmt.Sprintf("%s%s:", constants.WorkCacheKeyPrefix, c.CachePrefix)
	if c.HttpClient == nil {
		c.HttpClient = NewHttpClient(
			WithHttpClientDebug(c.Debug),
		)
	}
	return c
}

// HasQuery 判断url中是否有参数
func (c *Client) HasQuery(url string) bool {
	return strings.Contains(url, "?")
}
func (c *Client) Get(url string, obj interface{}) ([]byte, error) {
	// 判断参数中是否有baseApi
	ak, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
Request:
	// 判断url中是否有参数
	if c.HasQuery(url) {
		url = fmt.Sprintf("%s&access_token=%s", url, ak)
	} else {
		url = fmt.Sprintf("%s?access_token=%s", url, ak)
	}

	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}

	// 没有传入obj则直接返回
	if obj == nil {
		return resp, nil
	}
	err = c.handleResp(resp, obj, url)
	if err != nil {
		// 判断错误是否是response.CommonError
		var e *response.CommonError
		if errors.As(err, &e) {
			// 判断是否是40001
			if e.GetErrCode() == constants.AccessToken40001 {
				// 判断刷新次数是否超过3次
				if c.refreshTokenCount > 3 {
					return resp, e
				}
				// 则刷新 access_token
				ak, err = c.RefreshAccessToken()
				if err != nil {
					return nil, err
				}
				// 增加刷新次数
				c.refreshTokenCount++
				goto Request
			}

		}
	}
	return resp, err
}

func (c *Client) Post(url string, body []byte, obj interface{}, header map[string]string) ([]byte, error) {
	// 判断参数中是否有baseApi
	ak, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
Request:

	// 判断url中是否有参数
	if c.HasQuery(url) {
		url = fmt.Sprintf("%s&access_token=%s", url, ak)
	} else {
		url = fmt.Sprintf("%s?access_token=%s", url, ak)
	}

	resp, err := c.HttpClient.Post(url, body, header)
	if err != nil {
		return nil, err
	}
	// 没有传入obj则直接返回
	if obj == nil {
		return resp, nil
	}

	err = c.handleResp(resp, obj, url)
	if err != nil {
		// 判断错误是否是response.CommonError
		var e *response.CommonError
		if errors.As(err, &e) {
			// 判断是否是40001
			if e.GetErrCode() == constants.AccessToken40001 {
				// 判断刷新次数是否超过3次
				if c.refreshTokenCount > 3 {
					return resp, e
				}
				// 则刷新 access_token
				ak, err = c.RefreshAccessToken()
				if err != nil {
					return nil, err
				}
				// 增加刷新次数
				c.refreshTokenCount++
				goto Request
			}

		}
	}
	return resp, err
}

func (c *Client) PostJSON(url string, params interface{}, obj interface{}) ([]byte, error) {
	// 判断参数中是否有baseApi
	ak, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// 打一个标签 用于判断是否是刷新access_token 刷了之后需要重新请求
Request:
	// 判断url中是否有参数
	if c.HasQuery(url) {
		url = fmt.Sprintf("%s&access_token=%s", url, ak)
	} else {
		url = fmt.Sprintf("%s?access_token=%s", url, ak)
	}
	resp, err := c.HttpClient.PostJSON(url, params)
	if err != nil {
		return nil, err
	}

	// 没有传入obj则直接返回
	if obj == nil {
		return resp, nil
	}

	err = c.handleResp(resp, obj, url)
	if err != nil {
		// 判断错误是否是response.CommonError
		var e *response.CommonError
		if errors.As(err, &e) {
			// 判断是否是40001
			if e.GetErrCode() == constants.AccessToken40001 {
				// 判断刷新次数是否超过3次
				if c.refreshTokenCount > 3 {
					return resp, e
				}
				// 则刷新 access_token
				ak, err = c.RefreshAccessToken()
				if err != nil {
					return nil, err
				}
				// 增加刷新次数
				c.refreshTokenCount++
				goto Request
			}

		}
	}
	return resp, err
}

// 处理server响应的方法
func (c *Client) handleResp(ret []byte, obj interface{}, apiName string) error {
	err := json.Unmarshal(ret, obj)
	if err != nil {
		return ObjDecodeErr
	}
	// 判断是否传入了obj
	responseObj := reflect.ValueOf(obj)
	if !responseObj.IsValid() {
		return StructNotHasCommonError
	}
	//commonError := responseObj.Elem().FieldByName("CommonError")
	//if !commonError.IsValid() || commonError.Kind() != reflect.Struct {
	//	return StructNotHasCommonError
	//}
	errCode := responseObj.Elem().FieldByName("ErrCode")
	errMsg := responseObj.Elem().FieldByName("ErrMsg")
	if !errCode.IsValid() || !errMsg.IsValid() {
		return StructNotHasCommonError
	}
	if errCode.Int() != 0 {
		e := new(response.CommonError)
		e.ErrMsg = errMsg.String()
		e.ErrCode = errCode.Int()
		e.ApiName = apiName
		return e
	}
	return nil
}

// GetCachePrefix 获取缓存前缀
func (c *Client) GetCachePrefix() string {
	return c.CachePrefix
}

// GetAccessToken 获取access_token
func (c *Client) GetAccessToken() (string, error) {

	// 先判断是否有值
	if c.AccessToken != "" {
		return c.AccessToken, nil
	}
	// 从缓存中获取
	accessToken, err := c.Store.GetAccessToken(fmt.Sprintf("%s%s", c.GetCachePrefix(), c.CorpID))
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == redis.Nil {
		// 获取access_token
		// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从微信服务器上获取到不同token
		c.AccessTokenLock.Lock()
		defer c.AccessTokenLock.Unlock()
		ac, err := c.getAccessTokenFromServer()
		if err != nil {
			return "", err
		}

		// 设置cache
		err = c.Store.SetAccessToken(fmt.Sprintf("%s%s", c.GetCachePrefix(), c.CorpID), ac.AccessToken, ac.ExpiresIn-1500)
		if err != nil {
			return "", err
		}
		return ac.AccessToken, nil
	}

	return accessToken, nil
}

// RefreshAccessToken 刷新access_token
func (c *Client) RefreshAccessToken() (string, error) {
	// 获取access_token
	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从微信服务器上获取到不同token
	c.AccessTokenLock.Lock()
	defer c.AccessTokenLock.Unlock()
	ac, err := c.getAccessTokenFromServer()
	if err != nil {
		return "", err
	}

	// 设置cache
	err = c.Store.SetAccessToken(fmt.Sprintf("%s%s", c.GetCachePrefix(), c.CorpID), ac.AccessToken, ac.ExpiresIn-1500)
	if err != nil {
		return "", err
	}
	return ac.AccessToken, nil
}

// GetAccessTokenResult 获取access_token结果
type GetAccessTokenResult struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// 从微信服务器获取access_token
func (c *Client) getAccessTokenFromServer() (*GetAccessTokenResult, error) {
	ret := new(GetAccessTokenResult)

	res, err := c.HttpClient.Get(fmt.Sprintf(
		constants.WorkAccessTokenURL,
		c.CorpID,
		c.CorpSecret,
	))

	if err != nil {
		return nil, err
	}
	err = c.handleResp(res, ret, constants.WorkAccessTokenURL)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
