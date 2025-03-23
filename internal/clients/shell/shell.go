package shell

import "github.com/abiosoft/ishell/v2"

type Shell struct {
	shell *ishell.Shell
}

func NewShell() *Shell {
	return &Shell{
		shell: ishell.New(),
	}
}

func (s *Shell) HandleMultipleInput(message string, options []string) int {
	return s.shell.Actions.MultiChoice(options, message)
}

func (s *Shell) HandleCHeckBoxInput(message string, options []string) []int {
	return s.shell.Actions.Checklist(options, message, []int{0})
}
