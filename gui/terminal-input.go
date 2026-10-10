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

	focused bool

	OnCtrlD func()
	OnUp    func()
	OnDown  func()
	OnPaste func(string)
}

type terminalEntryRenderer struct {
	entry       *TerminalEntry
	base        fyne.WidgetRenderer
	blockCursor *canvas.Rectangle
}

func NewTerminalEntry() *TerminalEntry {
	entry := &TerminalEntry{}

	entry.Wrapping = fyne.TextWrap(
		fyne.TextTruncateClip,
	)

	entry.ExtendBaseWidget(entry)

	return entry
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

	if !strings.Contains(text, "\n") {
		return false
	}

	if e.OnPaste == nil {
		return false
	}

	e.OnPaste(text)

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

func (e *TerminalEntry) TypedKey(
	event *fyne.KeyEvent,
) {
	switch event.Name {
	case fyne.KeyUp:
		if e.OnUp != nil {
			e.OnUp()
			return
		}

	case fyne.KeyDown:
		if e.OnDown != nil {
			e.OnDown()
			return
		}
	}

	e.Entry.TypedKey(event)
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
	canvas.Refresh(r.blockCursor)
}

func (r *terminalEntryRenderer) Objects() []fyne.CanvasObject {
	baseObjects := r.base.Objects()

	objects := make(
		[]fyne.CanvasObject,
		0,
		len(baseObjects)+1,
	)

	objects = append(objects, baseObjects...)
	objects = append(objects, r.blockCursor)

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
