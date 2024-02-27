package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var Ver1 int64 = 1

type Message struct {
	Ver    int64  `json:"ver,omitempty"`
	Seq    int64  `json:"seq,omitempty"`
	ID     int64  `json:"id,omitempty"`
	Action string `json:"action"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Data   *Data  `json:"data,omitempty"`
	Msg    string `json:"msg,omitempty"`

	Ticket string `json:"ticket,omitempty"`
	Sign   string `json:"sign,omitempty"`

	Extra map[string]string `json:"extra,omitempty"`
}

func NewMessage(seq int64, action Action, data interface{}) *Message {
	return &Message{
		Ver:    Ver1,
		Seq:    seq,
		Action: string(action),
	}
}

func NewEmptyMessage() *Message {
	return &Message{
		Ver:   Ver1,
		Data:  nil,
		Extra: nil,
	}
}

func (m *Message) GetSeq() int64 {
	return m.Seq
}

func (m *Message) GetAction() Action {
	return Action(m.Action)
}

func (m *Message) SetSeq(seq int64) {
	m.Seq = seq
}

func (m *Message) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("&Message{Ver:%d, Action:%s, Data:%s}", m.Ver, m.Action, m.Data)
}

// Data used to wrap message data
// server received a msg, the data type is []byte, it's waiting for deserialize to specified struct
// When server push a msg to client, the data type is specific struct
type Data struct {
	des interface{}
}

func NewData(d interface{}) *Data {
	return &Data{
		des: d,
	}
}

func (d *Data) UnmarshalJSON(bytes []byte) error {
	d.des = bytes
	return nil
}

func (d *Data) MarshalJSON() ([]byte, error) {
	bytes, ok := d.des.([]byte)
	if ok {
		return bytes, nil
	}
	return JsonCodec.Encode(d.des)
}

func (d *Data) GetData() interface{} {
	return d.des
}

// Deserialize attempts to decode the `des` field in a `Data` object into the provided interface{}
// returns an error if the operation fails
func (d *Data) Deserialize(i interface{}) error {
	if d == nil {
		return errors.New("data is nil")
	}
	s, ok := d.des.([]byte)
	if ok {
		return JsonCodec.Decode(s, i)
	} else {
		t1 := reflect.TypeOf(i)
		t2 := reflect.TypeOf(d.des)
		if t1 == t2 {
			reflect.ValueOf(i).Elem().Set(reflect.ValueOf(d.des).Elem())
			return nil
		}
	}
	return errors.New("deserialize message data failed")
}

func (d *Data) String() string {
	b, ok := d.des.([]byte)
	var s interface{}
	if ok {
		s = string(b)
	} else {
		if d.des == nil {
			s = "<nil>"
		} else {
			s, _ = json.Marshal(d.des)
		}
	}
	return fmt.Sprintf("%s", s)
}
