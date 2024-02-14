package message

import "strings"

type Action string

const (
	ActionHello     Action = "hello"
	ActionHeartbeat        = "heartbeat"

	ActionChatMessage = "message.chat"

	ActionGroupMessage = "message.group"

	// ========= ACK ==========
	ActionAckRequest  = "ack.request"
	ActionAckGroupMsg = "ack.group.msg"
	ActionAckMessage  = "ack.message"

	ActionNotifyUnknownAction = "notify.unknown.action"
	ActionNotifyError         = "notify.error"
	ActionNotifySuccess       = "notify.success"

	ActionInternalOnline  = "internal.online"
	ActionInternalOffline = "internal.offline"
)

func (a Action) IsInternal() bool {
	return strings.HasPrefix(string(a), "internal.")
}
