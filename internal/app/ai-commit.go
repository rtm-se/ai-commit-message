package app

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	config_reader "github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
	"github.com/rtm-se/ai-commit-message/internal/clients/gemini"
	"github.com/rtm-se/ai-commit-message/internal/clients/git"
	"github.com/rtm-se/ai-commit-message/internal/clients/ollama"
	"github.com/rtm-se/ai-commit-message/internal/clients/shell"
	"github.com/rtm-se/ai-commit-message/internal/clients/spinner"
)

type AppAICommit struct {
	gitClient *git_client.GitClient
	config    *config_reader.Config
	LLMClient LLM
	shell     *shell.Shell
}

type LLM interface {
	ChangeModel(model string)
	GetResponse(ctx context.Context, fullPrompt string) (string, error)
	GetCurrentModelName() string
	GetAvailableModels() []string
}

func NewLLMClient(ctx context.Context, config *config_reader.Config) (LLM, error) {
	switch config.LLMClientName {
	case gemini.LLMClientName:
		if config.LLMKeys[gemini.LLMClientName] == "" {
			return nil, fmt.Errorf("gemini: LLM client key is required")
		}
		return gemini.NewGeminiClient(ctx, config.LLMKeys[gemini.LLMClientName], config.Model), nil
	case ollama.LLMClientName:
		//TODO: Uncomment this when there will be remote ollama with auth support
		//if config.LLMKeys[ollama.LLMClientName] == "" {
		//	return nil, fmt.Errorf("ollama: LLM client key is required")
		//}
		return ollama.NewOllamaClient(config.Model, config.LLMEndpoint), nil
	}
	return nil, fmt.Errorf("[NewLLMClient]LLM client name is invalid")
}
func NewApp(ctx context.Context, config *config_reader.Config) (*AppAICommit, error) {
	gc := git_client.NewGitClient(config)
	sh := shell.NewShell()
	llmClient, err := NewLLMClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("NewLLMClient: %v", err)
	}
	return &AppAICommit{
		gitClient: gc,
		config:    config,
		LLMClient: llmClient,
		shell:     sh,
	}, nil
}

func (a *AppAICommit) prepareDiff() (string, error) {

	return a.gitClient.GetDiff(), nil
}

func (a *AppAICommit) GetCommitPrefix() string {
	if a.config.CustomPrefix != "" {
		return fmt.Sprintf("[%v]", a.config.CustomPrefix)
	}
	if a.config.RepeatPrefix {
		previousCommitPrefix := a.gitClient.GePreviousCommitPrefix()
		if previousCommitPrefix != "" {
			return fmt.Sprintf("[%v]", previousCommitPrefix)
		}
	} // if no prefix | repeat message flags are set, use the ticket number from git branch name
	currentGitBranch := a.gitClient.GetBranch()
	trimmed := strings.SplitAfter(currentGitBranch, "/")
	ticket := trimmed[len(trimmed)-1]
	log.Println("Prefix detected as", ticket)
	return fmt.Sprintf("[%v]", ticket)
}

func (a *AppAICommit) GetCommitMessage() (commitMessage string) {
	a.changeModel()
	for {
		commitMessage = a.generateCommitMessage()
		if !a.shouldUseDifferentModel(commitMessage) {
			break
		}
		a.changeModel()
	}
	return commitMessage
}

func (a *AppAICommit) generateCommitMessage() string {
	InitialCommitMessage, err := a.CreateCommit()
	if err != nil {
		panic(err)
	}
	// TODO: fix this for long commits and re-enable it
	//filteredCommit := a.choseWhatLinesToKeepInCommit(InitialCommitMessage)
	//commitMessage := strings.Join(filteredCommit, "")
	if a.ShouldLoopResponse(InitialCommitMessage) {
		loopedCommitMessage := a.LoopForFeedback(InitialCommitMessage)
		log.Println(loopedCommitMessage)
		return loopedCommitMessage
	}
	return InitialCommitMessage
}

//TODO: fix this to work with big amount of messages
//func (a *AppAICommit) choseWhatLinesToKeepInCommit(commitMessage []string) []string {
//	if len(commitMessage) <= 1 {
//		return commitMessage
//	}
//	options := []string{
//		"yes",
//		"no",
//	}
//	resp := a.shell.HandleMultipleInput("Do you want to filter lines for the messages:", options)
//	switch resp {
//	case 0:
//		return a.filterCommitMessage(commitMessage)
//	case 1:
//		return commitMessage
//	default:
//		panic("Unexpected response")
//	}
//
//}

//func (a *AppAICommit) filterCommitMessage(commitMessages []string) []string {
//	resp := a.shell.HandleCHeckBoxInput("What lines should we keep, mark lines to delete, ctrl+c to keep all", commitMessages)
//	var filteredCommitMessages []string
//	switch len(resp) {
//	case 0:
//		return commitMessages
//	default:
//		for i, _ := range commitMessages {
//			if slices.Contains(resp, i) {
//				continue
//			}
//			filteredCommitMessages = append(filteredCommitMessages, commitMessages[i])
//		}
//		return filteredCommitMessages
//	}
//
//}

func (a *AppAICommit) shouldFilterMessagesByLength(messageLength int) bool {
	if a.config.AutoRejectLongMessages <= 0 {
		return false
	}
	if messageLength > a.config.AutoRejectLongMessages {
		return true
	}
	return false
}

func (a *AppAICommit) regenerateCommitMessageWithLengthConstraints(message string) (string, error) {
	response, err := a.getResponseFromLLM(a.config.RegenerateForLengthPrompt + strconv.Itoa(a.config.AutoRejectLongMessages) + "\n" + message)
	if err != nil {
		return "", fmt.Errorf("[regenerateCommitMessageWithLengthConstraints] %v", err)
	}
	if a.shouldFilterMessagesByLength(len(response)) {
		log.Printf("Regenerated message is too long:\n%v", response)
		return "", fmt.Errorf("[regenerateCommitMessageWithLengthConstraints]regenerated message rejected by length: %v", len(response))
	}

	return response, nil

}
func (a *AppAICommit) CreateCommit() (string, error) {
	prompts := a.getPrompts()
	if a.config.SeparateDiff {
		log.Printf("Detected prompts: %v", len(prompts))
	}
	commitMessage := strings.Builder{}
	spn := spinner.NewSpinner()
	for _, prompt := range prompts {
		go spn.Spin()
		partialCommitMessage, err := a.getResponseFromLLM(prompt)
		spn.Stop()
		if err != nil {
			return "", fmt.Errorf("[CreateCommit] %v", err)
		}
		if a.shouldFilterMessagesByLength(len(partialCommitMessage)) {
			log.Printf("Regenerating message:\n %v", partialCommitMessage)
			partialCommitMessage, err = a.regenerateCommitMessageWithLengthConstraints(partialCommitMessage)
			log.Printf("Regenerated message:\n %v", partialCommitMessage)
			if err != nil {
				continue
			}
		}
		if len(prompts) > 1 {
			log.Println(partialCommitMessage)
		}

		commitMessage.WriteString(partialCommitMessage)
		commitMessage.WriteString("\n")
	}
	return commitMessage.String(), nil
}

func (a *AppAICommit) getResponseFromLLM(prompt string) (string, error) {
	ctx := context.Background()
	LLMResponse, err := a.LLMClient.GetResponse(ctx, prompt)
	if err != nil {
		return "", err
	}
	if a.config.CLeanThinkBlock {
		LLMResponse = a.deleteThinkBlockFromModelResponse(LLMResponse)
	}
	return LLMResponse, nil
}

func (a *AppAICommit) deleteThinkBlockFromModelResponse(response string) string {
	ss := strings.SplitAfter(response, "</think>")
	return strings.Replace(ss[len(ss)-1], "\n", "", 2)
}

func (a *AppAICommit) StageAllFiles() {
	log.Println("Staging files")
	err := a.gitClient.Stage()
	if err != "" {
		panic(err)
	}
}

func (a *AppAICommit) CommitWithMessage(message string) {
	s := a.gitClient.Commit(message)
	log.Println(s)
	log.Printf("Changes Committed")
}

func (a *AppAICommit) getLoopPrompt(commitMessage string) string {
	return fmt.Sprintf(a.config.LoopPrompt + "\n" + commitMessage)
}
func (a *AppAICommit) LoopForFeedback(commitMessage string) string {
	loopPrompt := a.getLoopPrompt(commitMessage)
	spn := spinner.NewSpinner()
	go spn.Spin()
	loopedBack, err := a.getResponseFromLLM(loopPrompt)
	if err != nil {
		panic(err)
	}
	spn.Stop()
	return loopedBack
}
