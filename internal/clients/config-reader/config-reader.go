package config_reader

import (
	"flag"
	"fmt"
)

type Config struct {
	Prompt          string
	Model           string
	CLeanThinkBlock bool
	SeparateDiff    bool
}

type configBuilder struct {
	model           *string
	cleanThinkBlock *bool
	separateDiff    *bool
}

func NewConfigBuilder() *configBuilder {
	return &configBuilder{}
}

func (builder *configBuilder) SetModelFromFlag() *configBuilder {
	builder.model = flag.String("model", "mistral", "Ollama model you want to use; default: mistral")
	return builder
}

func (builder *configBuilder) SetSeparateFilesFromFlag() *configBuilder {
	builder.separateDiff = flag.Bool("separate-diff-into-files", true, "feed whole diff into llm or separate into chunks")
	return builder
}

func (builder *configBuilder) SetCleanThinkBlock() *configBuilder {
	builder.cleanThinkBlock = flag.Bool("clean-think", false, "should clean <think></think> block form model response")
	return builder
}

func (builder *configBuilder) BuildConfig() *Config {
	flag.Parse()
	fmt.Println("model: " + *builder.model)
	return &Config{
		Model:           *builder.model,
		CLeanThinkBlock: *builder.cleanThinkBlock,
		Prompt: "Write a professional short git commit message based on the a diff below in English language\n" +
			"Do not preface the commit with anything, use the present tense, return the full sentence, and use the conventional commits specification (<type in lowercase>: <subject>):\n",
	}
}
