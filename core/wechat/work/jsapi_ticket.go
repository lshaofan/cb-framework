package work

import (
	"fmt"
	"github.com/lshaofan/cb-framework/core/wechat/constants"
	"github.com/redis/go-redis/v9"
	"sync"
)

type JsapiTicketResponse struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// GetJsapiTicket 获取企业的jsapi_ticket
func (c *Client) GetJsapiTicket() (string, error) {

	// 先判断是否有值
	if c.AccessToken != "" {
		return c.AccessToken, nil
	}
	var ticket string
	var err error

	// 从缓存中获取
	ticket, err = c.Store.GetAccessToken(fmt.Sprintf("%sticket:%s", c.GetCachePrefix(), c.CorpID))
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == redis.Nil {
		// 加锁
		m := sync.Mutex{}

		m.Lock()
		defer m.Unlock()
		ret := new(JsapiTicketResponse)
		_, err = c.Get(fmt.Sprintf(
			constants.WorkJsapiTicketURL,
		), ret)
		if err != nil {
			return "", err
		}
		// 设置 cache
		err = c.Store.SetJsapiTicket(fmt.Sprintf("%sticket:%s", c.GetCachePrefix(), c.CorpID), ret.Ticket, int64(ret.ExpiresIn-1500))
		if err != nil {
			return "", err
		}
		ticket = ret.Ticket
	}
	return ticket, nil
}
