package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TerminalEntry struct {
	widget.Entry

	focused bool
}

func NewTerminalEntry() *TerminalEntry {
	entry := &TerminalEntry{}

	entry.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)

	entry.ExtendBaseWidget(entry)

	return entry
}

func (e *TerminalEntry) FocusGained() {
	e.focused = true
	e.Entry.FocusGained()
}

func (e *TerminalEntry) FocusLost() {
	e.focused = false
	e.Entry.FocusLost()
}

func (e *TerminalEntry) CreateRenderer() fyne.WidgetRenderer {
	baseRenderer := e.Entry.CreateRenderer()

	e.Entry.ExtendBaseWidget(e)

	blockCursor := canvas.NewRectangle(
		color.NRGBA{
			R: 235,
			G: 235,
			B: 235,
			A: 180,
		},
	)

	renderer := &terminalEntryRenderer{
		entry:       e,
		base:        baseRenderer,
		blockCursor: blockCursor,
	}

	renderer.updateCursor()

	return renderer
}

type terminalEntryRenderer struct {
	entry       *TerminalEntry
	base        fyne.WidgetRenderer
	blockCursor *canvas.Rectangle
}

func (r *terminalEntryRenderer) Layout(size fyne.Size) {
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

	objects := make([]fyne.CanvasObject, 0, len(baseObjects)+1)
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

	pos := r.entry.CursorPosition()

	textSize := r.entry.Theme().Size(theme.SizeNameText)

	r.blockCursor.Move(
		fyne.NewPos(
			pos.X,
			pos.Y,
		),
	)

	r.blockCursor.Resize(
		fyne.NewSize(
			8,
			textSize+4,
		),
	)
}
