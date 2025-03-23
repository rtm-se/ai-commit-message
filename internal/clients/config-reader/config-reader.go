package config_reader

import (
	"flag"
	"log"

	constants "github.com/rtm-se/ai-commit-message/internal"
)

type Config struct {
	Prompt          string
	Model           string
	CLeanThinkBlock bool
	SeparateDiff    bool
	LoopPrompt      string
	Loop            bool
	LLMEndpoint     string
	Interactive     bool
}

type configBuilder struct {
	model           *string
	cleanThinkBlock *bool
	separateDiff    *bool
	loop            *bool
	interactive     *bool

	llmEndpoint *string
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

func (builder *configBuilder) SetLoopFromFlag() *configBuilder {
	builder.loop = flag.Bool("loop", false, "feed response into llm again to gain shortened result")
	return builder
}

func (builder *configBuilder) SetApiEndpointFromFlag() *configBuilder {
	builder.llmEndpoint = flag.String("llm-endpoint", "http://localhost:11434", "llm endpoint to use; default: ollama default api endpoint")
	return builder
}

func (builder *configBuilder) SetCleanThinkBlock() *configBuilder {
	builder.cleanThinkBlock = flag.Bool("clean-think", false, "should clean <think></think> block form model response")
	return builder
}

func (builder *configBuilder) SetInteractive() *configBuilder {
	builder.interactive = flag.Bool("interactive", true, "will prompt with options step by step during process")
	return builder
}

func (builder *configBuilder) BuildConfig() *Config {
	flag.Parse()
	if *builder.interactive {
		log.Println("Starting in interactive mode")
	}
	log.Printf("model: %v \n", *builder.model)
	log.Printf("clean think block: %v, \n", *builder.cleanThinkBlock)
	log.Printf("separate diff: %v \n", *builder.separateDiff)
	if *builder.loop {
		log.Printf("loop: %v \n", *builder.loop)
	}
	return &Config{
		Model:           *builder.model,
		CLeanThinkBlock: *builder.cleanThinkBlock,
		SeparateDiff:    *builder.separateDiff,
		Loop:            *builder.loop,
		Prompt:          constants.Prompt,
		LoopPrompt:      constants.LoopPrompt,
		LLMEndpoint:     *builder.llmEndpoint,
		Interactive:     *builder.interactive,
	}
}
