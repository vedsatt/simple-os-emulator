package shell

import "strings"

func (s *Shell) Parse(cmd string) []string {
	return strings.Fields(cmd)
}
