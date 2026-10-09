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
	RevMode       bool

	CurrDir *vfs.Node
	PrevDir *vfs.Node
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
		CurrDir:    VFS.Root(),
	}

	return sh
}

func (s *Shell) Execute(command string) (string, bool) {
	if s.RevMode {
		return s.revCmd(nil, command), false
	}

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
	case "pwd":
		return s.pwdCmd(args[1:]), false
	case "whoami":
		return s.whoamiCmd(args[1:]), false
	case "rev":
		return s.revCmd(args[1:], ""), false
	case "exit":
		return "", true
	}

	return fmt.Sprintf("bash: %s: command not found", cmd), false
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

func (s *Shell) currentPath() string {
	if s.CurrDir == s.VFS.Root() {
		return "/"
	}

	parts := make([]string, 0)
	curr := s.CurrDir

	for curr != s.VFS.Root() {
		parts = append([]string{curr.Name}, parts...)
		curr = curr.Parent
	}

	path := "/" + strings.Join(parts, "/")

	if path == "/home/user" {
		return "~"
	}

	if strings.HasPrefix(path, "/home/user/") {
		return "~" + strings.TrimPrefix(path, "/home/user")
	}

	return path
}

func (s *Shell) PromptString() string {
	return fmt.Sprintf("%s:%s$", s.Prompt, s.currentPath())
}
