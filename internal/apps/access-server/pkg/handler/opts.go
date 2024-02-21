package handler

import "github.com/go-redis/redis/v8"

type MessageHandlerOptions struct {
	// MessageStore chat message store
	//MessageStore store.MessageStore
	// TODO redis 连接
	RedisCli redis.Client

	// DontInitDefaultHandler true will not init default action offlineMessageHandler, MessageHandlerImpl.InitDefaultHandler
	//DontInitDefaultHandler bool

	// NotifyOnErr true express notify client on server error.
	NotifyOnErr bool
}
