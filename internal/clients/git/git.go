package git_client

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	config_reader "github.com/rtm-se/ai-commit-message/internal/clients/config-reader"
)

const DiffBlockRE = "@@[^a-z]+[+|-]+[0-9]+[^a-z]+@@"
const MessagePatternRE = "^\\[(\\S*)\\].*"

type GitClient struct {
	IgnorePattern []config_reader.IgnoreFilesPattern
}

func NewGitClient(config *config_reader.Config) *GitClient {
	return &GitClient{
		IgnorePattern: config.IgnorePatterns,
	}
}

func (g *GitClient) GetDiff() string {
	c, b := exec.Command("git", "diff", "--staged"), new(strings.Builder)
	c.Stdout = b
	c.Run()
	s := strings.TrimRight(b.String(), "\n")
	return s
}

func (g *GitClient) Stage() string {
	c, b := exec.Command("git", "add", "."), new(strings.Builder)
	c.Stdout = b
	c.Run()
	s := strings.TrimRight(b.String(), "\n")
	return s
}
func (g *GitClient) Commit(message string) string {
	c, b := exec.Command("git", "commit", "-m", message), new(strings.Builder)
	c.Stdout = b
	c.Stderr = b
	c.Run()
	s := strings.TrimRight(b.String(), "\n")
	return s
}

func (g *GitClient) GetBranch() string {
	c, b := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD"), new(strings.Builder)
	c.Stdout = b
	c.Run()
	s := strings.TrimRight(b.String(), "\n")
	return s
}

func (g *GitClient) GetSeparatedDiffByFiles() []string {
	return strings.Split(g.GetDiff(), "diff --git")
}

func (g *GitClient) GetSeparatedDiffByBlocks() []string {
	re := regexp.MustCompile(DiffBlockRE)
	return re.Split(g.GetDiff(), -1)
}

func (g *GitClient) StageFile(fileName string) {
	c, b := exec.Command("git", "add", fileName), new(strings.Builder)
	c.Stderr = b
	err := c.Run()
	if err != nil {
		log.Fatal(fmt.Sprintf("%v: %v", b.String(), err))
	}
}

func (g *GitClient) ResetToPreviousCommit(soft bool) string {
	args := []string{"reset", "HEAD~1"}
	if soft {
		args = append(args, "--soft")
	}
	c, b := exec.Command("git", args...), new(strings.Builder)
	c.Stdout = b
	c.Run()
	s := strings.TrimRight(b.String(), "\n")
	return s
}

func (g *GitClient) getStatus() string {
	c, b := exec.Command("git", "status", "--short", "--no-renames"), new(strings.Builder)
	c.Stdout = b
	c.Run()
	s := strings.TrimRight(b.String(), "\n")
	return s
}

func (g *GitClient) getFileNamesFromStatus(status string) []string {
	splitStatus := strings.Split(status, "\n")
	files := make([]string, len(splitStatus))
	for _, line := range splitStatus {
		splitLine := strings.Split(line, " ")
		files = append(files, splitLine[len(splitLine)-1])
	}
	return files
}

func (g *GitClient) GetIgnoredGroups() map[string][]string {
	gitStatus := g.getStatus()
	filesFromStatus := g.getFileNamesFromStatus(gitStatus)
	filesPerIgnoreGroups := map[string][]string{}
	for _, pattern := range g.IgnorePattern {
		matchedFile := g.filterFilesByPattern(*pattern.Patterns, filesFromStatus)
		filesPerIgnoreGroups[pattern.Message] = matchedFile
	}
	return filesPerIgnoreGroups
}
func (g *GitClient) filterFilesByPattern(pattern regexp.Regexp, files []string) []string {
	var matchedFiles []string
	for _, file := range files {
		if pattern.MatchString(file) {
			matchedFiles = append(matchedFiles, file)
		}
	}
	return matchedFiles
}

// GePreviousCommitPrefix returns the prefix of previous commit
func (g *GitClient) GePreviousCommitPrefix() string {
	c, b := exec.Command("git", "log", "-n", "1", "--pretty=tformat:%s"), new(strings.Builder)
	c.Stdout = b
	c.Run()
	re := regexp.MustCompile(MessagePatternRE)
	matches := re.FindAllStringSubmatch(b.String(), -1)
	if matches[0] != nil {
		return matches[0][1]
	} else {
		return ""
	}

}

func (g *GitClient) UnstageAll() {
	c, b := exec.Command("git", "restore", "--staged", "."), new(strings.Builder)
	c.Stdout = b
	c.Run()
}
