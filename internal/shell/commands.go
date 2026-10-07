package shell

import (
	"fmt"
	"strings"
)

func (s *Shell) lsCmd(args []string) string {
	stdoutArgs := strings.Join(args, " ")
	return fmt.Sprintf("ls %s", stdoutArgs)
}

func (s *Shell) cdCmd(args []string) string {
	stdoutArgs := strings.Join(args, " ")
	return fmt.Sprintf("cd %s", stdoutArgs)
}
