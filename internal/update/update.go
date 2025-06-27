package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	githubAPI      = "https://api.github.com/repos/Fascinax/turbotilt/releases/latest"
	updateCheckTTL = 24 * time.Hour
)

// Release represents GitHub release information
type Release struct {
	TagName string    `json:"tag_name"`
	Assets  []Asset   `json:"assets"`
	URL     string    `json:"html_url"`
	Date    time.Time `json:"published_at"`
}

// Asset represents a release asset
type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int    `json:"size"`
}

// CheckForUpdates checks for newer versions of the app
func CheckForUpdates(currentVersion string) (*Release, bool) {
	// Skip update check if running from dev build
	if currentVersion == "dev" || currentVersion == "" {
		return nil, false
	}

	// Only check for updates once a day
	if !shouldCheckForUpdate() {
		return nil, false
	}

	// Remove 'v' prefix if present for comparison
	currentVersion = strings.TrimPrefix(currentVersion, "v")

	// Fetch the latest release info
	release, err := getLatestRelease()
	if err != nil {
		return nil, false
	}

	// Update the timestamp of the last check
	updateLastCheckTime()

	// Compare versions (simplistic version comparison)
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	return release, isNewer(currentVersion, latestVersion)
}

// shouldCheckForUpdate determines if we should check for updates
func shouldCheckForUpdate() bool {
	timestampFile := getUpdateTimestampPath()
	info, err := os.Stat(timestampFile)
	if err != nil {
		return true // File doesn't exist or error, check for update
	}

	// Check if the TTL has expired
	return time.Since(info.ModTime()) > updateCheckTTL
}

// updateLastCheckTime updates the timestamp of the last update check
func updateLastCheckTime() {
	timestampFile := getUpdateTimestampPath()
	os.MkdirAll(filepath.Dir(timestampFile), 0755)
	f, err := os.Create(timestampFile)
	if err != nil {
		return
	}
	defer f.Close()
}

// getUpdateTimestampPath returns the path to the update timestamp file
func getUpdateTimestampPath() string {
	var configDir string
	switch runtime.GOOS {
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), "turbotilt")
	case "darwin":
		configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "turbotilt")
	default: // Linux and others
		configDir = filepath.Join(os.Getenv("HOME"), ".config", "turbotilt")
	}
	
	return filepath.Join(configDir, "update_timestamp")
}

// getLatestRelease fetches the latest release info from GitHub API
func getLatestRelease() (*Release, error) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(githubAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch latest release: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release Release
	if err := json.Unmarshal(data, &release); err != nil {
		return nil, err
	}

	return &release, nil
}

// isNewer checks if latestVersion is newer than currentVersion
func isNewer(currentVersion, latestVersion string) bool {
	// Simple string comparison (assumes semver format like "1.2.3")
	return latestVersion > currentVersion
}

// GetUpdateAssetURL returns the download URL for the appropriate update asset
func GetUpdateAssetURL(release *Release) string {
	platform := runtime.GOOS
	arch := runtime.GOARCH

	// Convert architecture names
	if arch == "386" {
		arch = "x86"
	} else if arch == "amd64" {
		arch = "x64"
	}

	for _, asset := range release.Assets {
		// Match asset name with platform and architecture
		if strings.Contains(asset.Name, platform) && strings.Contains(asset.Name, arch) {
			return asset.DownloadURL
		}
	}

	// If we can't find a specific match, return the release URL
	return release.URL
}
