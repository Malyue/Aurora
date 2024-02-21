package handler

import (
	_client "Aurora/internal/apps/access-server/pkg/client"
	_message "Aurora/internal/apps/access-server/pkg/message"
)

func (d *MessageHandlerImpl) dispatch(uid string, m *_message.Message) bool {
	// client id will set in EnqueueMessage
	id := _client.NewID("", uid, "")
	err := d.def.GetClientInterface().EnqueueMessage(id, m)
	var ok = true
	if err != nil {
		d.ctx.Logger.Errorf("dispatch message error : %s", err)
		ok = false
	}
	return ok
}

func (d *MessageHandlerImpl) dispatchOffline(c *_client.Info, message *_message.ChatMessage) error {
	// TODO store offline message
	return nil
}

func (d *MessageHandlerImpl) dispatchOnline(c *_client.Info, message *_message.ChatMessage) error {
	message.From = c.ID.UID()
	dispatchMsg := _message.NewMessage(-1, _message.ActionChatMessage, message)
	return d.def.GetClientInterface().EnqueueMessage(c.ID, dispatchMsg)
}
