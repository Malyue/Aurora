package handler

import (
	_client "Aurora/internal/apps/access-server/pkg/client"
	_message "Aurora/internal/apps/access-server/pkg/message"
	"Aurora/internal/apps/access-server/pkg/subscription"
	"Aurora/internal/apps/access-server/svc"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
)

type Options struct {
	NotifyServerError     bool `yaml:"notify_server_error"`
	MaxMessageConcurrency int  `yaml:"max_message_concurrency"`
}

// MessageInterfaceImpl implementation of the handler interface
type MessageInterfaceImpl struct {
	execPool *ants.Pool

	hc *handlerChain

	// notifyOnSrvErr notify client on server error
	notifyOnSrvErr bool
	subscription   subscription.Interface
	gate           _client.Gateway

	ctx *svc.ServerCtx
}

// NewDefaultImpl
func NewDefaultImpl(ctx *svc.ServerCtx, opt *Options) (*MessageInterfaceImpl, error) {
	ret := MessageInterfaceImpl{
		hc:             &handlerChain{},
		notifyOnSrvErr: opt.NotifyServerError,
		ctx:            ctx,
	}

	var err error
	ret.execPool, err = ants.NewPool(
		opt.MaxMessageConcurrency,
		ants.WithNonblocking(true),
		ants.WithPanicHandler(ret.onMessageHandlerPanic),
		ants.WithPreAlloc(false),
	)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (d *MessageInterfaceImpl) Handle(clientInfo *_client.Info, msg *_message.Message) error {
	// if the message is not internal, it has the sender
	if !msg.GetAction().IsInternal() {
		msg.From = clientInfo.ID.UID()
	}

	// start a goroutine to handle
	err := d.execPool.Submit(func() {
		// use handle chain to handle the msg
		handled := d.hc.handle(d, clientInfo, msg)
		if !handled {
			if !msg.GetAction().IsInternal() {
				r := _message.NewMessage(msg.GetSeq(), _message.ActionNotifyUnknownAction, msg.GetAction())
				_ = d.gate.EnqueueMessage(clientInfo.ID, r)
			}
			d.ctx.Logger.Warnf("action is not handled: %s", msg.GetAction())
		}
	})
	if err != nil {
		d.onHandleMessageError(clientInfo, msg, err)
		return err
	}
	return nil
}

func (d *MessageInterfaceImpl) onMessageHandlerPanic(i interface{}) {
	if d.ctx.Logger != nil {
		d.ctx.Logger.Errorf("MessageInterfaceImpl panic : %v", i)
		return
	}
	logrus.Errorf("MessageInterfaceImpl panic : %v", i)
}

func (d *MessageInterfaceImpl) onHandleMessageError(clientInfo *_client.Info, msg *_message.Message, err error) {
	if d.notifyOnSrvErr {
		_ = d.gate.EnqueueMessage(clientInfo.ID, _message.NewMessage(-1, _message.ActionNotifyError, err.Error()))
	}
}

func (d *MessageInterfaceImpl) AddHandler(i MessageHandler) {
	d.hc.add(i)
}

func (d *MessageInterfaceImpl) SetGate(gate _client.Gateway) {
	d.gate = gate
}

func (d *MessageInterfaceImpl) SetSubscription(s subscription.Interface) {
	d.subscription = s
}

func (d *MessageInterfaceImpl) SetNotifyErrorOnServer(enable bool) {
	d.notifyOnSrvErr = enable
}

func (d *MessageInterfaceImpl) GetClientInterface() _client.Gateway {
	return d.gate
}

func (d *MessageInterfaceImpl) GetGroupInterface() subscription.Interface {
	return d.subscription
}
