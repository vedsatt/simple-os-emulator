package vfs

type NodeType string

const (
	File NodeType = "file"
	Dir  NodeType = "dir"
)

type Node struct {
	Name    string
	Type    NodeType
	Content string
	Parent  *Node
	Childs  map[string]*Node
}
