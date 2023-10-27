package miniprogram

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStore struct {
	Client *redis.Client
	ctx    context.Context
}

func (r *RedisStore) GetJsapiTicket(key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisStore) SetJsapiTicket(key string, jsapiTicket string, expires int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisStore) SetAccessToken(key string, accessToken string, expires int64) error {
	return r.Client.Set(r.ctx, key, accessToken,
		time.Duration(expires)*time.Second,
	).Err()
}

func (r *RedisStore) GetAccessToken(key string) (string, error) {
	return r.Client.Get(r.ctx, key).Result()
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		Client: client,
		ctx:    context.Background(),
	}
}
