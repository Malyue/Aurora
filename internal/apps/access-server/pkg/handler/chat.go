package handler

import (
	_client "Aurora/internal/apps/access-server/pkg/client"
	_message "Aurora/internal/apps/access-server/pkg/message"
)

// handleChatMessage
func (d *MessageHandlerImpl) handleChatMessage(c *_client.Info, m *_message.Message) error {
	msg := &_message.ChatMessage{}
	// unmarshal message to chat message
	if !d.unmarshalData(m, msg) {
		return nil
	}

	// get id from client
	msg.From = c.ID.UID()
	msg.To = m.To

	// if id equals 0 means the message is not receive by server
	if msg.Mid == 0 && m.Action != _message.ActionChatMessageResend {
		// TODO store it in mysql

	}

	// send to client ack, it means the server get the msg from client
	err := d.ackChatMessage(c, msg)
	if err != nil {
		d.ctx.Logger.Errorf("ack chat message error : %s", err)
	}

	// init a msg and sends to receiver
	pushMsg := _message.NewMessage(0, _message.ActionChatMessage, msg)
	if !d.dispatch(msg.To, pushMsg) {
		// if send error, it means the client is offline
		// ack notify message
		err := d.ackNotifyMessage(c, msg)
		if err != nil {
			d.ctx.Logger.Errorf("ack notify message error : %s", err)
		}
		return d.dispatchOffline(c, msg)
	}
	return nil
}
