package shell

import (
	"fmt"

	"github.com/vedsatt/simple-os-emulator/internal/config"
	"github.com/vedsatt/simple-os-emulator/internal/vfs"
)

type Shell struct {
	VFS *vfs.VFS
}

func InitShell(VFS *vfs.VFS, cfg *config.Config) *Shell {
	sh := &Shell{
		VFS: VFS,
	}

	return sh
}

func (s *Shell) Execute(req string) (string, bool) {
	args := s.Parse(req)

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
