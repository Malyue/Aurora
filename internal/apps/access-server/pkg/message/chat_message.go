package message

type ChatMessage struct {
	// CliMid client message id to identity unique a message
	// for identity a message and wait for the server ack receipt and return `mid` for it
	CliMid string `json:"cliMid,omitempty"`
	// Mid server message id in the db
	// when a client sends a message for the first time or client retry to send a message
	// that the server does not ack, the 'Mid' is empty
	// if this field is not empty that this message is server acked, need not store to db again
	Mid int64 `json:"mid,omitempty"`
	// Seq message sequence for a chat, use to check message whether the message lost.
	Seq     int64  `json:"seq,omitempty"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	Type    int32  `json:"type,omitempty"`
	Content string `json:"content,omitempty"`
	SendAt  int64  `json:"sendAt,omitempty"`
}
