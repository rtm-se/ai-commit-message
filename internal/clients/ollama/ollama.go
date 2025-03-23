package ollama

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jonathanhecl/gollama"
)

type apiModel struct {
	Name string `json:"name"`
}

type modelsResponse struct {
	Models []apiModel `json:"models"`
}

type OllamaClient struct {
	client      *gollama.Gollama
	apiEndpoint string
}

func NewOllamaClient(model string, endpoint string) *OllamaClient {
	lama := &OllamaClient{
		client:      nil,
		apiEndpoint: endpoint,
	}
	lama.ChangeModel(model)
	return lama
}

func (lama *OllamaClient) GetResponse(fullPrompt string) string {
	ctx := context.Background()
	output, err := lama.client.Chat(ctx, fullPrompt)
	if err != nil {
		panic(err)
	}
	return output.Content
}
func (lama *OllamaClient) ChangeModel(model string) {
	lama.client = gollama.New(model)
	log.Printf("Changed ollama model to %v", model)
}

func (lama *OllamaClient) GetCurrentModelName() string {
	return lama.client.ModelName
}

// TODO: handle unmarshle error
func (lama *OllamaClient) GetAvailableModels() []string {
	modelRSP := modelsResponse{}
	resp := lama.do("/api/tags")
	err := json.Unmarshal(resp, &modelRSP)
	if err != nil {
		panic(err)
	}
	if len(modelRSP.Models) == 0 {
		panic("No models available in Ollama response, make sure you have at least one model locally installed")
	}
	models := make([]string, len(modelRSP.Models))
	for i, model := range modelRSP.Models {
		models[i] = model.Name
	}
	return models
}

func (lama *OllamaClient) do(endpoint string) []byte {
	rsp, err := http.Get(lama.apiEndpoint + endpoint)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	response, errr := ioutil.ReadAll(rsp.Body)
	if errr != nil {
		panic(errr)
	}
	return response
}
