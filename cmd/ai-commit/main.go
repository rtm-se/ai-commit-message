package main

import (
	"log"

	"github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
)
import "github.com/rtm-se/ai-commit-message/internal/app"

func main() {
	log.Println("ai-commit started")
	config := config_reader.NewConfig()
	a := app.NewApp(config)
	a.StageAllFiles()
	commitMessage := a.CreateCommit()
	log.Println(commitMessage)
	a.CommitWithMessage(commitMessage)
}
