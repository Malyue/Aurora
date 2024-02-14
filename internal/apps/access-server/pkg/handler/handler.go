package handler

import (
	_client "Aurora/internal/apps/access-server/pkg/client"
	_message "Aurora/internal/apps/access-server/pkg/message"
	"Aurora/internal/apps/access-server/pkg/store"
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
	log            *logrus.Logger
}

type MessageHandlerImpl struct {
	def *MessageInterfaceImpl
	// TODO get message model
	store     store.MessageStore
	userState *UserState
}

func NewHandler(opt *Options, log *logrus.Logger) (*MessageInterfaceImpl, error) {
	ret := MessageInterfaceImpl{
		hc:             &handlerChain{},
		notifyOnSrvErr: opt.NotifyServerError,
		log:            log,
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

func (m *MessageInterfaceImpl) onMessageHandlerPanic(i interface{}) {
	if m.log != nil {
		m.log.Error("MessageInterfaceImpl panic : %v", i)
		return
	}
	logrus.Error("MessageInterfaceImpl panic : %v", i)
}

func (m *MessageInterfaceImpl) Handle(clientInfo *_client.Info, msg *_message.Message) error {
	if !msg.GetAction().IsInternal() {
		msg.From = clientInfo.ID.UID()
	}

	err := m.execPool.Submit(func() {
		handled := m.hc.handle(m, clientInfo, msg)
		if !handled {
			if !msg.GetAction().IsInternal() {
				r := _message.NewMessage(msg.GetSeq(), _message.ActionNotifyUnknownAction, msg.GetAction())
			}
		}
	})
}
