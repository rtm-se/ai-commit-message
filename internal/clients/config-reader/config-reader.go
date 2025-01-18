package config_reader

type Config struct {
	Prompt       string
	GitReposPath string
	Model        string
}

func NewConfig() *Config {
	return &Config{
		Prompt: "Write a professional short git commit message based on the a diff below in English language\n" +
			"Do not preface the commit with anything, use the present tense, return the full sentence, and use the conventional commits specification (<type in lowercase>: <subject>):\n",
		GitReposPath: "",
		Model:        "mistral",
	}
}
