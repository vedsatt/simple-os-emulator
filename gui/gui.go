package gui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/vedsatt/simple-os-emulator/internal/shell"
)

type terminalTheme struct{}

type terminalState struct {
	window   fyne.Window
	shell    *shell.Shell
	terminal *TerminalEntry

	history      []string
	historyIndex int
}

func (t *terminalTheme) Color(
	name fyne.ThemeColorName,
	variant fyne.ThemeVariant,
) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case theme.ColorNameInputBackground:
		return color.Transparent
	case theme.ColorNameInputBorder:
		return color.Transparent
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t *terminalTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		if style.Bold {
			return sfMonoSemibold
		}

		return sfMonoRegular
	}

	return theme.DefaultTheme().Font(style)
}

func (t *terminalTheme) Icon(
	name fyne.ThemeIconName,
) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *terminalTheme) Size(
	name fyne.ThemeSizeName,
) float32 {
	switch name {
	case theme.SizeNameText:
		return 12
	case theme.SizeNameInputBorder:
		return 0
	case theme.SizeNameInnerPadding:
		return 0
	case theme.SizeNamePadding:
		return 2
	}

	return theme.DefaultTheme().Size(name)
}

func CreateWindow(
	a fyne.App,
	sh *shell.Shell,
) (fyne.Window, *TerminalEntry) {
	a.Settings().SetTheme(&terminalTheme{})

	state := createTerminalState(a, sh)

	setupTerminalCallbacks(state)
	renderInitialTerminal(state)

	return state.window, state.terminal
}

func createTerminalState(
	a fyne.App,
	sh *shell.Shell,
) *terminalState {
	window := a.NewWindow(
		sh.Prompt + " — VFS: " + sh.VFS.Name,
	)

	window.Resize(
		fyne.NewSize(550, 400),
	)

	terminal := NewTerminalEntry()
	terminal.TextStyle = terminalTextStyle()

	window.SetContent(terminal)

	return &terminalState{
		window:       window,
		shell:        sh,
		terminal:     terminal,
		history:      make([]string, 0),
		historyIndex: 0,
	}
}

func setupTerminalCallbacks(state *terminalState) {
	state.terminal.OnSubmit = func(command string) {
		submitCommand(state, command)
	}

	state.terminal.OnUp = func() {
		historyUp(state)
	}

	state.terminal.OnDown = func() {
		historyDown(state)
	}

	state.terminal.OnCtrlD = func() {
		handleCtrlD(state)
	}

	state.terminal.OnPaste = func(text string) {
		runPastedCommands(state, text)
	}
}

func renderInitialTerminal(state *terminalState) {
	transcript := buildScriptTranscript(state)

	prompt := ""

	if !state.shell.ScriptShouldExit() {
		prompt = currentPrompt(state.shell)
	}

	state.terminal.LoadTranscript(
		transcript,
		prompt,
	)
}

func buildScriptTranscript(
	state *terminalState,
) string {
	results := state.shell.GetScriptResult()

	if len(results) == 0 {
		return ""
	}

	if results[0].Error() != nil {
		return results[0].Error().Error()
	}

	var builder strings.Builder

	for i, result := range results {
		if i > 0 {
			builder.WriteByte('\n')
		}

		builder.WriteString(state.shell.PromptString())
		builder.WriteByte(' ')
		builder.WriteString(result.Command())

		if result.String() != "" {
			builder.WriteByte('\n')
			builder.WriteString(result.String())
		}
	}

	return builder.String()
}

func submitCommand(
	state *terminalState,
	command string,
) bool {
	wasRevMode := state.shell.RevMode

	if !wasRevMode && strings.TrimSpace(command) != "" {
		addToHistory(state, command)
	}

	result, shouldExit := state.shell.Execute(command)

	if shouldExit {
		state.window.Close()
		return true
	}

	state.terminal.CommitCommand(
		result,
		currentPrompt(state.shell),
	)

	return false
}

func runPastedCommands(
	state *terminalState,
	text string,
) {
	text = state.terminal.CurrentCommand() + text

	for _, command := range strings.Split(text, "\n") {
		if strings.TrimSpace(command) == "" {
			continue
		}

		state.terminal.SetCurrentCommand(command)

		if submitCommand(state, command) {
			return
		}
	}
}

func handleCtrlD(state *terminalState) {
	if !state.shell.RevMode {
		return
	}

	state.shell.RevMode = false

	state.terminal.ResetCurrentLine(
		currentPrompt(state.shell),
	)
}

func historyUp(state *terminalState) {
	if len(state.history) == 0 {
		return
	}

	if state.historyIndex > 0 {
		state.historyIndex--
	}

	state.terminal.SetCurrentCommand(
		state.history[state.historyIndex],
	)
}

func historyDown(state *terminalState) {
	if len(state.history) == 0 {
		return
	}

	if state.historyIndex < len(state.history)-1 {
		state.historyIndex++

		state.terminal.SetCurrentCommand(
			state.history[state.historyIndex],
		)

		return
	}

	state.historyIndex = len(state.history)
	state.terminal.SetCurrentCommand("")
}

func addToHistory(
	state *terminalState,
	command string,
) {
	state.history = append(
		state.history,
		command,
	)

	state.historyIndex = len(state.history)
}

func currentPrompt(
	sh *shell.Shell,
) string {
	if sh.RevMode {
		return ""
	}

	return sh.PromptString() + " "
}

func terminalTextStyle() fyne.TextStyle {
	return fyne.TextStyle{
		Monospace: true,
		Bold:      true,
	}
}
