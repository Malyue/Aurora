package redis

const KeyRedisOfflineMsgPrefix = "im.msg.offline"

func (r *RedisClient) SetOfflineMessage(userId string) error {
	key := KeyRedisOfflineMsgPrefix + userId
}

func GetOfflineMessage(userid int64) {

}
