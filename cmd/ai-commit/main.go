package main

import (
	"context"
	"log"

	"github.com/rtm-se/ai-commit-message/internal/app"
	"github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
)

func main() {
	ctx := context.Background()
	builder := config_reader.NewConfigBuilder()
	builder.SetModelFromFlag().SetCleanThinkBlock().SetSeparateFilesFromFlag().SetLoopFromFlag().SetApiEndpointFromFlag().SetInteractive().SetAutoRejectLongMessages().SetLLMClient().CollectApiKeys()
	cfg := builder.BuildConfig()
	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("ai-commit started")
	a.StageAllFiles()
	commitMessage := a.GetCommitMessage()
	if !a.ShouldCommit(commitMessage) {
		log.Println("Won't commit message, exiting...")
		return
	}
	a.CommitWithMessage(a.GetCommitPrefix() + commitMessage)
}
