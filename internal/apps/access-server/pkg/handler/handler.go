package handler

import (
	_client "Aurora/internal/apps/access-server/pkg/client"
	_message "Aurora/internal/apps/access-server/pkg/message"
	"Aurora/internal/apps/access-server/pkg/subscription"
	"Aurora/internal/apps/access-server/svc"
)

type MessageHandlerImpl struct {
	def *MessageInterfaceImpl
	// TODO get message model
	//store store.MessageStore
	//userState *UserState
	ctx *svc.ServerCtx
}

func NewHandlerWithOptions(gateway _client.Gateway, ctx *svc.ServerCtx, opts *MessageHandlerOptions) (*MessageHandlerImpl, error) {
	impl, err := NewDefaultImpl(ctx, &Options{
		NotifyServerError:     true,
		MaxMessageConcurrency: 10_0000,
	})
	if err != nil {
		return nil, err
	}
	impl.SetNotifyErrorOnServer(opts.NotifyOnErr)
	ret := &MessageHandlerImpl{
		def: impl,
		// client
	}
	ret.InitDefaultHandler(nil)
	//if !opts.DontInitDefaultHandler {
	//	ret.InitDefaultHandler(nil)
	//}
	return ret, nil
}

func (d *MessageHandlerImpl) InitDefaultHandler(callback func(action _message.Action, fn HandlerFunc) HandlerFunc) {
	m := map[_message.Action]HandlerFunc{
		_message.ActionHeartbeat:  d.handleHeartbeat,
		_message.ActionAckRequest: d.handleAckRequest,
	}
	for action, handlerFunc := range m {
		if callback != nil {
			handlerFunc = callback(action, handlerFunc)
		}
		// set action handler in handler chain
		d.def.AddHandler(NewActionHandler(action, handlerFunc))
	}

	//d.def.AddHandler(&ClientCustomMessageHandler{})
}

func (d *MessageHandlerImpl) AddHandler(i MessageHandler) {
	d.def.AddHandler(i)
}

func (d *MessageHandlerImpl) Handle(clientInfo *_client.Info, message *_message.Message) error {
	return d.def.Handle(clientInfo, message)
}

func (d *MessageHandlerImpl) SetGate(gate _client.Gateway) {
	d.def.SetGate(gate)
}

func (d *MessageHandlerImpl) SetSubscription(s subscription.Interface) {
	d.def.SetSubscription(s)
}

func (d *MessageHandlerImpl) enqueueMessage(id _client.ID, message *_message.Message) {
	err := d.def.GetClientInterface().EnqueueMessage(id, message)
	if err != nil {
		d.ctx.Logger.Errorf("EnqueueMessage err : %s", err)
	}
}

func (d *MessageHandlerImpl) unmarshalData(message *_message.Message, to interface{}) bool {
	err := message.Data.Deserialize(to)
	if err != nil {
		d.ctx.Logger.Errorf("sender chat senderMsg error : %s", err)
		return false
	}
	return true
}
