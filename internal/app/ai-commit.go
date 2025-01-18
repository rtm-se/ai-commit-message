package app

import (
	"fmt"

	config_reader "github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
	"github.com/rtm-se/ai-commit-message/internal/clients/git"
	"github.com/rtm-se/ai-commit-message/internal/clients/ollama"
)

type AppAICommit struct {
	gitClient *git_client.GitCLient
	config    *config_reader.Config
	Ollama    *ollama.OllamaClient
}

func NewApp(config *config_reader.Config) *AppAICommit {
	gc := git_client.NewGitClient()
	ol := ollama.NewOllamaClient(config.Model)
	return &AppAICommit{
		gitClient: gc,
		config:    config,
		Ollama:    ol,
	}
}

func (a *AppAICommit) prepareDiff() (string, error) {
	return a.gitClient.GetDiff(), nil
}

func (a *AppAICommit) prepareFullPrompt() string {
	diff := a.gitClient.GetDiff()
	return fmt.Sprintf(a.config.Prompt + diff)
}

func (a *AppAICommit) CreateCommit() string {
	prompt := a.prepareFullPrompt()
	a.Ollama.GetResponse(prompt)
	return ""
}
