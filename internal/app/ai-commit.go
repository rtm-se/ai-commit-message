package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	config_reader "github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
	"github.com/rtm-se/ai-commit-message/internal/clients/git"
	"github.com/rtm-se/ai-commit-message/internal/clients/ollama"
	"github.com/rtm-se/ai-commit-message/internal/clients/spinner"
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

func (a *AppAICommit) prepareFullPrompt() []string {
	fullPrompt := make([]string, 0)
	diff := a.gitClient.GetDiff()
	fullPrompt = append(fullPrompt, fmt.Sprintf(a.config.Prompt+diff))
	return fullPrompt
}

func (a *AppAICommit) preparePromptsByFiles() []string {
	diffs := a.gitClient.GetSeparatedDiffs()
	prompts := make([]string, 0)
	for _, diff := range diffs {
		if diff == "" {
			continue
		}
		prompts = append(prompts, fmt.Sprintf(a.config.Prompt+"\n"+diff))
	}
	return prompts
}

func (a *AppAICommit) getCommitPrefix() string {
	currentGitBranch := a.gitClient.GetBranch()
	trimmed := strings.SplitAfter(currentGitBranch, "/")
	ticket := trimmed[len(trimmed)-1]
	log.Println("Prefix detected as", ticket)
	return fmt.Sprintf("[%v]", ticket)
}

func (a *AppAICommit) getPrompts() []string {
	if !a.config.SeparateDiff {
		return a.prepareFullPrompt()
	}
	return a.preparePromptsByFiles()
}

func (a *AppAICommit) CreateCommit() string {
	prompts := a.getPrompts()
	if a.config.SeparateDiff {
		log.Printf("Detected %v files to prompt", len(prompts))
	}
	commitMessage := strings.Builder{}
	commitMessage.WriteString(a.getCommitPrefix())
	spn := spinner.NewSpinner()
	for _, prompt := range prompts {
		go spn.Spin()
		partialCommitMessage := a.getResponseFromLLM(prompt)
		spn.Stop()
		if len(prompts) > 1 {
			log.Println(partialCommitMessage)
		}
		commitMessage.WriteString(partialCommitMessage)
		commitMessage.WriteString("\n")
	}
	return commitMessage.String()
}

func (a *AppAICommit) getResponseFromLLM(prompt string) string {
	LLMResponse := a.Ollama.GetResponse(prompt)
	if a.config.CLeanThinkBlock {
		LLMResponse = a.deleteThinkBlockFromModelResponse(LLMResponse)
	}
	return LLMResponse
}
func (a *AppAICommit) ShouldCommit() bool {
	fmt.Println("should we commit with the message above?")
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func (a *AppAICommit) deleteThinkBlockFromModelResponse(response string) string {
	ss := strings.SplitAfter(response, "</think>")
	return strings.Replace(ss[len(ss)-1], "\n", "", 2)
}

func (a *AppAICommit) StageAllFiles() {
	log.Printf("Staging files")
	err := a.gitClient.Stage()
	if err != "" {
		panic(err)
	}
}

func (a *AppAICommit) CommitWithMessage(message string) {
	a.gitClient.Commit(message)
	log.Printf("Changes Committed")
}
