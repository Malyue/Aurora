package message

import (
	"encoding/json"
	"testing"
)

func TestDecode(t *testing.T) {
	s := Msg{
		Type: Auth,
		Msg: &AuthMessage{
			ID: 1,
		},
	}

	buf, err := json.Marshal(s)
	if err != nil {
		t.Error(err)
	}
	msg, err := Decode(buf)
	if err != nil {
		t.Error(err)
	}
	t.Log(msg)
}
