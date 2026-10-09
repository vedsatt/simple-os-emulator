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
	Name string
}

func InitVFS(cfg *config.Config) (*VFS, error) {
	path := filepath.Join("data/vfs", cfg.VFSPath) + ".json"

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("err with opening vfs file: %w", err)
	}
	defer file.Close()

	var root *Node
	if err := json.NewDecoder(file).Decode(&root); err != nil {
		return nil, fmt.Errorf("err with decoding vfs file: %w", err)
	}

	setParent(root, nil)

	return &VFS{root: root, Name: cfg.VFSPath}, nil
}

func setParent(node *Node, parent *Node) {
	if node == nil {
		return
	}

	node.Parent = parent

	for _, child := range node.Childs {
		setParent(child, node)
	}
}

func (v *VFS) Root() *Node {
	return v.root
}
