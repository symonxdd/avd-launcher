package main

import (
	"avd-launcher/app"
	"avd-launcher/app/manager"
	"avd-launcher/app/services"
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	myApp := app.NewApp()
	avdManager := manager.NewAvdManager()
	systemService := services.NewSystemService()
	updateService := services.NewUpdateService()

	err := wails.Run(&options.App{
		Title:  "AVD Launcher",
		Width:  950,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Frameless:        true,
		OnStartup: func(ctx context.Context) {
			app.SetContext(myApp, ctx)
			manager.SetContext(avdManager, ctx)
			services.SetContext(systemService, ctx)
		},
		LogLevel: logger.INFO,

		Bind: []interface{}{
			myApp,
			avdManager,
			systemService,
			updateService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
