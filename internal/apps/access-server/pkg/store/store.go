package store

import _message "Aurora/internal/apps/access-server/pkg/message"

type MessageStore interface {
	StoreMessage(message *_message.ChatMessage) error
	StoreOffline(message *_message.ChatMessage) error
}

type SubscriptionStore interface {
	//NextSegmentSequence(id)
}
