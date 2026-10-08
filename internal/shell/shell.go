package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vedsatt/simple-os-emulator/internal/config"
	"github.com/vedsatt/simple-os-emulator/internal/vfs"
)

type Shell struct {
	VFS           *vfs.VFS
	Prompt        string
	scriptPath    string
	scriptResults []CommandResult
}

type CommandResult struct {
	cmd    string
	output string
	exit   bool
	err    error
}

func (c *CommandResult) Command() string {
	return c.cmd
}

func (c *CommandResult) String() string {
	return c.output
}

func (c *CommandResult) ShouldExit() bool {
	return c.exit
}

func (c *CommandResult) Error() error {
	return c.err
}

func (s *Shell) GetScriptResult() []CommandResult {
	return s.scriptResults
}

func InitShell(VFS *vfs.VFS, cfg *config.Config) *Shell {
	sh := &Shell{
		VFS:        VFS,
		Prompt:     cfg.Prompt,
		scriptPath: cfg.Script,
	}

	return sh
}

func (s *Shell) Execute(command string) (string, bool) {
	args := s.Parse(command)

	if len(args) == 0 {
		return "", false
	}

	cmd := args[0]
	switch cmd {
	case "ls":
		return s.lsCmd(args[1:]), false
	case "cd":
		return s.cdCmd(args[1:]), false
	case "exit":
		return "", true
	}

	return fmt.Sprintf("zsh: command not found: %s", cmd), false
}

func (s *Shell) ExecuteScript() {
	file, err := os.Open(s.scriptPath)
	if err != nil {
		res := CommandResult{
			cmd:    "",
			output: "",
			exit:   false,
			err:    fmt.Errorf("error with opening file: %w", err),
		}

		s.scriptResults = []CommandResult{res}
		return
	}
	defer file.Close()

	output := []CommandResult{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		cmd := strings.TrimSpace(scanner.Text())
		if cmd != "" && !strings.HasPrefix(cmd, "//") {
			exec, exit := s.Execute(cmd)
			res := CommandResult{
				cmd:    cmd,
				output: exec,
				exit:   exit,
			}

			output = append(output, res)

			if exit {
				break
			}
		}
	}

	s.scriptResults = output
}

func (s *Shell) ScriptShouldExit() bool {
	if len(s.scriptResults) == 0 {
		return false
	}

	return s.scriptResults[len(s.scriptResults)-1].ShouldExit()
}
