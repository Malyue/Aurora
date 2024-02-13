package client

import (
	_message "Aurora/internal/apps/access-server/pkg/message"
	_sonyflake "Aurora/internal/apps/access-server/pkg/sonyflake"
)

var (
	ClientTypeRobot = 1
	ClientTypeUser  = 2
)

type Client interface {
	SetUserID(id string)
	IsRunning() bool
	EnqueueMessage(message *_message.Message) error
	Exit()
	Run()
	GetInfo() Info
	AddMessageInterceptor(interceptor MessageInterceptor)
}

func NewID() string {
	return _sonyflake.NextStringID()
}
