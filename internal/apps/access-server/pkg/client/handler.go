package client

import _message "Aurora/internal/apps/access-server/pkg/message"

type MessageHandler func(cliInfo *Info, message *_message.Message)

func (c *UserClient) handleHello(m *_message.Message) {
	hello := _message.Hello{}
	err := m.Data.Deserialize(&hello)
	if err != nil {
		_ = c.EnqueueMessage(_message.NewMessage(
			0, _message.ActionNotifyError, "invalid handleHello message"))
		return
	}
	c.info.Version = hello.ClientVersion
}
