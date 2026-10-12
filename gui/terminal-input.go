package gui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TerminalEntry struct {
	widget.Entry

	focused      bool
	commandStart int

	OnSubmit func(string)
	OnCtrlD  func()
	OnUp     func()
	OnDown   func()
	OnPaste  func(string)
}

type terminalEntryRenderer struct {
	entry       *TerminalEntry
	base        fyne.WidgetRenderer
	blockCursor *canvas.Rectangle
}

func NewTerminalEntry() *TerminalEntry {
	entry := &TerminalEntry{}

	entry.MultiLine = true
	entry.Wrapping = fyne.TextWrapOff

	entry.ExtendBaseWidget(entry)

	return entry
}

func (e *TerminalEntry) LoadTranscript(
	transcript string,
	prompt string,
) {
	text := transcript

	if text != "" && prompt != "" {
		text += "\n"
	}

	text += prompt

	e.SetText(text)

	e.commandStart = len(
		[]rune(text),
	)

	e.moveCursorToEnd()
}

func (e *TerminalEntry) CurrentCommand() string {
	runes := []rune(e.Text)

	if e.commandStart >= len(runes) {
		return ""
	}

	return string(
		runes[e.commandStart:],
	)
}

func (e *TerminalEntry) SetCurrentCommand(
	command string,
) {
	runes := []rune(e.Text)

	prefix := string(
		runes[:e.commandStart],
	)

	e.SetText(prefix + command)
	e.moveCursorToEnd()
}

func (e *TerminalEntry) CommitCommand(
	output string,
	prompt string,
) {
	text := e.Text

	if output != "" {
		text += "\n" + output
	}

	text += "\n" + prompt

	e.SetText(text)

	e.commandStart = len(
		[]rune(text),
	)

	e.moveCursorToEnd()
}

func (e *TerminalEntry) ResetCurrentLine(
	prompt string,
) {
	runes := []rune(e.Text)

	text := string(
		runes[:e.commandStart],
	)

	text += prompt

	e.SetText(text)

	e.commandStart = len(
		[]rune(text),
	)

	e.moveCursorToEnd()
}

func (e *TerminalEntry) FocusGained() {
	e.focused = true

	e.Entry.FocusGained()
	e.Refresh()
}

func (e *TerminalEntry) FocusLost() {
	e.focused = false

	e.Entry.FocusLost()
	e.Refresh()
}

func (e *TerminalEntry) TypedRune(r rune) {
	e.insertText(string(r))
}

func (e *TerminalEntry) TypedKey(
	event *fyne.KeyEvent,
) {
	switch event.Name {
	case fyne.KeyReturn:
		if e.OnSubmit != nil {
			e.OnSubmit(
				e.CurrentCommand(),
			)
		}

	case fyne.KeyBackspace:
		e.backspace()

	case fyne.KeyDelete:
		e.deleteRune()

	case fyne.KeyUp:
		if e.OnUp != nil {
			e.OnUp()
		}

	case fyne.KeyDown:
		if e.OnDown != nil {
			e.OnDown()
		}

	case fyne.KeyHome:
		e.setCursorOffset(
			e.commandStart,
		)

	case fyne.KeyEnd:
		e.moveCursorToEnd()

	default:
		e.Entry.TypedKey(event)
		e.protectHistoryCursor()
	}
}

func (e *TerminalEntry) TypedShortcut(
	shortcut fyne.Shortcut,
) {
	if e.handlePaste(shortcut) {
		return
	}

	if e.handleCtrlD(shortcut) {
		return
	}

	e.Entry.TypedShortcut(shortcut)
}

func (e *TerminalEntry) handlePaste(
	shortcut fyne.Shortcut,
) bool {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		return false
	}

	text := paste.Clipboard.Content()

	if strings.Contains(text, "\n") {
		if e.OnPaste != nil {
			e.OnPaste(text)
		}

		return true
	}

	e.insertText(text)

	return true
}

func (e *TerminalEntry) handleCtrlD(
	shortcut fyne.Shortcut,
) bool {
	custom, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		return false
	}

	if custom.KeyName != fyne.KeyD {
		return false
	}

	if custom.Modifier != fyne.KeyModifierControl {
		return false
	}

	if e.OnCtrlD != nil {
		e.OnCtrlD()
	}

	return true
}

func (e *TerminalEntry) insertText(
	text string,
) {
	runes := []rune(e.Text)

	offset := e.CursorTextOffset()

	if offset < e.commandStart {
		offset = len(runes)
	}

	insert := []rune(text)

	result := make(
		[]rune,
		0,
		len(runes)+len(insert),
	)

	result = append(
		result,
		runes[:offset]...,
	)

	result = append(
		result,
		insert...,
	)

	result = append(
		result,
		runes[offset:]...,
	)

	e.SetText(
		string(result),
	)

	e.setCursorOffset(
		offset + len(insert),
	)
}

func (e *TerminalEntry) backspace() {
	runes := []rune(e.Text)

	offset := e.CursorTextOffset()

	if offset <= e.commandStart {
		return
	}

	result := append(
		[]rune{},
		runes[:offset-1]...,
	)

	result = append(
		result,
		runes[offset:]...,
	)

	e.SetText(
		string(result),
	)

	e.setCursorOffset(
		offset - 1,
	)
}

func (e *TerminalEntry) deleteRune() {
	runes := []rune(e.Text)

	offset := e.CursorTextOffset()

	if offset < e.commandStart {
		return
	}

	if offset >= len(runes) {
		return
	}

	result := append(
		[]rune{},
		runes[:offset]...,
	)

	result = append(
		result,
		runes[offset+1:]...,
	)

	e.SetText(
		string(result),
	)

	e.setCursorOffset(offset)
}

func (e *TerminalEntry) protectHistoryCursor() {
	if e.CursorTextOffset() >= e.commandStart {
		return
	}

	e.setCursorOffset(
		e.commandStart,
	)
}

func (e *TerminalEntry) moveCursorToEnd() {
	e.setCursorOffset(
		len([]rune(e.Text)),
	)
}

func (e *TerminalEntry) setCursorOffset(
	offset int,
) {
	runes := []rune(e.Text)

	if offset < 0 {
		offset = 0
	}

	if offset > len(runes) {
		offset = len(runes)
	}

	row := 0
	column := 0

	for _, r := range runes[:offset] {
		if r == '\n' {
			row++
			column = 0
			continue
		}

		column++
	}

	e.CursorRow = row
	e.CursorColumn = column

	e.Refresh()
}

func (e *TerminalEntry) CreateRenderer() fyne.WidgetRenderer {
	baseRenderer := e.Entry.CreateRenderer()

	e.Entry.ExtendBaseWidget(e)

	renderer := &terminalEntryRenderer{
		entry:       e,
		base:        baseRenderer,
		blockCursor: newBlockCursor(),
	}

	renderer.updateCursor()

	return renderer
}

func newBlockCursor() *canvas.Rectangle {
	return canvas.NewRectangle(
		color.NRGBA{
			R: 235,
			G: 235,
			B: 235,
			A: 180,
		},
	)
}

func (r *terminalEntryRenderer) Layout(
	size fyne.Size,
) {
	r.base.Layout(size)
	r.updateCursor()
}

func (r *terminalEntryRenderer) MinSize() fyne.Size {
	return r.base.MinSize()
}

func (r *terminalEntryRenderer) Refresh() {
	r.base.Refresh()
	r.updateCursor()

	canvas.Refresh(
		r.blockCursor,
	)
}

func (r *terminalEntryRenderer) Objects() []fyne.CanvasObject {
	baseObjects := r.base.Objects()

	objects := make(
		[]fyne.CanvasObject,
		0,
		len(baseObjects)+1,
	)

	objects = append(
		objects,
		baseObjects...,
	)

	objects = append(
		objects,
		r.blockCursor,
	)

	return objects
}

func (r *terminalEntryRenderer) Destroy() {
	r.base.Destroy()
}

func (r *terminalEntryRenderer) updateCursor() {
	if !r.entry.focused {
		r.blockCursor.Hide()
		return
	}

	r.blockCursor.Show()

	r.updateCursorPosition()
	r.updateCursorSize()
}

func (r *terminalEntryRenderer) updateCursorPosition() {
	position := r.entry.CursorPosition()

	r.blockCursor.Move(
		fyne.NewPos(
			position.X,
			position.Y,
		),
	)
}

func (r *terminalEntryRenderer) updateCursorSize() {
	textSize := r.entry.Theme().Size(
		theme.SizeNameText,
	)

	r.blockCursor.Resize(
		fyne.NewSize(
			8,
			textSize+4,
		),
	)
}
