package git_client

import (
	"os/exec"
	"strings"
)

type GitCLient struct {
}

func NewGitClient() *GitCLient {
	return &GitCLient{}
}

func (g *GitCLient) GetDiff() string {
	c, b := exec.Command("git", "diff", "--staged"), new(strings.Builder)
	c.Stdout = b
	c.Run()
	s := strings.TrimRight(b.String(), "\n")
	return s
}

func (g *GitCLient) Stage() string {
	c, b := exec.Command("git", "add", "."), new(strings.Builder)
	c.Stdout = b
	c.Run()
	s := strings.TrimRight(b.String(), "\n")
	return s
}
func (g *GitCLient) Commit(message string) string {
	c, b := exec.Command("git", "commit", "-m", message), new(strings.Builder)
	c.Stdout = b
	c.Run()
	s := strings.TrimRight(b.String(), "\n")
	return s
}

func (g *GitCLient) GetBranch() string {
	c, b := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD"), new(strings.Builder)
	c.Stdout = b
	c.Run()
	s := strings.TrimRight(b.String(), "\n")
	return s
}
