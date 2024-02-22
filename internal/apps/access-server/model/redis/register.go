package redis

import (
	"context"
)

func (r *RedisClient) SetRegisterRouterInfo(userid string, device string, gate string) error {
	key := KeyClientRoutePrefix + userid + ":" + device
	return r.Set(context.Background(), key, gate)
}

func (r *RedisClient) GetRegisterInfo(key string) (string, error) {
	return r.Get(context.Background(), KeyClientRoutePrefix+key)
}

func (r *RedisClient) DelRegisterInfo(key string) error {
	_, err := r.Client.Del(context.Background(), KeyClientRoutePrefix+key).Result()
	return err
}
