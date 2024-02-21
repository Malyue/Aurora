package handler

import (
	_client "Aurora/internal/apps/access-server/pkg/client"
	_message "Aurora/internal/apps/access-server/pkg/message"
)

func (d *MessageHandlerImpl) handleHeartbeat(clientInfo *_client.Info, message *_message.Message) error {
	return nil
}
