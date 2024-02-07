package message

import (
	_sonyflake "Aurora/internal/apps/access-server/pkg/sonyflake"
	"encoding/json"
	"time"
)

type BroadcastMessage struct {
	ID       uint64
	Payload  []byte
	SendTime *time.Time
}

func NewBroadcastMessage(payload []byte, time *time.Time) *BroadcastMessage {
	return &BroadcastMessage{
		ID:       _sonyflake.NextID(),
		Payload:  payload,
		SendTime: time,
	}
}

func ParseBroadcastMessage(data []byte) (*BroadcastMessage, error) {
	var msg *BroadcastMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return NewBroadcastMessage(msg.Payload, msg.SendTime), nil
}

func (b *BroadcastMessage) GetMsgID() uint64 {
	return b.ID
}

func (b *BroadcastMessage) GetMsgType() MsgType {
	return Broadcast
}

func (b *BroadcastMessage) GetPayload() []byte {
	return b.Payload
}

func (b *BroadcastMessage) GetReceiverID() ([]string, bool) {
	return []string{}, true
}
