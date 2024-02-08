package gpt

import (
	"context"

	"Aurora/internal/apps/ai-proxy/svc"
	"github.com/sashabaranov/go-openai"
)

type Instance struct {
	Proxy string
}

func (gpt Instance) Chat(svcCtx *svc.ServerCtx, token string) (*openai.ChatCompletionResponse, error) {
	clientConfig := openai.DefaultConfig(token)
	clientConfig.BaseURL = gpt.Proxy
	client := openai.NewClientWithConfig(clientConfig)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "!",
				},
			},
		})
	if err != nil {
		return nil, err
	}

	return &resp, err

	//return resp.Choices[0].Message.Content, nil

}

func (gpt Instance) GetModel(token string) (openai.ModelsList, error) {
	clientConfig := openai.DefaultConfig(token)
	clientConfig.BaseURL = gpt.Proxy
	client := openai.NewClientWithConfig(clientConfig)
	list, err := client.ListModels(context.Background())
	return list, err
}
