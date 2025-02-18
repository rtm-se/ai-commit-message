package main

import (
	"log"

	"github.com/rtm-se/ai-commit-message/internal/app"
	"github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
)

func main() {
	builder := config_reader.NewConfigBuilder()
	builder.SetModelFromFlag().SetCleanThinkBlock().SetSeparateFilesFromFlag().SetLoopFromFlag()
	cfg := builder.BuildConfig()
	a := app.NewApp(cfg)
	log.Println("ai-commit started")
	a.StageAllFiles()
	commitMessage := a.CreateCommit()
	log.Println(commitMessage)

	if a.ShouldLoopResponse() {
		loopedCommitMessage := a.LoopForFeedback(commitMessage)
		log.Println(loopedCommitMessage)
	}
	if !a.ShouldCommit() {
		log.Println("Won't commit message, exiting...")
		return
	}
	a.CommitWithMessage(a.GetCommitPrefix() + commitMessage)
}
