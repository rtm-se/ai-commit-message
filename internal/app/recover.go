package app

func (a *AppAICommit) Recover() {
	a.gitClient.UnstageAll()
}
