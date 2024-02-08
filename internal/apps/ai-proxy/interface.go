package ai_proxy

type Interface interface {
	GetAnswer()
}

type Msg struct {
	Content string
}
