package gui

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed fonts/SF-Mono-Regular.otf
var sfMonoRegularBytes []byte

//go:embed fonts/SF-Mono-Semibold.otf
var sfMonoSemiboldBytes []byte

var sfMonoRegular = fyne.NewStaticResource(
	"SFMono-Regular.otf",
	sfMonoRegularBytes,
)

var sfMonoSemibold = fyne.NewStaticResource(
	"SFMono-Semibold.otf",
	sfMonoSemiboldBytes,
)
