package main

import (
	"embed"
	"encoding/json"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed secrets.json
var secrets embed.FS

type secretsStruct struct {
	Key string `json:"api_key"`
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "CurseForge Modpack Unpacker",
		Width:  400,
		Height: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		DisableResize:    true,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func APIKey() string {
	// Read the secrets.json file
	file, err := secrets.Open("secrets.json")
	if err != nil {
		return ""
	}
	defer file.Close()

	// Decode the JSON data
	var secret secretsStruct
	if err := json.NewDecoder(file).Decode(&secret); err != nil {
		return ""
	}

	return secret.Key
}
