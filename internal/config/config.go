package config

import "flag"

const (
	defaultVFSPath = "deep"
	defaultPrompt  = "user@localhost"
	defaultScript  = ""
)

type Config struct {
	VFSPath string
	Prompt  string
	Script  string
}

func GetConfig() *Config {
	cfg := &Config{}

	flag.StringVar(
		&cfg.VFSPath,
		"vfs",
		defaultVFSPath,
		"path to VFS file",
	)

	flag.StringVar(
		&cfg.Prompt,
		"prompt",
		defaultPrompt,
		"shell prompt",
	)

	flag.StringVar(
		&cfg.Script,
		"script",
		defaultScript,
		"path to startup script",
	)

	flag.Parse()

	return cfg
}
