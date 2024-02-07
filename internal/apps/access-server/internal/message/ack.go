package message

import "time"

type AckMessage struct {
	ID       uint64
	AckID    uint64
	SendTime *time.Time
}
