package handler

import (
	_client "Aurora/internal/apps/access-server/pkg/client"
	_message "Aurora/internal/apps/access-server/pkg/message"
)

type ActionHandler struct {
	action _message.Action
	fn     HandlerFunc
}

func NewActionHandler(action _message.Action, fn HandlerFunc) *ActionHandler {
	return &ActionHandler{
		action: action,
		fn:     fn,
	}
}

func (a *ActionHandler) Handle(h *MessageInterfaceImpl, clientInfo *_client.Info, message *_message.Message) bool {
	// check if the action is true, otherwise return false directly
	if message.GetAction() == a.action {
		err := a.fn(clientInfo, message)
		if err != nil {
			h.onHandleMessageError(clientInfo, message, err)
		}
		return true
	}
	return false
}
