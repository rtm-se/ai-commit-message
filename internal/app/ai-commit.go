package app

import (
	"fmt"
	"log"
	"strings"

	config_reader "github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
	"github.com/rtm-se/ai-commit-message/internal/clients/git"
	"github.com/rtm-se/ai-commit-message/internal/clients/ollama"
	"github.com/rtm-se/ai-commit-message/internal/clients/shell"
	"github.com/rtm-se/ai-commit-message/internal/clients/spinner"
)

type AppAICommit struct {
	gitClient *git_client.GitCLient
	config    *config_reader.Config
	Ollama    *ollama.OllamaClient
	shell     *shell.Shell
}

func NewApp(config *config_reader.Config) *AppAICommit {
	gc := git_client.NewGitClient()
	ol := ollama.NewOllamaClient(config.Model, config.LLMEndpoint)
	sh := shell.NewShell()
	return &AppAICommit{
		gitClient: gc,
		config:    config,
		Ollama:    ol,
		shell:     sh,
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
	diffs := a.gitClient.GetSeparatedDiffByFiles()
	prompts := make([]string, 0)
	for _, diff := range diffs {
		if diff == "" {
			continue
		}
		prompts = append(prompts, fmt.Sprintf(a.config.Prompt+"\n"+diff))
	}
	return prompts
}

func (a *AppAICommit) preparePromptsByBlocks() []string {
	diffs := a.gitClient.GetSeparatedDiffByBlocks()
	prompts := make([]string, 0)
	for _, diff := range diffs {
		if diff == "" {
			continue
		}
		prompts = append(prompts, fmt.Sprintf(a.config.Prompt+"\n"+diff))
	}
	return prompts
}

func (a *AppAICommit) GetCommitPrefix() string {
	currentGitBranch := a.gitClient.GetBranch()
	trimmed := strings.SplitAfter(currentGitBranch, "/")
	ticket := trimmed[len(trimmed)-1]
	log.Println("Prefix detected as", ticket)
	return fmt.Sprintf("[%v]", ticket)
}

func (a *AppAICommit) getPrompts() []string {
	options := []string{"full diff", "separated by files diff", "separated by blocks diff"}
	resp := a.shell.HandleMultipleInput("What kind of model prompting should be used?\n", options)
	switch resp {
	case 0:
		return a.prepareFullPrompt()
	case 1:
		return a.preparePromptsByFiles()
	case 2:
		return a.preparePromptsByBlocks()
	default:
		panic("Unexpected response")
	}
}

func (a *AppAICommit) changeModel() {
	options := a.Ollama.GetAvailableModels()
	//TODO: track edge-case on exit from shell
	resp := a.shell.HandleMultipleInput("What model do you want to use:", options)
	a.Ollama.ChangeModel(options[resp])
}

func (a *AppAICommit) shouldUseDifferentModel(commitMessage string) bool {
	options := []string{"no", "yes"}
	questionMessage := fmt.Sprintf("%v\nCurrent Model: %v Do you want to try different model?", commitMessage, a.Ollama.GetCurrentModelName())
	resp := a.shell.HandleMultipleInput(questionMessage, options)
	switch resp {
	case 0:
		return false
	case 1:
		return true
	default:
		panic("Unexpected response")
	}
}

func (a *AppAICommit) GetCommitMessage() (commitMessage string) {
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
	InitialCommitMessage := a.CreateCommit()
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

func (a *AppAICommit) CreateCommit() string {
	prompts := a.getPrompts()
	if a.config.SeparateDiff {
		log.Printf("Detected prompts: %v", len(prompts))
	}
	commitMessage := strings.Builder{}
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

func (a *AppAICommit) ShouldCommit(commitMessage string) bool {
	options := []string{"yes", "no"}
	questionMessage := fmt.Sprintf("%v\nshould we commit with the message above?", commitMessage)
	responseOptionID := a.shell.HandleMultipleInput(questionMessage, options)
	switch responseOptionID {
	case 0:
		return true
	case 1:
		return false
	default:
		panic("unrecognized output")
	}
}

func (a *AppAICommit) ShouldLoopResponse(commitMessage string) bool {
	if a.config.Interactive {
		message := fmt.Sprintf("%v\nShould we loop back output to LLM?", commitMessage)
		options := []string{"yes", "no"}
		resp := a.shell.HandleMultipleInput(message, options)
		switch resp {
		case 0:
			return true
		case 1:
			return false
		}
	}
	return a.config.Loop
}
func (a *AppAICommit) deleteThinkBlockFromModelResponse(response string) string {
	ss := strings.SplitAfter(response, "</think>")
	return strings.Replace(ss[len(ss)-1], "\n", "", 2)
}

func (a *AppAICommit) Test() {
	a.Ollama.GetAvailableModels()
}
func (a *AppAICommit) StageAllFiles() {
	log.Println("Staging files")
	err := a.gitClient.Stage()
	if err != "" {
		panic(err)
	}
}

func (a *AppAICommit) CommitWithMessage(message string) {
	a.gitClient.Commit(message)
	log.Printf("Changes Committed")
}

func (a *AppAICommit) getLoopPrompt(commitMessage string) string {
	return fmt.Sprintf(a.config.LoopPrompt + "\n" + commitMessage)
}
func (a *AppAICommit) LoopForFeedback(commitMessage string) string {
	loopPrompt := a.getLoopPrompt(commitMessage)
	spn := spinner.NewSpinner()
	go spn.Spin()
	loopedBack := a.getResponseFromLLM(loopPrompt)
	spn.Stop()
	return loopedBack
}
