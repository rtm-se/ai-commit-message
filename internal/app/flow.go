package app

import "fmt"

func (a *AppAICommit) changeModel() {
	options := a.LLMClient.GetAvailableModels()
	//TODO: track edge-case on exit from shell
	resp := a.shell.HandleMultipleInput("What model do you want to use:", options)
	a.LLMClient.ChangeModel(options[resp])
}

func (a *AppAICommit) shouldUseDifferentModel(commitMessage string) bool {
	options := []string{"no", "yes"}
	questionMessage := fmt.Sprintf("%v\nCurrent Model: %v Do you want to try different model?", commitMessage, a.LLMClient.GetCurrentModelName())
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
		options := []string{"no", "yes"}
		resp := a.shell.HandleMultipleInput(message, options)
		switch resp {
		case 0:
			return false
		case 1:
			return true
		}
	}
	return a.config.Loop
}

func (a *AppAICommit) ShouldCommitIgnoredFiles() bool {
	if a.config.Interactive {
		message := fmt.Sprintf("Files from ignored groups detected, should the by commited by their group name?")
		options := []string{"yes", "no"}
		resp := a.shell.HandleMultipleInput(message, options)
		switch resp {
		case 0:
			return true
		case 1:
			return false
		}
	}
	return false
}
