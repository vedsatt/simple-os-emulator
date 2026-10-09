package shell

func (s *Shell) pwdCmd(args []string) string {
	if len(args) > 0 {
		return "pwd: too many arguments"
	}

	return s.currentPath()
}
