package shell

func (s *Shell) exitCmd(args []string) (string, bool) {
	if len(args) != 0 {
		return "usage: exit", false
	}

	return "", true
}
