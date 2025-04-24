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

// Structure for the Curseforge API response
type CurseForgeMod struct {
	Data struct {
		ID                   int    `json:"id"`
		GameID               int    `json:"gameId"`
		Name                 string `json:"name"`
		Slug                 string `json:"slug"`
		Summary              string `json:"summary"`
		Status               int    `json:"status"`
		DownloadCount        int    `json:"downloadCount"`
		IsFeatured           bool   `json:"isFeatured"`
		ClassID              int    `json:"classId"`
		WebsiteURL           string `json:"websiteUrl"`
		DateCreated          string `json:"dateCreated"`
		DateModified         string `json:"dateModified"`
		DateReleased         string `json:"dateReleased"`
		AllowModDistribution bool   `json:"allowModDistribution"`
		GamePopularityRank   int    `json:"gamePopularityRank"`
		IsAvailable          bool   `json:"isAvailable"`
		HasCommentsEnabled   bool   `json:"hasCommentsEnabled"`
		ThumbsUpCount        int    `json:"thumbsUpCount"`

		Links struct {
			WebsiteURL string  `json:"websiteUrl"`
			WikiURL    string  `json:"wikiUrl"`
			IssuesURL  *string `json:"issuesUrl"`
			SourceURL  *string `json:"sourceUrl"`
		} `json:"links"`

		Categories []struct {
			ID               int    `json:"id"`
			GameID           int    `json:"gameId"`
			Name             string `json:"name"`
			Slug             string `json:"slug"`
			URL              string `json:"url"`
			IconURL          string `json:"iconUrl"`
			DateModified     string `json:"dateModified"`
			IsClass          bool   `json:"isClass"`
			ClassID          int    `json:"classId"`
			ParentCategoryID int    `json:"parentCategoryId"`
		} `json:"categories"`

		Authors []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			URL       string `json:"url"`
			AvatarURL string `json:"avatarUrl"`
		} `json:"authors"`

		Logo struct {
			ID           int    `json:"id"`
			ModID        int    `json:"modId"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			ThumbnailURL string `json:"thumbnailUrl"`
			URL          string `json:"url"`
		} `json:"logo"`

		LatestFiles []struct {
			ID            int      `json:"id"`
			GameID        int      `json:"gameId"`
			ModID         int      `json:"modId"`
			IsAvailable   bool     `json:"isAvailable"`
			DisplayName   string   `json:"displayName"`
			FileName      string   `json:"fileName"`
			ReleaseType   int      `json:"releaseType"`
			FileStatus    int      `json:"fileStatus"`
			FileDate      string   `json:"fileDate"`
			FileLength    int64    `json:"fileLength"`
			DownloadCount int      `json:"downloadCount"`
			DownloadURL   string   `json:"downloadUrl"`
			GameVersions  []string `json:"gameVersions"`
			Dependencies  []struct {
				ModID        int `json:"modId"`
				RelationType int `json:"relationType"`
			} `json:"dependencies"`
			Modules []struct {
				Name        string `json:"name"`
				Fingerprint uint32 `json:"fingerprint"`
			} `json:"modules"`
		} `json:"latestFiles"`

		LatestFilesIndexes []struct {
			GameVersion       string `json:"gameVersion"`
			FileID            int    `json:"fileId"`
			Filename          string `json:"filename"`
			ReleaseType       int    `json:"releaseType"`
			GameVersionTypeID int    `json:"gameVersionTypeId"`
			ModLoader         int    `json:"modLoader"`
		} `json:"latestFilesIndexes"`
	} `json:"data"`
}
