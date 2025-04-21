package app

import "fmt"

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
