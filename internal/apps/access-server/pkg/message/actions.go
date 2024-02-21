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
	ActionAckNotify   = "ack.notify"

	ActionNotifyUnknownAction = "notify.unknown.action"
	ActionNotifyError         = "notify.error"
	ActionNotifySuccess       = "notify.success"
	ActionNotifyForbidden     = "notify.forbidden"

	ActionInternalOnline  = "internal.online"
	ActionInternalOffline = "internal.offline"

	ActionAuthenticate = "authenticate"
)

func (a Action) IsInternal() bool {
	return strings.HasPrefix(string(a), "internal.")
}
