package ollama

import (
	"context"

	"github.com/jonathanhecl/gollama"
)

type OllamaClient struct {
	client *gollama.Gollama
}

func NewOllamaClient(model string) *OllamaClient {
	ollama := gollama.New(model)
	return &OllamaClient{client: ollama}
}

func (lama *OllamaClient) GetResponse(fullPrompt string) string {
	ctx := context.Background()
	output, err := lama.client.Chat(ctx, fullPrompt)
	if err != nil {
		panic(err)
	}
	return output.Content
}
