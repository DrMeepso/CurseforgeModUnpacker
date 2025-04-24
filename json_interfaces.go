package main

// Root structure
type ModpackManifest struct {
	Minecraft       Minecraft `json:"minecraft"`
	ManifestType    string    `json:"manifestType"`
	ManifestVersion int       `json:"manifestVersion"`
	Name            string    `json:"name"`
	Version         string    `json:"version"`
	Author          string    `json:"author"`
	Files           []ModFile `json:"files"`
	Overrides       string    `json:"overrides"`
}

// Nested structure for "minecraft"
type Minecraft struct {
	Version        string      `json:"version"`
	ModLoaders     []ModLoader `json:"modLoaders"`
	RecommendedRam int         `json:"recommendedRam"`
}

// Nested structure for each mod loader
type ModLoader struct {
	ID      string `json:"id"`
	Primary bool   `json:"primary"`
}

// Structure for each file in "files"
type ModFile struct {
	ProjectID int  `json:"projectID"`
	FileID    int  `json:"fileID"`
	Required  bool `json:"required"`
}
