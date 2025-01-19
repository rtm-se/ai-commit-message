package app

import (
	"fmt"
	"log"
	"strings"

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

func (a *AppAICommit) getCommitPrefix() string {
	currentGitBranch := a.gitClient.GetBranch()
	trimmed := strings.SplitAfter(currentGitBranch, "/")
	ticket := trimmed[len(trimmed)-1]
	log.Println("Prefix detected as", ticket)
	return fmt.Sprintf("[%v]", ticket)
}
func (a *AppAICommit) CreateCommit() string {
	prompt := a.prepareFullPrompt()
	commitMessage := a.Ollama.GetResponse(prompt)
	prefix := a.getCommitPrefix()
	return prefix + commitMessage
}

func (a *AppAICommit) StageAllFiles() {
	stagedFiles := a.gitClient.Stage()
	log.Printf("Staged files\n %v", stagedFiles)
}

func (a *AppAICommit) CommitWithMessage(message string) {
	a.gitClient.Commit(message)
	log.Printf("Changes Committed")
}
