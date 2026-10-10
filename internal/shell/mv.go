package shell

import (
	"fmt"
	"path"

	"github.com/vedsatt/simple-os-emulator/internal/vfs"
)

func (s *Shell) mvCmd(args []string) string {
	if len(args) != 2 {
		return "usage: mv source destination"
	}

	source := args[0]
	dest := args[1]

	sourceNode, err := s.resolvePath(source)
	if err != nil {
		return fmt.Sprintf("mv: %s: %v", source, err)
	}

	if sourceNode == s.VFS.Root() {
		return "mv: cannot move '/': operation not permitted"
	}

	destNode, err := s.resolvePath(dest)
	if err != nil {
		return s.moveToNewPath(sourceNode, source, dest)
	}

	if destNode.Type == vfs.Dir {
		return s.moveIntoDir(sourceNode, destNode, source, dest)
	}

	return s.replaceFile(sourceNode, destNode, source, dest)
}

func (s *Shell) moveToNewPath(
	sourceNode *vfs.Node,
	source string,
	dest string,
) string {
	destParentPath := path.Dir(dest)
	newName := path.Base(dest)

	destParent, err := s.resolvePath(destParentPath)
	if err != nil {
		return fmt.Sprintf(
			"mv: %s: no such file or directory",
			destParentPath,
		)
	}

	if destParent.Type != vfs.Dir {
		return fmt.Sprintf("mv: %s: not a directory", destParentPath)
	}

	if sourceNode.Type == vfs.Dir && isSubdir(sourceNode, destParent) {
		return subdirError(source, dest)
	}

	delete(sourceNode.Parent.Childs, sourceNode.Name)

	sourceNode.Name = newName
	sourceNode.Parent = destParent
	destParent.Childs[newName] = sourceNode

	return ""
}

func (s *Shell) moveIntoDir(
	sourceNode *vfs.Node,
	destNode *vfs.Node,
	source string,
	dest string,
) string {
	if sourceNode.Type == vfs.Dir && isSubdir(sourceNode, destNode) {
		return subdirError(source, dest)
	}

	existing := destNode.Childs[sourceNode.Name]

	if existing != nil && existing != sourceNode {
		if existing.Type == vfs.Dir {
			return fmt.Sprintf(
				"mv: cannot overwrite directory '%s'",
				sourceNode.Name,
			)
		}

		if sourceNode.Type == vfs.Dir {
			return fmt.Sprintf(
				"mv: cannot overwrite non-directory '%s' with directory '%s'",
				existing.Name,
				source,
			)
		}

		delete(destNode.Childs, existing.Name)
	}

	delete(sourceNode.Parent.Childs, sourceNode.Name)

	sourceNode.Parent = destNode
	destNode.Childs[sourceNode.Name] = sourceNode

	return ""
}

func (s *Shell) replaceFile(
	sourceNode *vfs.Node,
	destNode *vfs.Node,
	source string,
	dest string,
) string {
	if sourceNode.Type == vfs.Dir {
		return fmt.Sprintf(
			"mv: cannot overwrite non-directory '%s' with directory '%s'",
			dest,
			source,
		)
	}

	if sourceNode == destNode {
		return ""
	}

	sourceParent := sourceNode.Parent
	destParent := destNode.Parent
	newName := destNode.Name

	delete(sourceParent.Childs, sourceNode.Name)
	delete(destParent.Childs, destNode.Name)

	sourceNode.Name = newName
	sourceNode.Parent = destParent
	destParent.Childs[newName] = sourceNode

	return ""
}

func isSubdir(source *vfs.Node, dest *vfs.Node) bool {
	curr := dest

	for curr != nil {
		if curr == source {
			return true
		}

		curr = curr.Parent
	}

	return false
}

func subdirError(source string, dest string) string {
	return fmt.Sprintf(
		"mv: cannot move '%s' to a subdirectory of itself, '%s'",
		source,
		dest,
	)
}
