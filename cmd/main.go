package main

import (
	"log"

	"fyne.io/fyne/v2/app"

	"github.com/vedsatt/simple-os-emulator/gui"
	"github.com/vedsatt/simple-os-emulator/internal/config"
	"github.com/vedsatt/simple-os-emulator/internal/shell"
	"github.com/vedsatt/simple-os-emulator/internal/vfs"
)

func main() {
	cfg := config.GetConfig()

	vfs, err := vfs.InitVFS(cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	shell := shell.InitShell(vfs, cfg)

	a := app.New()

	window, input := gui.CreateWindow(a, shell)

	window.Show()

	window.Canvas().Focus(input)

	a.Run()
}
