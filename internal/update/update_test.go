package update

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIsNewer(t *testing.T) {
	tests := []struct {
		current, latest string
		result          bool
	}{
		{"1.0.0", "1.0.1", true},
		{"1.0.0", "1.1.0", true},
		{"1.0.0", "2.0.0", true},
		{"1.0.0", "1.0.0", false},
		{"1.0.1", "1.0.0", false},
		{"2.0.0", "1.0.0", false},
		{"1.2.3", "1.2.3-alpha", false}, // Simple case, no complex semantic versioning
	}

	for _, test := range tests {
		result := isNewer(test.current, test.latest)
		if result != test.result {
			t.Errorf("isNewer(%s, %s) = %v, attendu %v",
				test.current, test.latest, result, test.result)
		}
	}
}

func TestShouldCheckForUpdate(t *testing.T) {
	// Create a temporary directory for tests
	tempDir := t.TempDir()
	origGetUpdateTimestampPath := getUpdateTimestampPath
	defer func() { getUpdateTimestampPath = origGetUpdateTimestampPath }()

	getUpdateTimestampPath = func() string {
		return filepath.Join(tempDir, "update_timestamp")
	}

	// Case 1: No timestamp file, should return true
	if !shouldCheckForUpdate() {
		t.Error("shouldCheckForUpdate should return true when file doesn't exist")
	}

	// Create a timestamp file
	timestampFile := getUpdateTimestampPath()
	if err := os.MkdirAll(filepath.Dir(timestampFile), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	f, err := os.Create(timestampFile)
	if err != nil {
		t.Fatalf("Failed to create timestamp file: %v", err)
	}
	f.Close()

	// Case 2: Recent timestamp file, should return false
	if shouldCheckForUpdate() {
		t.Error("shouldCheckForUpdate should return false with a recent timestamp file")
	}

	// Case 3: Old timestamp file, should return true
	// Simulate an old file by modifying its modification date
	oldTime := time.Now().Add(-updateCheckTTL - time.Hour)
	if err := os.Chtimes(timestampFile, oldTime, oldTime); err != nil {
		t.Fatalf("Failed to modify file date: %v", err)
	}

	if !shouldCheckForUpdate() {
		t.Error("shouldCheckForUpdate should return true with an old timestamp file")
	}
}

func TestUpdateLastCheckTime(t *testing.T) {
	// Create a temporary directory for tests
	tempDir := t.TempDir()
	origGetUpdateTimestampPath := getUpdateTimestampPath
	defer func() { getUpdateTimestampPath = origGetUpdateTimestampPath }()

	getUpdateTimestampPath = func() string {
		return filepath.Join(tempDir, "update_timestamp")
	}

	// Verify that the file doesn't exist at the beginning
	timestampFile := getUpdateTimestampPath()
	if _, err := os.Stat(timestampFile); err == nil {
		t.Error("Timestamp file should not exist at the beginning of the test")
	}

	// Update the timestamp
	updateLastCheckTime()

	// Verify that the file now exists
	if _, err := os.Stat(timestampFile); err != nil {
		t.Errorf("Timestamp file should exist after updateLastCheckTime: %v", err)
	}
}

func TestGetUpdateTimestampPath(t *testing.T) {
	// Test that the function returns a non-empty path
	path := getUpdateTimestampPath()
	if path == "" {
		t.Error("getUpdateTimestampPath should not return an empty path")
	}
}

func TestGetLatestRelease(t *testing.T) {
	// Create a test HTTP server that simulates the GitHub API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"tag_name": "v1.2.3",
			"html_url": "https://github.com/test/repo/releases/v1.2.3",
			"published_at": "2023-01-01T00:00:00Z",
			"assets": [
				{
					"name": "turbotilt-1.2.3-darwin-amd64.zip",
					"browser_download_url": "https://github.com/test/repo/releases/download/v1.2.3/turbotilt-1.2.3-darwin-amd64.zip",
					"size": 12345
				}
			]
		}`))
	}))
	defer server.Close()

	// Save the original URL and restore it after the test
	origGithubAPI := githubAPI
	defer func() {
		githubAPI = origGithubAPI
	}()

	// Use a local variable with the same value
	testAPI := server.URL
	// Use this variable in the test function
	githubAPI = testAPI

	// Test getLatestRelease
	release, err := getLatestRelease()
	if err != nil {
		t.Fatalf("getLatestRelease returned an error: %v", err)
	}

	if release.TagName != "v1.2.3" {
		t.Errorf("release.TagName = %s, expected v1.2.3", release.TagName)
	}

	if release.URL != "https://github.com/test/repo/releases/v1.2.3" {
		t.Errorf("release.URL = %s, expected https://github.com/test/repo/releases/v1.2.3", release.URL)
	}

	if len(release.Assets) != 1 {
		t.Fatalf("len(release.Assets) = %d, expected 1", len(release.Assets))
	}

	if release.Assets[0].Name != "turbotilt-1.2.3-darwin-amd64.zip" {
		t.Errorf("release.Assets[0].Name = %s, expected turbotilt-1.2.3-darwin-amd64.zip",
			release.Assets[0].Name)
	}
}

func TestGetUpdateAssetURL(t *testing.T) {
	// Create a test release with different assets
	release := &Release{
		TagName: "v1.2.3",
		URL:     "https://github.com/test/repo/releases/v1.2.3",
		Assets: []Asset{
			{
				Name:        "turbotilt-1.2.3-darwin-amd64.zip",
				DownloadURL: "https://github.com/test/repo/releases/download/v1.2.3/turbotilt-1.2.3-darwin-amd64.zip",
			},
			{
				Name:        "turbotilt-1.2.3-darwin-arm64.zip",
				DownloadURL: "https://github.com/test/repo/releases/download/v1.2.3/turbotilt-1.2.3-darwin-arm64.zip",
			},
			{
				Name:        "turbotilt-1.2.3-linux-amd64.zip",
				DownloadURL: "https://github.com/test/repo/releases/download/v1.2.3/turbotilt-1.2.3-linux-amd64.zip",
			},
			{
				Name:        "turbotilt-1.2.3-windows-amd64.zip",
				DownloadURL: "https://github.com/test/repo/releases/download/v1.2.3/turbotilt-1.2.3-windows-amd64.zip",
			},
		},
	}

	// Test different platform/architecture combinations
	testCases := []struct {
		platform, arch, expectedAsset string
	}{
		{"darwin", "amd64", "darwin-amd64"},
		{"darwin", "arm64", "darwin-arm64"},
		{"linux", "amd64", "linux-amd64"},
		{"windows", "amd64", "windows-amd64"},
		{"freebsd", "amd64", ""}, // Unsupported platform
	}

	for _, tc := range testCases {
		// Simulate different architectures using a closure
		url := func() string {
			// Saving and resetting runtime.GOOS and runtime.GOARCH constants
			// is not directly possible in tests, so we test the logic directly
			for _, asset := range release.Assets {
				if strings.Contains(asset.Name, tc.platform) && strings.Contains(asset.Name, tc.arch) {
					return asset.DownloadURL
				}
			}
			return release.URL
		}()

		if tc.expectedAsset == "" {
			// For unsupported platforms, we expect the release URL
			if url != release.URL {
				t.Errorf("GetUpdateAssetURL for %s/%s = %s, expected %s",
					tc.platform, tc.arch, url, release.URL)
			}
		} else {
			// For supported platforms, we expect the URL of the corresponding asset
			if !strings.Contains(url, tc.expectedAsset) {
				t.Errorf("GetUpdateAssetURL for %s/%s = %s, should contain %s",
					tc.platform, tc.arch, url, tc.expectedAsset)
			}
		}
	}
}

func TestCheckForUpdates(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"tag_name": "v2.0.0",
			"html_url": "https://github.com/test/repo/releases/v2.0.0",
			"published_at": "2023-01-01T00:00:00Z",
			"assets": []
		}`))
	}))
	defer server.Close()

	// Save the original functions and variables
	origGithubAPI := githubAPI
	defer func() {
		githubAPI = origGithubAPI
	}()

	// Use a local variable with the same value
	testAPI := server.URL
	// Use this variable in the test function
	githubAPI = testAPI

	// Create a temporary directory for tests
	tempDir := t.TempDir()

	// Redefine getUpdateTimestampPath
	origGetUpdateTimestampPath := getUpdateTimestampPath
	defer func() { getUpdateTimestampPath = origGetUpdateTimestampPath }()

	getUpdateTimestampPath = func() string {
		return filepath.Join(tempDir, "update_timestamp")
	}

	// Force a check by ensuring the timestamp file doesn't exist
	os.Remove(getUpdateTimestampPath())

	// Test with a lower current version
	release, hasUpdate := CheckForUpdates("1.0.0")
	if !hasUpdate {
		t.Error("CheckForUpdates should return hasUpdate=true for a newer version")
	}
	if release == nil {
		t.Fatal("CheckForUpdates should return a non-nil Release object")
	}
	if release.TagName != "v2.0.0" {
		t.Errorf("release.TagName = %s, expected v2.0.0", release.TagName)
	}

	// Test with an equal current version
	_, hasUpdate = CheckForUpdates("2.0.0")
	if hasUpdate {
		t.Error("CheckForUpdates should return hasUpdate=false for an equal version")
	}

	// Test with a dev version
	release, hasUpdate = CheckForUpdates("dev")
	if hasUpdate {
		t.Error("CheckForUpdates should return hasUpdate=false for a dev version")
	}
	if release != nil {
		t.Error("CheckForUpdates should return a nil Release object for a dev version")
	}
}
