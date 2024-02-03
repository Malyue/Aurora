package message

type Msg struct {
	Type MsgType `json:"type"`
	Msg  Interface
}

type Interface interface {
	GetMsgID() int64
	GetMsgType() MsgType
	GetPayload() []byte
	GetReceiverType() int64
}

type Receiver int
type MsgType string

const (
	Person    MsgType = "person"
	Group     MsgType = "group"
	Broadcast MsgType = "broadcast"
	Ack       MsgType = "ack"
	Auth      MsgType = "auth"
)

type ContentType int

const (
	TEXT = 1
	FILE = 2
)

func IfMsgTypeAllowed(t MsgType) bool {
	if t == Auth || t == Ack || t == Broadcast || t == Group || t == Person {
		return true
	}
	return false
}
