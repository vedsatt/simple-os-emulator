package shell

import "strings"

func (s *Shell) whoamiCmd(args []string) string {
	if len(args) != 0 {
		return "usage: whoami"
	}

	idx := strings.Index(s.Prompt, "@")
	return s.Prompt[:idx]
}
