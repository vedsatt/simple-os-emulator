package shell

import (
	"fmt"
	"strings"

	"github.com/vedsatt/simple-os-emulator/internal/vfs"
)

func (s *Shell) cdCmd(args []string) string {
	if len(args) > 1 {
		return "cd: too many arguments"
	}

	if len(args) == 1 && args[0] == "-" {
		if s.PrevDir == nil {
			return "cd: OLDPWD not set"
		}

		s.CurrDir, s.PrevDir = s.PrevDir, s.CurrDir
		return s.currentPath()
	}

	path := ""

	if len(args) == 1 {
		path = args[0]
	}

	dir, err := s.resolvePath(path)
	if err != nil {
		return fmt.Sprintf("cd: %s: no such directory", path)
	}

	if dir.Type == vfs.File {
		return fmt.Sprintf("cd: %s: not a directory", path)
	}

	s.CurrDir, s.PrevDir = dir, s.CurrDir
	return ""
}

func (s *Shell) resolvePath(path string) (*vfs.Node, error) {
	if path == "~" || path == "" {
		path = "/home/user"
	} else if strings.HasPrefix(path, "~/") {
		path = "/home/user" + path[1:]
	}

	pathParts := strings.Split(path, "/")
	absolute := strings.HasPrefix(path, "/")
	var currDir *vfs.Node

	if absolute {
		currDir = s.VFS.Root()
		pathParts = pathParts[1:]
	} else {
		currDir = s.CurrDir
	}

	err := fmt.Errorf("no such file or directory")
	for i := range pathParts {
		if pathParts[i] == "" || pathParts[i] == "." {
			continue
		}

		if pathParts[i] == ".." {
			if currDir != s.VFS.Root() {
				currDir = currDir.Parent
			}

		} else {
			if currDir.Type == vfs.File {
				return nil, err
			}

			next, found := currDir.Childs[pathParts[i]]

			if !found {
				return nil, err
			}
			currDir = next
		}
	}

	return currDir, nil
}
