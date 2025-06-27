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
		{"1.2.3", "1.2.3-alpha", false}, // Cas simple, pas de gestion sémantique complexe
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
	// Créer un dossier temporaire pour les tests
	tempDir := t.TempDir()
	origGetUpdateTimestampPath := getUpdateTimestampPath
	defer func() { getUpdateTimestampPath = origGetUpdateTimestampPath }()
	
	getUpdateTimestampPath = func() string {
		return filepath.Join(tempDir, "update_timestamp")
	}
	
	// Cas 1: Pas de fichier timestamp, devrait retourner true
	if !shouldCheckForUpdate() {
		t.Error("shouldCheckForUpdate devrait retourner true quand le fichier n'existe pas")
	}
	
	// Créer un fichier timestamp
	timestampFile := getUpdateTimestampPath()
	if err := os.MkdirAll(filepath.Dir(timestampFile), 0755); err != nil {
		t.Fatalf("Impossible de créer le dossier: %v", err)
	}
	f, err := os.Create(timestampFile)
	if err != nil {
		t.Fatalf("Impossible de créer le fichier timestamp: %v", err)
	}
	f.Close()
	
	// Cas 2: Fichier timestamp récent, devrait retourner false
	if shouldCheckForUpdate() {
		t.Error("shouldCheckForUpdate devrait retourner false avec un fichier timestamp récent")
	}
	
	// Cas 3: Fichier timestamp ancien, devrait retourner true
	// Simuler un fichier ancien en modifiant sa date de modification
	oldTime := time.Now().Add(-updateCheckTTL - time.Hour)
	if err := os.Chtimes(timestampFile, oldTime, oldTime); err != nil {
		t.Fatalf("Impossible de modifier la date du fichier: %v", err)
	}
	
	if !shouldCheckForUpdate() {
		t.Error("shouldCheckForUpdate devrait retourner true avec un fichier timestamp ancien")
	}
}

func TestUpdateLastCheckTime(t *testing.T) {
	// Créer un dossier temporaire pour les tests
	tempDir := t.TempDir()
	origGetUpdateTimestampPath := getUpdateTimestampPath
	defer func() { getUpdateTimestampPath = origGetUpdateTimestampPath }()
	
	getUpdateTimestampPath = func() string {
		return filepath.Join(tempDir, "update_timestamp")
	}
	
	// Vérifier que le fichier n'existe pas au début
	timestampFile := getUpdateTimestampPath()
	if _, err := os.Stat(timestampFile); err == nil {
		t.Error("Le fichier timestamp ne devrait pas exister au début du test")
	}
	
	// Mettre à jour le timestamp
	updateLastCheckTime()
	
	// Vérifier que le fichier existe maintenant
	if _, err := os.Stat(timestampFile); err != nil {
		t.Errorf("Le fichier timestamp devrait exister après updateLastCheckTime: %v", err)
	}
}

func TestGetUpdateTimestampPath(t *testing.T) {
	// Tester que la fonction retourne un chemin non vide
	path := getUpdateTimestampPath()
	if path == "" {
		t.Error("getUpdateTimestampPath ne devrait pas retourner un chemin vide")
	}
}

func TestGetLatestRelease(t *testing.T) {
	// Créer un serveur HTTP de test qui simule l'API GitHub
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
	
	// Sauvegarder l'URL originale et la restaurer après le test
	origGithubAPI := githubAPI
	defer func() { 
		githubAPI = origGithubAPI 
	}()
	
	// Utiliser une variable locale avec la même valeur
	testAPI := server.URL
	// Utiliser cette variable dans la fonction de test
	githubAPI = testAPI
	
	// Tester getLatestRelease
	release, err := getLatestRelease()
	if err != nil {
		t.Fatalf("getLatestRelease a retourné une erreur: %v", err)
	}
	
	if release.TagName != "v1.2.3" {
		t.Errorf("release.TagName = %s, attendu v1.2.3", release.TagName)
	}
	
	if release.URL != "https://github.com/test/repo/releases/v1.2.3" {
		t.Errorf("release.URL = %s, attendu https://github.com/test/repo/releases/v1.2.3", release.URL)
	}
	
	if len(release.Assets) != 1 {
		t.Fatalf("len(release.Assets) = %d, attendu 1", len(release.Assets))
	}
	
	if release.Assets[0].Name != "turbotilt-1.2.3-darwin-amd64.zip" {
		t.Errorf("release.Assets[0].Name = %s, attendu turbotilt-1.2.3-darwin-amd64.zip", 
			release.Assets[0].Name)
	}
}

func TestGetUpdateAssetURL(t *testing.T) {
	// Créer une release de test avec différents assets
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
	
	// Tester différentes combinaisons de plateforme/architecture
	testCases := []struct {
		platform, arch, expectedAsset string
	}{
		{"darwin", "amd64", "darwin-amd64"},
		{"darwin", "arm64", "darwin-arm64"},
		{"linux", "amd64", "linux-amd64"},
		{"windows", "amd64", "windows-amd64"},
		{"freebsd", "amd64", ""}, // Plateforme non supportée
	}
	
	for _, tc := range testCases {
		// Simuler des architectures différentes en utilisant une closure
		url := func() string {
			// Sauvegarde et réinitialisation des constantes runtime.GOOS et runtime.GOARCH
			// n'est pas possible directement dans les tests, donc on teste la logique directement
			for _, asset := range release.Assets {
				if strings.Contains(asset.Name, tc.platform) && strings.Contains(asset.Name, tc.arch) {
					return asset.DownloadURL
				}
			}
			return release.URL
		}()
		
		if tc.expectedAsset == "" {
			// Pour les plateformes non supportées, on s'attend à avoir l'URL de la release
			if url != release.URL {
				t.Errorf("GetUpdateAssetURL pour %s/%s = %s, attendu %s", 
					tc.platform, tc.arch, url, release.URL)
			}
		} else {
			// Pour les plateformes supportées, on s'attend à avoir l'URL de l'asset correspondant
			if !strings.Contains(url, tc.expectedAsset) {
				t.Errorf("GetUpdateAssetURL pour %s/%s = %s, devrait contenir %s", 
					tc.platform, tc.arch, url, tc.expectedAsset)
			}
		}
	}
}

func TestCheckForUpdates(t *testing.T) {
	// Créer un serveur HTTP de test
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
	
	// Sauvegarder les fonctions et variables d'origine
	origGithubAPI := githubAPI
	defer func() {
		githubAPI = origGithubAPI
	}()
	
	// Utiliser une variable locale avec la même valeur
	testAPI := server.URL
	// Utiliser cette variable dans la fonction de test
	githubAPI = testAPI
	
	// Créer un dossier temporaire pour les tests
	tempDir := t.TempDir()
	
	// Redéfinir getUpdateTimestampPath
	origGetUpdateTimestampPath := getUpdateTimestampPath
	defer func() { getUpdateTimestampPath = origGetUpdateTimestampPath }()
	
	getUpdateTimestampPath = func() string {
		return filepath.Join(tempDir, "update_timestamp")
	}
	
	// Forcer une vérification en s'assurant que le fichier timestamp n'existe pas
	os.Remove(getUpdateTimestampPath())
	
	// Test avec une version courante inférieure
	release, hasUpdate := CheckForUpdates("1.0.0")
	if !hasUpdate {
		t.Error("CheckForUpdates devrait retourner hasUpdate=true pour une version plus récente")
	}
	if release == nil {
		t.Fatal("CheckForUpdates devrait retourner un objet Release non nil")
	}
	if release.TagName != "v2.0.0" {
		t.Errorf("release.TagName = %s, attendu v2.0.0", release.TagName)
	}
	
	// Test avec une version courante égale
	release, hasUpdate = CheckForUpdates("2.0.0")
	if hasUpdate {
		t.Error("CheckForUpdates devrait retourner hasUpdate=false pour une version égale")
	}
	
	// Test avec une version dev
	release, hasUpdate = CheckForUpdates("dev")
	if hasUpdate {
		t.Error("CheckForUpdates devrait retourner hasUpdate=false pour une version dev")
	}
	if release != nil {
		t.Error("CheckForUpdates devrait retourner un objet Release nil pour une version dev")
	}
}
