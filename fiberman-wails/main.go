package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	apiHandler, app, err := newDesktopRuntime()
	if err != nil {
		log.Fatal(err)
	}

	if err := wails.Run(&options.App{
		Title:            "FiberMan",
		Width:            1440,
		Height:           960,
		MinWidth:         1180,
		MinHeight:        760,
		BackgroundColour: options.NewRGB(14, 18, 26),
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: apiHandler,
		},
		Linux: &linux.Options{
			ProgramName: "FiberMan",
		},
	}); err != nil {
		log.Fatal(err)
	}
}
