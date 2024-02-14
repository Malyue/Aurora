package client

import (
	_message "Aurora/internal/apps/access-server/pkg/message"
)

var (
	ClientTypeRobot = 1
	ClientTypeUser  = 2
)

type Client interface {
	SetID(id ID)
	IsRunning() bool
	EnqueueMessage(message *_message.Message) error
	Exit()
	Run()
	GetInfo() Info
	AddMessageInterceptor(interceptor MessageInterceptor)
}

//
//func NewID() string {
//	return _sonyflake.NextStringID()
//}
