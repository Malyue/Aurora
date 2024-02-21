package handler

import (
	"Aurora/internal/apps/access-server/pkg/client"
	_client "Aurora/internal/apps/access-server/pkg/client"
	_message "Aurora/internal/apps/access-server/pkg/message"
	"Aurora/internal/apps/access-server/pkg/subscription"
)

type HandlerFunc func(cliInfo *client.Info, message *_message.Message) error

type MessageHandler interface {
	// Handle handles the message, returns true if the message is handled
	// otherwise the message is delegated to next offlineMessageHandler
	Handle(h *MessageInterfaceImpl, clientInfo *_client.Info, message *_message.Message) bool
}

type Handler interface {
	Handle(clientInfo *client.Info, msg *_message.Message) error

	AddHandler(i MessageHandler)

	SetSubscription(g subscription.Interface)
}
