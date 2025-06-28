package render

import (
	"os"
	"strings"
	"testing"
)

func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && strings.Contains(s, substr)
}

// TestGenerateDockerfile tests the generation of Dockerfile
func TestGenerateDockerfile(t *testing.T) {
	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "turbotilt-render-test-*")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save the current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get current working directory: %v", err)
	}

	// Change to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Unable to change to temporary directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Test 1: Generation for Spring Boot
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

		// Check that the file exists
		if _, err := os.Stat("Dockerfile"); os.IsNotExist(err) {
			t.Errorf("Dockerfile was not created")
		}

		// Read the generated file
		fileContent, err := os.ReadFile("Dockerfile")
		if err != nil {
			t.Errorf("Unable to read Dockerfile: %v", err)
		}

		// Check that the content contains specific elements for Spring
		if !containsString(string(fileContent), "eclipse-temurin") || !containsString(string(fileContent), "COPY --from=build") {
			t.Errorf("Dockerfile does not contain the expected elements for Spring")
		}

		// Clean up
		os.Remove("Dockerfile")
	})
}
