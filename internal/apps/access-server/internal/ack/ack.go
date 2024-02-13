package ack

import (
	"Aurora/internal/apps/access-server/pkg/message"
)

// retain an ack set
// if it gets the response from the client, it means the msg has sent success
// otherwise, the msg is lost or the ack is lost
// set a seq number to remove duplicate msg in the client
// and if it is out of the ticker, resend the msg

type AckQueue struct {
	AckMsgs []*message.Message
}

func (a *AckQueue) Ack() {

}
