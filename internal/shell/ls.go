package shell

import (
	"sort"
	"strings"

	"github.com/vedsatt/simple-os-emulator/internal/vfs"
)

type lsRes struct {
	path string
	node *vfs.Node
}

type dirRes struct {
	path   string
	stdout []string
}

type lsErr struct {
	path string
	err  error
}

type lsFiles struct {
	path string
}

func (s *Shell) lsCmd(args []string) string {
	lsDirs := make([]lsRes, 0)

	errs := make([]lsErr, 0)
	files := make([]lsFiles, 0)
	dirs := make([]dirRes, 0)

	// добавляем текущую ноду
	if len(args) == 0 {
		lsDirs = append(lsDirs, lsRes{path: "", node: s.CurrDir})
	} else {
		// для каждого пути ищем ноду и добавляем либо ошибку нахождения либо саму ноду
		for i := range len(args) {
			node, err := s.resolvePath(args[i])

			if err != nil {
				errs = append(errs, lsErr{
					path: args[i],
					err:  err,
				})
			} else {
				lsDirs = append(lsDirs, lsRes{
					path: args[i],
					node: node,
				})
			}
		}
	}

	// обрабатываем найденные пути
	for i := range lsDirs {
		out, isFile := getLsItems(lsDirs[i].node)
		if isFile {
			files = append(files, lsFiles{path: lsDirs[i].path})
		} else {
			dirs = append(dirs, dirRes{path: lsDirs[i].path, stdout: out})
		}
	}

	// сюда в форматирование закинем ошибки, файли и директории
	return formatLs(errs, files, dirs)
}

// если дир - ее содержимое, иначе файл
func getLsItems(node *vfs.Node) ([]string, bool) {
	if node.Type == vfs.File {
		return []string{node.Name}, true
	}

	out := make([]string, 0)
	for k := range node.Childs {
		out = append(out, k)
	}

	sort.Strings(out)

	return out, false
}

func formatLsItems(items []string) string {
	if len(items) == 0 {
		return ""
	}

	maxLen := 0
	for _, item := range items {
		if len(item) > maxLen {
			maxLen = len(item)
		}
	}

	columnWidth := maxLen + 3
	terminalWidth := 80

	cols := min(max(terminalWidth/columnWidth, 1), len(items))

	var builder strings.Builder

	for i, item := range items {
		builder.WriteString(item)

		if (i+1)%cols == 0 || i == len(items)-1 {
			if i != len(items)-1 {
				builder.WriteString("\n")
			}
			continue
		}

		spaces := columnWidth - len(item)
		builder.WriteString(strings.Repeat(" ", spaces))
	}

	return builder.String()
}

func formatLs(errs []lsErr, files []lsFiles, dirs []dirRes) string {
	builder := strings.Builder{}

	// ошибки
	for i := range errs {
		builder.WriteString("ls: ")
		builder.WriteString(errs[i].path)
		builder.WriteString(": ")
		builder.WriteString(errs[i].err.Error())

		if i < len(errs)-1 {
			builder.WriteString("\n")
		}
	}

	if len(errs) > 0 && (len(files) > 0 || len(dirs) > 0) {
		builder.WriteString("\n\n")
	}

	// файлы
	for i := range files {
		builder.WriteString(files[i].path)

		if i < len(files)-1 {
			builder.WriteString("   ")
		}
	}

	if len(files) > 0 && len(dirs) > 0 {
		builder.WriteString("\n\n")
	}

	// директории
	total := len(errs) + len(files) + len(dirs)

	for i := range dirs {
		if total > 1 {
			builder.WriteString(dirs[i].path)
			builder.WriteString(":\n")
		}

		builder.WriteString(formatLsItems(dirs[i].stdout))

		if i < len(dirs)-1 {
			builder.WriteString("\n\n")
		}
	}

	return builder.String()
}
