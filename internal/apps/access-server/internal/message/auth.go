package message

type AuthMessage struct {
	ID    int64
	Token []byte
}

func (a *AuthMessage) GetMsgID() int64 {
	return a.ID
}

func (a *AuthMessage) GetMsgType() MsgType {
	return Auth
}

func (a *AuthMessage) GetReceiverID() []string {
	return []string{}
}

func (a *AuthMessage) GetPayload() []byte {
	return a.Token
}
