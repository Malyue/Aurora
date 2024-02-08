package ai_proxy

type AiType string

const (
	GPT3DOT5   AiType = "gpt-3.5-turbo"
	GPT4       AiType = "gpt-4"
	GPTCodex   AiType = "gpt-codex"
	Embeddings AiType = ""
)
