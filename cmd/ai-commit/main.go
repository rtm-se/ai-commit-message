package main

import (
	"log"

	"github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
)
import "github.com/rtm-se/ai-commit-message/internal/app"

func main() {
	log.Println("ai-commit started")
	builder := config_reader.NewConfigBuilder()
	builder.SetModelFromFlag().SetCleanThinkBlock()
	cfg := builder.BuildConfig()
	a := app.NewApp(cfg)
	a.StageAllFiles()
	commitMessage := a.CreateCommit()
	log.Println(commitMessage)
	a.CommitWithMessage(commitMessage)
}
