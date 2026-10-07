package gui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/vedsatt/simple-os-emulator/internal/shell"
)

type terminalTheme struct{}

func (t *terminalTheme) Color(
	name fyne.ThemeColorName,
	variant fyne.ThemeVariant,
) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{
			R: 18,
			G: 18,
			B: 18,
			A: 255,
		}

	case theme.ColorNameForeground:
		return color.NRGBA{
			R: 235,
			G: 235,
			B: 235,
			A: 255,
		}

	case theme.ColorNameInputBackground:
		return color.Transparent

	case theme.ColorNameInputBorder:
		return color.Transparent
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t *terminalTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *terminalTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *terminalTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInputBorder:
		return 0

	case theme.SizeNameInnerPadding:
		return 0

	case theme.SizeNamePadding:
		return 2
	}

	return theme.DefaultTheme().Size(name)
}

func CreateWindow(a fyne.App, shell *shell.Shell) (fyne.Window, *widget.Entry) {
	a.Settings().SetTheme(&terminalTheme{})

	window := a.NewWindow("Simple OS Emulator — VFS: default")
	window.Resize(fyne.NewSize(800, 500))

	terminal := container.NewVBox()

	scroll := container.NewVScroll(terminal)

	window.SetContent(scroll)

	var firstInput *widget.Entry

	var createPrompt func() *widget.Entry

	createPrompt = func() *widget.Entry {
		prompt := newTerminalLabel("user@localhost:~$ ")

		input := widget.NewEntry()
		input.TextStyle = fyne.TextStyle{
			Monospace: true,
		}

		inputRow := container.NewBorder(
			nil,
			nil,
			prompt,
			nil,
			input,
		)

		terminal.Add(inputRow)

		input.OnSubmitted = func(command string) {
			terminal.Remove(inputRow)

			completedPrompt := newTerminalLabel("user@localhost:~$ ")
			completedCommand := newTerminalLabel(command)

			completedRow := container.NewBorder(
				nil,
				nil,
				completedPrompt,
				nil,
				completedCommand,
			)

			terminal.Add(completedRow)

			result, shouldExit := executeStub(command, shell)

			if shouldExit {
				window.Close()
				return
			}

			if result != "" {
				terminal.Add(
					newTerminalLabel(result),
				)
			}

			newInput := createPrompt()

			window.Canvas().Focus(newInput)

			scroll.ScrollToBottom()
		}

		return input
	}

	firstInput = createPrompt()

	return window, firstInput
}

func newTerminalLabel(text string) *widget.Label {
	label := widget.NewLabel(text)

	label.TextStyle = fyne.TextStyle{
		Monospace: true,
	}

	return label
}

func executeStub(command string, shell *shell.Shell) (string, bool) {
	command = strings.TrimSpace(command)

	return shell.Execute(command)
}
