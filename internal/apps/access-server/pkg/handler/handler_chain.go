package handler

import (
	_client "Aurora/internal/apps/access-server/pkg/client"
	_message "Aurora/internal/apps/access-server/pkg/message"
)

type handlerChain struct {
	h    MessageHandler
	next *handlerChain
}

func (hc handlerChain) add(i MessageHandler) {
	if hc.next == nil {
		hc.next = &handlerChain{
			h: i,
		}
	} else {
		hc.next.add(i)
	}
}

// handle use handlerChain
func (hc handlerChain) handle(h *MessageInterfaceImpl, clientInfo *_client.Info, message *_message.Message) bool {
	if hc.h != nil && hc.h.Handle(h, clientInfo, message) {
		return true
	}
	if hc.next != nil {
		return hc.next.handle(h, clientInfo, message)
	}
	return false
}
