package config_reader

import (
	"flag"
	"log"
	"os"
	"regexp"

	constants "github.com/rtm-se/ai-commit-message/internal"
	"github.com/rtm-se/ai-commit-message/internal/clients/gemini"
	"github.com/rtm-se/ai-commit-message/internal/clients/ollama"
)

type IgnoreFilesPattern struct {
	Message  string
	Patterns *regexp.Regexp
}
type Config struct {
	Prompt                    string
	Model                     string
	CLeanThinkBlock           bool
	SeparateDiff              bool
	LoopPrompt                string
	Loop                      bool
	LLMEndpoint               string
	Interactive               bool
	RegenerateForLengthPrompt string
	AutoRejectLongMessages    int
	LLMClientName             string
	CustomPrefix              string
	RepeatPrefix              bool
	LLMKeys                   map[string]string
	IgnorePatterns            []IgnoreFilesPattern
	flagsOverConfig           bool
}

type configBuilder struct {
	model                  *string
	cleanThinkBlock        *bool
	separateDiff           *bool
	loop                   *bool
	interactive            *bool
	autoRejectLongMessages *int
	llmEndpoint            *string
	llmKeys                map[string]string
	customPrefix           *string
	repeatPrefix           *bool
	ignorePatterns         []IgnoreFilesPattern
	flagsOverConfig        *bool
}

func NewConfigBuilder() *configBuilder {
	return &configBuilder{}
}

func (builder *configBuilder) CollectApiKeys() *configBuilder {
	builder.llmKeys = map[string]string{
		ollama.LLMClientName: "",
		gemini.LLMClientName: os.Getenv("GEMINI_API_KEY"),
	}
	return builder
}

func (builder *configBuilder) SetLLMClient() *configBuilder {
	builder.model = flag.String("llm-client", "ollama", "llm provider you want to use")
	return builder
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

func (builder *configBuilder) CustomPrefixFromFlag() *configBuilder {
	builder.customPrefix = flag.String("prefix", "", "custom prefix for commit will be inserted into pattern '[{{prefix}}]generated commit message' ")
	return builder
}

func (builder *configBuilder) RepeatPrefixFromFlag() *configBuilder {
	builder.repeatPrefix = flag.Bool("repeat-prefix", false, "repeat prefix from last commit message")
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

func (builder *configBuilder) SetAutoRejectLongMessages() *configBuilder {
	builder.autoRejectLongMessages = flag.Int("auto-reject-length", 150, "will reject messages long than certain length; <= 0 will allow any message length")
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
	log.Printf("loop: %v \n", *builder.loop)
	log.Printf("custom prefix: %v \n", *builder.customPrefix)
	log.Printf("repeat prefix: %v \n", *builder.repeatPrefix)
	if *builder.autoRejectLongMessages > 0 {
		log.Printf("rejecting message longer than: %v symbols \n", *builder.autoRejectLongMessages)
	}
	if *builder.loop {
		log.Printf("loop: %v \n", *builder.loop)
	}
	return &Config{
		Model:                     *builder.model,
		CLeanThinkBlock:           *builder.cleanThinkBlock,
		SeparateDiff:              *builder.separateDiff,
		Loop:                      *builder.loop,
		Prompt:                    constants.Prompt,
		LoopPrompt:                constants.LoopPrompt,
		RegenerateForLengthPrompt: constants.RegenerateForLengthPrompt,
		LLMEndpoint:               *builder.llmEndpoint,
		AutoRejectLongMessages:    *builder.autoRejectLongMessages,
		Interactive:               *builder.interactive,
		LLMKeys:                   builder.llmKeys,
		CustomPrefix:              *builder.customPrefix,
		RepeatPrefix:              *builder.repeatPrefix,
		LLMClientName:             *builder.model,
		IgnorePatterns:            builder.ignorePatterns,
		flagsOverConfig:           *builder.flagsOverConfig,
	}
}
