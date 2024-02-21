package message

// AckRequest receiver to server
type AckRequest struct {
	Seq  int64  `json:"seq,omitempty"`
	Mid  int64  `json:"mid,omitempty"`
	From string `json:"from,omitempty"`
}

// AckMessage server to sender (notify client to server)
type AckMessage struct {
	Seq    int64  `json:"seq,omitempty"`
	CliMid string `json:"cliMid,omitempty"`
	Mid    int64  `json:"mid,omitempty"`
}

// AckNotify sender to receiver
type AckNotify struct {
	Mid int64 `json:"mid,omitempty"`
}
