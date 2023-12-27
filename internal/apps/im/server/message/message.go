package message

const (
	TEXT = 1
	FILE = 2
)

type Message struct {
	MessageID   int64
	IfGroup     bool
	Type        int
	ReceiverID  string
	GroupID     string
	Subscribers []string
	Payload     []byte

	// -------- queue --------
	private    int64
	index      int
	retryCount int
}
