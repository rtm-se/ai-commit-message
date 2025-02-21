# AI Commit Message Application

This repository contains the source code for an AI-powered commit message
generator named `ai-commit-message`. The application is designed to
automate the creation of meaningful and descriptive commit messages based
on the changes made in a project.

## Getting Started

These instructions will get you a copy of the project up and running on
your local machine for development and testing purposes.

### Dependencies

- Go (version 1.23 or later) - Download from https://golang.org/dl/
- Ollama with mistrall or deepseek-r1 installed - Download from https://ollama.com/
  after installing run the following
   ```
  ollama run deepseek-r1
   ```
  or
   ```
  ollama run mistral
   ```

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/rtm-se/ai-commit-message.git
   ```
2. Navigate to the project directory:
   ```
   cd ai-commit-message
   ```
3. Build and run the application:
   ```
   go build ./cmd/ai-commit/main.go
   ```
4. navigate to your directory in which you wanna generate a commit and execute the file
   ```
   ~/path-to-executable/main
   ```
   To use model with deepseek you can use following flags, by default it will use mistral model 
   ```
   ~/path/to/executable/main -model=deepseek-r1 -clean-think=true
   ```

upcoming changes:
-
- ~~1 - feeding diff by files should yield better results~~
- ~~2 - feedback loop on commit messages~~
- 3 - trying different models in one go
- 4 - add flags for different behaviours
- 5 - add support for remote models
- 6 - support opensource local model providers
- 7 - separate commits into several chunks for even better results
- 8 - add better documentation for flags options
- 9 - improve prompts for better results