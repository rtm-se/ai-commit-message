package git_client

import (
	"os/exec"
	"regexp"
	"strings"
)

const DiffBlockRE = "@@[^a-z]+[+|-]+[0-9]+[^a-z]+@@"

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

func (g *GitCLient) GetSeparatedDiffByFiles() []string {
	return strings.Split(g.GetDiff(), "diff --git")
}

func (g *GitCLient) GetSeparatedDiffByBlocks() []string {
	re := regexp.MustCompile(DiffBlockRE)
	return re.Split(g.GetDiff(), -1)
}

func (g *GitCLient) ResetToPreviousCommit(soft bool) string {
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
