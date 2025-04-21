package gemini

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const LLMClientName = "gemini"

type GeminiClient struct {
	client *genai.Client
	apiKey string
	model  *genai.GenerativeModel
	ctx    context.Context
}

func NewGeminiClient(ctx context.Context, apiKey, model string) *GeminiClient {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	gc := GeminiClient{
		client: client,
		apiKey: apiKey,
		ctx:    ctx,
	}
	gc.ChangeModel(model)
	return &gc
}

func (gc *GeminiClient) GetAvailableModels() []string {

	iter := gc.client.ListModels(gc.ctx)
	models := []string{}
	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		models = append(models, m.Name)
	}
	return models
}
func (gc *GeminiClient) Close() {
	err := gc.client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
func (gc *GeminiClient) GetCurrentModelName() string {
	modelInfo, _ := gc.model.Info(gc.ctx)
	return modelInfo.Name
}

func (gc *GeminiClient) ChangeModel(model string) {
	gc.model = gc.client.GenerativeModel(model)
}

func (gc *GeminiClient) GetResponse(ctx context.Context, fullPrompt string) (string, error) {
	chatSession := gc.model.StartChat()
	iter := chatSession.SendMessageStream(ctx, genai.Text(fullPrompt))
	message := strings.Builder{}
	for {
		response, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		collectResponse(response, &message)
	}
	return message.String(), nil
}

func collectResponse(resp *genai.GenerateContentResponse, stringWriter *strings.Builder) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				stringWriter.WriteString(fmt.Sprintf("%v", part))
			}
		}
	}
}
