package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main initializes the application, registers services, creates the main window,
// and runs the app, logging any error that occurs.
func main() {
	app := application.New(application.Options{
		Name:        "PokeFanLauncher",
		Description: "A launcher for fan-made Pokémon games",
		Services: []application.Service{
			application.NewService(&AppService{}),
			application.NewService(&ConfigService{}),
			application.NewService(&GameService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "PokeFanLauncher",
		Width:            1100,
		Height:           720,
		MinWidth:         800,
		MinHeight:        540,
		BackgroundColour: application.NewRGB(15, 17, 26),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
