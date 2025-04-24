package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/pkg/browser"
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

			if IncludeOverrides {
				// make a /mods folder in the output folder
				err := os.MkdirAll(outputFolder+"/mods", 0755)
				IfErrorReturn(err, "Error creating mods folder", a.ctx)
				PrintOutput("Creating mods folder", a.ctx)
			}

			BeginUnpack(mainifest.Files, outputFolder, a.ctx, zipReader)

		}
	}
	return true, nil
}

var IncludeOverrides bool = false

func (a *App) SetIncludeOverrides(value bool) {
	IncludeOverrides = value
	if IncludeOverrides {
		PrintOutput("Including overrides", a.ctx)
	} else {
		PrintOutput("Not including overrides", a.ctx)
	}
}

func BeginUnpack(mods []ModFile, outputDir string, ctx context.Context, zipReader *zip.Reader) {

	// create a slice of the modfile array so we can remove files from it
	UndownloadedMods := make([]ModFile, len(mods))
	ListMutex := sync.Mutex{}
	copy(UndownloadedMods, mods)

	UnpackWorker := func() {
		for len(UndownloadedMods) > 0 {
			// get the first mod in the list
			mod := UndownloadedMods[0]
			ListMutex.Lock()
			// remove the mod from the list
			UndownloadedMods = UndownloadedMods[1:]
			ListMutex.Unlock()

			// Querry the API for the mod's info
			url := "https://api.curseforge.com/v1/mods/" + fmt.Sprint(mod.ProjectID)

			req, err := http.NewRequest("GET", url, nil)
			IfErrorReturn(err, "Error creating request", ctx)
			req.Header.Set("x-api-key", APIKey())

			client := &http.Client{}
			resp, err := client.Do(req)
			IfErrorReturn(err, "Error making request", ctx)
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				PrintOutput("Error: "+fmt.Sprint(resp.StatusCode), ctx)
			}

			// read the response body
			json_response, err := ioutil.ReadAll(resp.Body)
			IfErrorReturn(err, "Error reading response body", ctx)
			// parse the json response
			var curseForgeMod CurseForgeMod
			if err := json.Unmarshal(json_response, &curseForgeMod); err != nil {
				PrintOutput("Error parsing json response", ctx)
				continue
			}
			PrintOutput("Downloading mod: "+curseForgeMod.Data.Name, ctx)

			// use a filestread to download the file
			// file url is: https://www.curseforge.com/api/v1/mods/${modID}/files/${fileID}/download

			downloadURL := "https://www.curseforge.com/api/v1/mods/" + fmt.Sprint(mod.ProjectID) + "/files/" + fmt.Sprint(mod.FileID) + "/download"
			req, err = http.NewRequest("GET", downloadURL, nil)
			IfErrorReturn(err, "Error creating request", ctx)

			req.Header.Set("x-api-key", APIKey())
			resp, err = client.Do(req)
			IfErrorReturn(err, "Error making request", ctx)

			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				PrintOutput("Error: "+fmt.Sprint(resp.StatusCode), ctx)
				continue
			}

			// write the file to disk
			filepath := outputDir + "/" + curseForgeMod.Data.Slug + ".jar"
			if IncludeOverrides {
				filepath = outputDir + "/mods/" + curseForgeMod.Data.Slug + ".jar"
			}
			out, err := os.Create(filepath)
			IfErrorReturn(err, "Error creating file for mod: "+curseForgeMod.Data.Name, ctx)
			defer out.Close()

			// write the response body to the file
			_, err = io.Copy(out, resp.Body)
			IfErrorReturn(err, "Error writing file for mod: "+curseForgeMod.Data.Name, ctx)

			finishedPercent := 100 - (len(UndownloadedMods) * 100 / len(mods))
			runtime.EventsEmit(ctx, "progress", finishedPercent)
		}
	}

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			PrintOutput("Spawning goroutine: "+fmt.Sprint(i), ctx)
			UnpackWorker()
		}()
	}
	wg.Wait()

	if !IncludeOverrides {
		runtime.EventsEmit(ctx, "finished", true)
		return
	}

	PrintOutput("Unpacking overrides...", ctx)

	// unzip the /overrides folder in the modpack zip to the output folder
	for _, f := range zipReader.File {

		filePathParts := strings.Split(f.Name, "/")
		if filePathParts[0] == "overrides" {

			// walk up the array making a folder if it doesn't exist
			for i := 1; i < len(filePathParts)-1; i++ {
				// get the slice of the path from 1 to i
				path := filePathParts[1 : i+1]
				err := os.MkdirAll(outputDir+"\\"+strings.Join(path, "\\"), 0755)
				IfErrorReturn(err, "Error creating folder: "+filePathParts[i], ctx)
			}
			// create the file
			path := filePathParts[1 : len(filePathParts)-1]
			outFile, err := os.Create(outputDir + "\\" + strings.Join(path, "\\") + "\\" + filePathParts[len(filePathParts)-1])
			IfErrorReturn(err, "Error creating file: "+filePathParts[len(filePathParts)-1], ctx)
			defer outFile.Close()

			// write the file to disk
			zipFile, err := f.Open()
			IfErrorReturn(err, "Error opening zip file: "+filePathParts[len(filePathParts)-1], ctx)
			defer zipFile.Close()
			_, err = io.Copy(outFile, zipFile)
			IfErrorReturn(err, "Error writing file: "+filePathParts[len(filePathParts)-1], ctx)
			PrintOutput("Unzipped: "+filePathParts[len(filePathParts)-1], ctx)

		}
	}

	PrintOutput("Unpacking finished!", ctx)

	runtime.EventsEmit(ctx, "finished", true)

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

func (a *App) OpenDonationPage() {
	// Open the donation page in the default browser
	err := browser.OpenURL("https://ko-fi.com/drmeepso")
	if err != nil {
		println("Error opening donation page:", err)
	}
	runtime.EventsEmit(a.ctx, "output", "Opening donation page...")

}
