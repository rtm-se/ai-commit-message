package app

func (a *AppAICommit) PreCommitActions() {
	a.commitIgnoredGroups()
}

func (a *AppAICommit) commitIgnoredGroups() {
	groups := a.gitClient.GetIgnoredGroups()
	var atLestOneGroupFound bool
	for _, group := range groups {
		if len(group) != 0 {
			atLestOneGroupFound = true
		}
	}
	if !atLestOneGroupFound {
		return
	}
	if !a.ShouldCommitIgnoredFiles() {
		return
	}
	for groupName, groupFiles := range groups {
		for _, file := range groupFiles {
			a.gitClient.StageFile(file)
		}
		a.gitClient.Commit(a.GetCommitPrefix() + groupName)
	}

}
