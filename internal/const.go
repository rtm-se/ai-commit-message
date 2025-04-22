package internal

const Prompt2 = "Write a professional short git commit message based on the a diff below in English language. Do not preface the commit with anything, use the present tense, return the full sentence, and use the conventional commits specification (<type in lowercase>: <subject>):\n"
const Prompt = "You are a senior software engineer.\nYou will be given a git diff.\nGenerate a single, concise one line commit message in imperative tense (e.g. “Fix bug”, “Add feature”) that accurately summarizes the changes. Do not preface the commit with anything, use the present tense, return the full sentence in one line, and use the conventional commits specification (<type in lowercase>: <subject>). Limit your output to at most 150 characters. **Output only the commit message nothing else**.\nDiff:\n"
const LoopPrompt = "Shorten following text to a professional short git commit message in English language, do not preface the commit with anything, use the present tense, return the full sentence, and use the conventional commits specification (<type in lowercase>: <subject>):\n"
const RegenerateForLengthPrompt = "You will be provided with a long commit message, you need to extract the core of the message so the sum of the wouldn't exceed: {length} symbols, do not add explanation in the final answer, only the commit message"

var LengthPrompt = "The length of the git commit should be no longer than %v symbols"

const ConfigFileName = ".auto-commit-settings.yaml"
