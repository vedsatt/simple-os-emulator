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

type terminalState struct {
	window     fyne.Window
	content    *fyne.Container
	scroll     *container.Scroll
	shell      *shell.Shell
	output     *widget.Label
	transcript string

	history      []string
	historyIndex int

	createPrompt func() *TerminalEntry
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

	renderScriptResults(state)

	if sh.ScriptShouldExit() {
		return state.window, nil
	}

	state.createPrompt = func() *TerminalEntry {
		return createPrompt(state)
	}

	firstInput := state.createPrompt()

	return state.window, firstInput
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

	output := newTerminalOutput()

	content := container.NewVBox(output)
	scroll := container.NewVScroll(content)

	window.SetContent(scroll)

	return &terminalState{
		window:       window,
		content:      content,
		scroll:       scroll,
		shell:        sh,
		output:       output,
		history:      make([]string, 0),
		historyIndex: 0,
	}
}

func newTerminalOutput() *widget.Label {
	output := widget.NewLabel("")

	output.TextStyle = terminalTextStyle()
	output.Selectable = true
	output.Wrapping = fyne.TextWrapOff

	return output
}

func renderScriptResults(state *terminalState) {
	results := state.shell.GetScriptResult()

	if len(results) == 0 {
		return
	}

	if results[0].Error() != nil {
		appendTranscript(
			state,
			results[0].Error().Error(),
		)
		return
	}

	for _, result := range results {
		appendCommand(
			state,
			result.Command(),
			false,
		)

		if result.String() != "" {
			appendTranscript(
				state,
				result.String(),
			)
		}
	}
}

func createPrompt(
	state *terminalState,
) *TerminalEntry {
	prompt := newTerminalLabel(
		currentPrompt(state.shell),
	)

	input := NewTerminalEntry()
	input.TextStyle = terminalTextStyle()

	inputRow := container.NewBorder(
		nil,
		nil,
		prompt,
		nil,
		input,
	)

	setupInputCallbacks(
		state,
		input,
		inputRow,
	)

	state.content.Add(inputRow)

	return input
}

func setupInputCallbacks(
	state *terminalState,
	input *TerminalEntry,
	inputRow *fyne.Container,
) {
	setupHistoryCallbacks(state, input)

	input.OnCtrlD = func() {
		handleCtrlD(state, inputRow)
	}

	input.OnSubmitted = func(command string) {
		runCommands(
			state,
			inputRow,
			command,
		)
	}

	input.OnPaste = func(text string) {
		runCommands(
			state,
			inputRow,
			input.Text+text,
		)
	}
}

func setupHistoryCallbacks(
	state *terminalState,
	input *TerminalEntry,
) {
	input.OnUp = func() {
		historyUp(state, input)
	}

	input.OnDown = func() {
		historyDown(state, input)
	}
}

func historyUp(
	state *terminalState,
	input *TerminalEntry,
) {
	if len(state.history) == 0 {
		return
	}

	if state.historyIndex > 0 {
		state.historyIndex--
	}

	input.SetText(
		state.history[state.historyIndex],
	)

	moveCursorToEnd(input)
}

func historyDown(
	state *terminalState,
	input *TerminalEntry,
) {
	if len(state.history) == 0 {
		return
	}

	if state.historyIndex < len(state.history)-1 {
		state.historyIndex++

		input.SetText(
			state.history[state.historyIndex],
		)
	} else {
		state.historyIndex = len(state.history)
		input.SetText("")
	}

	moveCursorToEnd(input)
}

func handleCtrlD(
	state *terminalState,
	inputRow *fyne.Container,
) {
	if !state.shell.RevMode {
		return
	}

	state.shell.RevMode = false

	state.content.Remove(inputRow)

	focusNewPrompt(state)
}

func runCommands(
	state *terminalState,
	inputRow *fyne.Container,
	text string,
) {
	state.content.Remove(inputRow)

	for _, command := range strings.Split(text, "\n") {
		if strings.TrimSpace(command) == "" {
			continue
		}

		if runCommand(state, command) {
			return
		}
	}

	focusNewPrompt(state)
}

func runCommand(
	state *terminalState,
	command string,
) bool {
	wasRevMode := state.shell.RevMode

	appendCommand(
		state,
		command,
		wasRevMode,
	)

	if !wasRevMode {
		addToHistory(state, command)
	}

	result, shouldExit := state.shell.Execute(command)

	if result != "" {
		appendTranscript(state, result)
	}

	if shouldExit {
		state.window.Close()
		return true
	}

	return false
}

func appendCommand(
	state *terminalState,
	command string,
	revMode bool,
) {
	prompt := ""

	if !revMode {
		prompt = state.shell.PromptString() + " "
	}

	appendTranscript(
		state,
		prompt+command,
	)
}

func appendTranscript(
	state *terminalState,
	text string,
) {
	if state.transcript != "" {
		state.transcript += "\n"
	}

	state.transcript += text

	state.output.SetText(
		state.transcript,
	)
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

func focusNewPrompt(
	state *terminalState,
) {
	input := state.createPrompt()

	state.window.Canvas().Focus(input)
	state.scroll.ScrollToBottom()
}

func currentPrompt(
	sh *shell.Shell,
) string {
	if sh.RevMode {
		return ""
	}

	return sh.PromptString() + " "
}

func moveCursorToEnd(
	input *TerminalEntry,
) {
	input.CursorColumn = len(
		[]rune(input.Text),
	)

	input.Refresh()
}

func terminalTextStyle() fyne.TextStyle {
	return fyne.TextStyle{
		Monospace: true,
		Bold:      true,
	}
}

func newTerminalLabel(
	text string,
) *widget.Label {
	label := widget.NewLabel(text)

	label.TextStyle = terminalTextStyle()

	return label
}
