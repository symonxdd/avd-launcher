package main

import (
	"avd-launcher/app"
	"avd-launcher/app/services"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	myApp := app.NewApp()
	updateService := services.NewUpdateService()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "AVD Launcher",
		Width:  950, // 1024
		Height: 600, // 768
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Frameless:        true,
		OnStartup:        myApp.Startup,
		LogLevel:         logger.INFO,

		// This allows the frontend to call methods from the backend
		Bind: []interface{}{
			myApp,
			updateService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
