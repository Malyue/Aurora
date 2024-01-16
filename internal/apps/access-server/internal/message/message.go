package message

import (
	"encoding/json"
	"errors"
)

type Receiver int
type MsgType string

const (
	Person    MsgType = "person"
	Group     MsgType = "group"
	Broadcast MsgType = "broadcast"
	Ack       MsgType = "ack"
	Auth      MsgType = "auth"
)

type ConntentType int

const (
	TEXT = 1
	FILE = 2
)

type Message struct {
	MessageID int64
	// MsgType ack/auth/broadcast/group/person
	MsgType MsgType
	// ContentType text/file/video/sound
	ContentType ConntentType
	ReceiverID  []string
	GroupID     []string
	Subscribers []string
	// Payload content
	Payload []byte

	// -------- queue --------
	private    int64
	index      int
	retryCount int
}

type AuthMessage struct {
	ID    int64
	Token string
}

// TODO add a sequence id such as snowflake to keep the seq id orderly

func MessageAck() {

}

func MessageSync() {

}

func MessageBroadcast() {

}

func GetAuthMessage(data []byte) (*AuthMessage, error) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	if msg.MsgType != Auth {
		return nil, errors.New("the msg is not auth")
	}

	authMessage := &AuthMessage{
		ID:    msg.MessageID,
		Token: string(msg.Payload),
	}

	return authMessage, err
}
