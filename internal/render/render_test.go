package render

import (
	"os"
	"strings"
	"testing"
)

func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && strings.Contains(s, substr)
}

// TestGenerateDockerfile teste la génération de Dockerfile
func TestGenerateDockerfile(t *testing.T) {
	// Création d”un répertoire temporaire pour les tests
	tempDir, err := os.MkdirTemp("", "turbotilt-render-test-*")
	if err != nil {
		t.Fatalf("Impossible de créer le répertoire temporaire: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Sauvegarde du répertoire de travail actuel
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Impossible d''obtenir le répertoire de travail: %v", err)
	}

	// Changer vers le répertoire temporaire
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Impossible de changer vers le répertoire temporaire: %v", err)
	}
	defer os.Chdir(originalDir)

	// Test 1: Génération pour Spring Boot
	t.Run("Spring Dockerfile", func(t *testing.T) {
		// Options would normally be used to generate the Dockerfile
		// but for testing, we”re using a hardcoded example
		opts := Options{
			ServiceName: "api",
			Framework:   "spring",
			AppName:     "spring-app",
			Port:        "8080",
			JDKVersion:  "17",
			DevMode:     true,
			Path:        ".",
		}

		// Generate the Dockerfile using the options
		err := GenerateDockerfile(opts)
		if err != nil {
			t.Fatalf("Failed to generate Dockerfile: %v", err)
		}

		// Vérifier que le fichier existe
		if _, err := os.Stat("Dockerfile"); os.IsNotExist(err) {
			t.Errorf("Le Dockerfile n''a pas été créé")
		}

		// Lire le fichier généré
		fileContent, err := os.ReadFile("Dockerfile")
		if err != nil {
			t.Errorf("Impossible de lire le Dockerfile: %v", err)
		}

		// Vérifier que le contenu contient des éléments spécifiques à Spring
		if !containsString(string(fileContent), "eclipse-temurin") || !containsString(string(fileContent), "COPY --from=build") {
			t.Errorf("Le Dockerfile ne contient pas les éléments attendus pour Spring")
		}

		// Nettoyer
		os.Remove("Dockerfile")
	})
}
