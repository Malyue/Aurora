package message

import (
	_sonyflake "Aurora/internal/apps/access-server/pkg/sonyflake"
	"encoding/json"
	"errors"
	"time"
)

type AuthMessage struct {
	ID       uint64
	Token    string
	SendTime *time.Time
}

func NewAuthMessage(token string, time *time.Time) *AuthMessage {
	return &AuthMessage{
		ID:       _sonyflake.NextID(),
		Token:    token,
		SendTime: time,
	}
}

func ParseAuthMessage(data []byte) (*AuthMessage, error) {
	var msg *AuthMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return NewAuthMessage(msg.Token, msg.SendTime), nil
}

func (a *AuthMessage) GetMsgID() uint64 {
	return a.ID
}

func (a *AuthMessage) GetMsgType() MsgType {
	return Auth
}

func (a *AuthMessage) GetReceiverID() ([]string, bool) {
	return []string{}, false
}

func (a *AuthMessage) GetPayload() []byte {
	return []byte(a.Token)
}

func HandlerAuthMessage(data []byte) (*AuthMessage, error) {
	msg := &AuthMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	if msg.GetMsgType() != Auth {
		return nil, errors.New("the msg is not auth")
	}

	return msg, err
}
