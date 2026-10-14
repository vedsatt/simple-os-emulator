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
		log.Fatal(err.Error())
	}

	log.Printf("VFS path: %s", cfg.VFSPath)
	log.Printf("Prompt: %s", cfg.Prompt)

	if cfg.Script != "" {
		log.Printf("Script: %s", cfg.Script)
	} else {
		log.Printf("Script: not set")
	}

	shell := shell.InitShell(vfs, cfg)

	if cfg.Script != "" {
		shell.ExecuteScript()
	}

	a := app.New()

	window, input := gui.CreateWindow(a, shell)

	window.Show()

	if shell.ScriptShouldExit() {
		window.Close()
		return
	}

	window.Canvas().Focus(input)

	a.Run()
}
