# AI Commit Message Application

This repository contains the source code for an AI-powered commit message
generator named `ai-commit-message`. The application is designed to
automate the creation of meaningful and descriptive commit messages based
on the changes made in a project.

## Getting Started

These instructions will get you a copy of the project up and running on
your local machine for development and testing purposes.

### Dependencies
- git (tested on version 2.48.1)
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

## installation: script
1. use provided following script in a repo root to clean build directory, rebuild project, and add the bin folder with binary to PATH
   ```bash
   source scripts.sh
   ```
    after doing so you will be able to just use `auto-commit` command in any git repo
    

## Installation: manual

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

## make file commands:
```bash
make install
```
Will try to compile project locally into ./bin directory

bad usecases: 
one single change in a long line

Pro-tips:
-
- Even tho deepseek-r1 is supported with local model it's not the best model for speed, use it on more trickier commits and be ready to clean up th final message
- 

upcoming changes:
-
- ~~1 - feeding diff by files should yield better results~~
- ~~2 - feedback loop on commit messages~~
- ~~3 - trying different models in one go~~
- 4 - add flags for different behaviors
- ~~5 - add support for remote models~~
- 6 - support opensource local model providers
- ~~7 - separate commits into several chunks for even better results~~
- 8 - add better documentation for flags options
- 9 ~~- improve prompts for better results~~
- 10 - separate config file
- 11 - support pattern for skipping files
- 12 - interactive mode (semi done, needs improvements)
- ~~13 - automatically reject messages longer than certain number of characters~~
- 14 - support auto git diff context for smaller changes to give llm more lines to process
- 15 - improve visibility on progress for creating commit message
- 16 - custom prefix or use prefix from previous message
- ~~17 - generate md doc for the diff~~(will be done as a plugin)
- 18 - mixed prompt generation (choosing if you want to split files into several contexts)
- 19 - support special diff splitters for better context
- 20 - plugin support
- 21 - add more remote model
- 22 - stage only certain files