package message

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
)
