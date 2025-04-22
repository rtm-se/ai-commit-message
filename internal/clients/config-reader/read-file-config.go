package config_reader

import (
	"github.com/rtm-se/ai-commit-message/internal"
	"gopkg.in/yaml.v3"

	"log"
	"os"
	"regexp"
	"strings"
)

type ignorePattern struct {
	Message  string   `yaml:"message"`
	Patterns []string `yaml:"patterns"`
}

type configFile struct {
	FlagsOverConfig bool            `yaml:"flags_over_config"`
	IgnorePatterns  []ignorePattern `yaml:"ignore_patterns"`
}

func readFileConfig() configFile {
	conf := &configFile{}
	home, err := os.UserHomeDir()
	yamlFile, err := os.ReadFile(home + "/" + internal.ConfigFileName)

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return *conf
}

func (builder *configBuilder) CollectSettingsFromConfigFile() *configBuilder {
	config := readFileConfig()
	builder.ignorePatterns = make([]IgnoreFilesPattern, len(config.IgnorePatterns))
	for i, iPattern := range config.IgnorePatterns {
		builder.ignorePatterns[i] = newIgnorePattern(iPattern.Message, iPattern.Patterns)
	}
	builder.flagsOverConfig = &config.FlagsOverConfig
	return builder
}

func newIgnorePattern(message string, patterns []string) IgnoreFilesPattern {
	return IgnoreFilesPattern{
		Message:  message,
		Patterns: regexp.MustCompile(strings.Join(patterns, "|")),
	}
}
