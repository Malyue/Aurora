package subscription

import _client "Aurora/internal/apps/access-server/pkg/client"

type Interface interface {
	PublishMessage(id ChanID, message Message) error
}

type Subscribe interface {
	Interface

	SetGateInterface(gate _client.Gateway)

	UpdateSubscriber(id ChanID, updates []Update) error

	UpdateChannel(id ChanID, update ChannelUpdate) error
}

type ChanID string

type SubscriberID string

type Update struct {
	Flag int64
	ID   SubscriberID

	Extra interface{}
}

type ChannelUpdate struct {
	Flag int64

	Extra interface{}
}

type Server interface {
	Subscribe
	Run() error
}
