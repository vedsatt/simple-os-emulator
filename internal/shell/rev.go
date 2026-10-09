package shell

import (
	"fmt"
	"strings"

	"github.com/vedsatt/simple-os-emulator/internal/vfs"
)

func (s *Shell) revCmd(args []string, command string) string {
	if len(args) == 0 && !s.RevMode {
		s.RevMode = true

		return ""
	}

	builder := strings.Builder{}
	if len(args) != 0 {
		for i := range args {
			if i > 0 {
				builder.WriteByte('\n')
			}

			file, err := s.resolvePath(args[i])

			if err != nil {
				builder.WriteString("rev: " + args[i] + ": " + err.Error())
				continue
			}

			revFile, err := processFile(file)
			if err != nil {
				builder.WriteString("rev: " + args[i] + ": " + err.Error())
				continue
			}

			builder.WriteString(revFile)
		}

		return builder.String()
	}

	if len(args) == 0 && s.RevMode {
		reversed := []rune(command)
		l, r := 0, len(reversed)-1
		for l < r {
			reversed[l], reversed[r] = reversed[r], reversed[l]
			l++
			r--
		}

		return string(reversed)
	}

	return "rev: undefined"
}

func processFile(file *vfs.Node) (string, error) {
	if file.Type == vfs.Dir {
		return "", fmt.Errorf("is a directory")
	}

	fileLines := strings.Split(file.Content, "\n")
	reversed := strings.Builder{}
	for i := range fileLines {
		runeLine := []rune(fileLines[i])

		l, r := 0, len(runeLine)-1
		for l < r {
			runeLine[l], runeLine[r] = runeLine[r], runeLine[l]
			l++
			r--
		}

		reversed.WriteString(string(runeLine))

		if i < len(fileLines)-1 {
			reversed.WriteByte('\n')
		}
	}

	return reversed.String(), nil
}
