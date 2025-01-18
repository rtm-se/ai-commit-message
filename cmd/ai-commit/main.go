package main

import "fmt"
import "github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
import "github.com/rtm-se/ai-commit-message/internal/app"

func main() {
	fmt.Println("ai-commit started")
	config := config_reader.NewConfig()
	a := app.NewApp(config)
	a.CreateCommit()
}
