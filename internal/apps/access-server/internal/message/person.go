package message

import (
	_sony "Aurora/internal/apps/access-server/pkg/sonyflake"
	"encoding/json"
	"time"
)

type PersonMessage struct {
	ID          uint64
	Content     []byte
	ContentType ContentType
	ReceiverID  string
	SendID      string
	SendTime    *time.Time
}

func NewPersonMessage(content []byte, contentType ContentType, receiverID string, sendID string, time *time.Time) *PersonMessage {
	return &PersonMessage{
		ID:          _sony.NextID(),
		Content:     content,
		ContentType: contentType,
		ReceiverID:  receiverID,
		SendID:      sendID,
		SendTime:    time,
	}
}

func ParsePersonMessage(data []byte) (*PersonMessage, error) {
	var msg *PersonMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return NewPersonMessage(msg.Content, msg.ContentType, msg.ReceiverID, msg.SendID, msg.SendTime), nil
}

func (p *PersonMessage) GetMsgID() uint64 {
	return p.ID
}

func (p *PersonMessage) GetMsgType() MsgType {
	return Person
}

func (p *PersonMessage) GetPayload() []byte {
	return p.Content
}

func (p *PersonMessage) GetReceiverID() ([]string, bool) {
	return []string{p.ReceiverID}, false
}
