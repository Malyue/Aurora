package handler

import (
	_client "Aurora/internal/apps/access-server/pkg/client"
	_message "Aurora/internal/apps/access-server/pkg/message"
)

const KeyRedisOfflineMsgPrefix = "im.msg.offline"

// handleAckRequest get ack from receiver and send ack response to sender
func (d *MessageHandlerImpl) handleAckRequest(clientInfo *_client.Info, message *_message.Message) error {
	ackMsg := &_message.AckRequest{}
	if !d.unmarshalData(message, ackMsg) {
		return nil
	}

	ackNotify := _message.NewMessage(0, _message.ActionAckNotify, ackMsg)

	// send to sender
	d.dispatch(ackMsg.From, ackNotify)
	return nil
}

func (d *MessageHandlerImpl) handleAckOffline(c *_client.Info, message *_message.Message) error {
	key := KeyRedisOfflineMsgPrefix + c.ID.UID()
	//result,err := db.Redis.Del(key).Result()

}

func (d *MessageHandlerImpl) ackChatMessage(c *_client.Info, message *_message.ChatMessage) error {
	ackMsg := _message.AckMessage{
		CliMid: message.CliMid,
		Mid:    message.Mid,
		Seq:    0,
	}
	ack := _message.NewMessage(0, _message.ActionAckMessage, &ackMsg)
	// send message
	return d.def.GetClientInterface().EnqueueMessage(c.ID, ack)
}

func (d *MessageHandlerImpl) ackNotifyMessage(c *_client.Info, message *_message.ChatMessage) error {
	ackNotify := _message.AckNotify{Mid: message.Mid}
	msg := _message.NewMessage(0, _message.ActionAckNotify, &ackNotify)
	return d.def.GetClientInterface().EnqueueMessage(c.ID, msg)
}
