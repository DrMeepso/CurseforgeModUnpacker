package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ZipFileDialog() (string, error) {
	// Open a file dialog to select a zip file

	options := runtime.OpenDialogOptions{
		Title: "Select File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Modpack Files", Pattern: "*.zip"},
		},
		DefaultDirectory: "",
	}

	filePath, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		return "", err
	}
	return filePath, nil

}

func (a *App) FolderDialog() (string, error) {
	// Open a folder dialog to select a folder

	options := runtime.OpenDialogOptions{
		Title:            "Select Folder",
		DefaultDirectory: "",
	}
	folderPath, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		return "", err
	}
	return folderPath, nil
}

func (a *App) VerifyModpackFile(filePath string) (string, error) {
	// Verify the modpack file
	// get the file, read the zip file and check if it contains the modpack.json file
	// if it does, return "ok"
	// if it doesn't, return an error

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return "Error Reading File", nil
	}
	defer file.Close() // make sure to close the file after reading

	// Check if the file is a zip file
	if file.Name() == "" {
		return "Error: File is not a zip file", fmt.Errorf("file is not a zip file")
	}

	// Check if the file contains the modpack.json file
	fileStats, err := file.Stat()
	if err != nil {
		return "Error: File is not a zip file", nil
	}
	zipReader, err := zip.NewReader(file, fileStats.Size())
	if err != nil {
		return "Error: File is not a zip file", nil
	}
	for _, f := range zipReader.File {
		println("\"" + f.Name + "\"")
		if f.Name == "manifest.json" {
			// read json file
			json_file, err := zipReader.Open(f.Name)
			if err != nil {
				return "Error: File is not a zip file", nil
			}
			defer json_file.Close() // make sure to close the file after reading

			// read the file
			var mainifest ModpackManifest
			if err := json.NewDecoder(json_file).Decode(&mainifest); err != nil {
				return "Error: File is not a modpack file", nil
			}

			println("Manifest Type: " + mainifest.ManifestType)
			println("Manifest Version: " + fmt.Sprint(mainifest.ManifestVersion))

			return "Info: " + mainifest.Name + " - " + mainifest.Version, nil
		}
	}

	return "Error: File is not a modpack file", nil
}

func (a *App) ShowErrorDialog(message string) {
	// Show an error dialog
	options := runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   "Error",
		Message: message,
	}
	runtime.MessageDialog(a.ctx, options)
}

func (a *App) RunUnpack(modfilePath string, outputFolder string) (bool, error) {

	// first tell the frontend to be on the downloading state
	runtime.EventsEmit(a.ctx, "stateChange", 1)

	// read the manifest.json file
	file, err := os.OpenFile(modfilePath, os.O_RDONLY, 0)
	IfErrorReturn(err, "Error Reading File", a.ctx)

	defer file.Close() // make sure to close the file after reading
	// Check if the file is a zip file
	if file.Name() == "" {
		return false, fmt.Errorf("file is not a zip file")
	}
	// Check if the file contains the modpack.json file
	fileStats, err := file.Stat()
	IfErrorReturn(err, "Error: File is not a zip file", a.ctx)

	zipReader, err := zip.NewReader(file, fileStats.Size())
	IfErrorReturn(err, "Error: File is not a zip file", a.ctx)

	// read the file
	for _, f := range zipReader.File {
		if f.Name == "manifest.json" {
			// read json file
			json_file, err := zipReader.Open(f.Name)
			IfErrorReturn(err, "Error: File is not a zip file", a.ctx)

			defer json_file.Close() // make sure to close the file after reading

			// read the file
			var mainifest ModpackManifest
			if err := json.NewDecoder(json_file).Decode(&mainifest); err != nil {
				return false, fmt.Errorf("error decoding json file")
			}

			PrintOutput("Modpack Name: "+mainifest.Name, a.ctx)
			PrintOutput("Modpack Version: "+mainifest.Version, a.ctx)
			PrintOutput("Modpack Author: "+mainifest.Author, a.ctx)

		}
	}
	return true, nil
}

func IfErrorReturn(err error, message string, ctx context.Context) bool {
	if err != nil {
		println(message)
		PrintOutput(message, ctx)
		return true
	}
	return false
}

func PrintOutput(message string, ctx context.Context) {
	println(message)
	runtime.EventsEmit(ctx, "output", message)
}
