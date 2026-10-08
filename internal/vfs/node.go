package vfs

type NodeType string

const (
	File NodeType = "file"
	Dir  NodeType = "dir"
)

type Node struct {
	Name    string           `json:"name"`
	Type    NodeType         `json:"type"`
	Content string           `json:"content,omitempty"`
	Parent  *Node            `json:"-"`
	Childs  map[string]*Node `json:"children,omitempty"`
}
