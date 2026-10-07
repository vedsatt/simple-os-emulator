package vfs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vedsatt/simple-os-emulator/internal/config"
)

type VFS struct {
	root *Node
}

func InitVFS(cfg *config.Config) (*VFS, error) {
	path := filepath.Join("data/vfs", cfg.VFSPath)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("err with opening vfs file: %w", err)
	}
	defer file.Close()

	var root *Node
	if err := json.NewDecoder(file).Decode(&root); err != nil {
		return nil, fmt.Errorf("err with decoding vfs file: %w", err)
	}

	return &VFS{root: root}, nil
}
