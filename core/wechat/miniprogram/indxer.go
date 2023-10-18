package miniprogram

import (
	"encoding/json"
	"errors"
	"fmt"
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

func WithCachePrefix(prefix string) Options {
	return func(c *Client) {
		c.CachePrefix = prefix
	}
}

func WithAppidAndSecret(appId string, appSecret string) Options {
	if appId == "" || appSecret == "" {
		panic("appid or app_secret is nil")
	}
	return func(c *Client) {
		c.AppId = appId
		c.AppSecret = appSecret
	}
}

type Client struct {
	HttpClient        interfaces.HttpClient
	Store             interfaces.Store
	AppId             string `json:"app_id" yaml:"app_id"`
	AppSecret         string `json:"app_secret" yaml:"app_secret"`
	CachePrefix       string `json:"cache_prefix" yaml:"cache_prefix"`
	AccessToken       string `json:"access_token" yaml:"access_token"`
	AccessTokenLock   *sync.Mutex
	refreshTokenCount int
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

func NewClient(o ...Options) *Client {
	c := &Client{
		AccessTokenLock: new(sync.Mutex),
	}
	for _, opt := range o {
		opt(c)
	}
	// 设置缓存前缀
	c.CachePrefix = fmt.Sprintf("%s%s:", constants.MiniProgramCacheKeyPrefix, c.CachePrefix)
	if c.HttpClient == nil {
		c.HttpClient = NewHttpClient()
	}
	return c
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
	commonError := responseObj.Elem().FieldByName("CommonError")
	if !commonError.IsValid() || commonError.Kind() != reflect.Struct {
		return StructNotHasCommonError
	}
	errCode := commonError.FieldByName("ErrCode")
	errMsg := commonError.FieldByName("ErrMsg")
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
	accessToken, err := c.Store.GetAccessToken(fmt.Sprintf("%s%s", c.GetCachePrefix(), c.AppId))
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
		err = c.Store.SetAccessToken(fmt.Sprintf("%s%s", c.GetCachePrefix(), c.AppId), ac.AccessToken, ac.ExpiresIn-1500)
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
	err = c.Store.SetAccessToken(fmt.Sprintf("%s%s", c.GetCachePrefix(), c.AppId), ac.AccessToken, ac.ExpiresIn-1500)
	if err != nil {
		return "", err
	}
	return ac.AccessToken, nil
}

// GetAccessTokenResult 获取access_token结果
type GetAccessTokenResult struct {
	response.CommonError
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// 从微信服务器获取access_token
func (c *Client) getAccessTokenFromServer() (*GetAccessTokenResult, error) {
	ret := new(GetAccessTokenResult)

	res, err := c.HttpClient.Get(fmt.Sprintf(
		constants.MiniProgramAccessTokenURL,
		c.AppId,
		c.AppSecret,
	))

	if err != nil {
		return nil, err
	}
	err = c.handleResp(res, ret, constants.MiniProgramAccessTokenURL)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Code2SessionResult 获取用户的openid和session_key 的结果
type Code2SessionResult struct {
	response.CommonError
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
}

// Code2Session 登陆凭证校验的结果
func (c *Client) Code2Session(code string) (*Code2SessionResult, error) {
	ret := new(Code2SessionResult)
	res, err := c.HttpClient.Get(fmt.Sprintf(
		constants.MiniProgramCode2SessionURL,
		c.AppId,
		c.AppSecret,
		code,
	))
	if err != nil {
		return nil, err
	}
	err = c.handleResp(res, ret, constants.MiniProgramCode2SessionURL)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
