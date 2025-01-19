package main

import "github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
import "github.com/rtm-se/ai-commit-message/internal/app"

func main() {
	println("ai-commit started")
	config := config_reader.NewConfig()
	a := app.NewApp(config)
	commit := a.CreateCommit()
	println(commit)
}
