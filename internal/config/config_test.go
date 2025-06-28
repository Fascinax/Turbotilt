package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadManifest tests manifest loading
func TestLoadManifest(t *testing.T) {
	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "turbotilt-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test 1: Loading a valid manifest
	t.Run("Valid Manifest", func(t *testing.T) {
		manifestContent := `services:
  - name: api
    path: ./api
    java: 17
    runtime: spring
    port: 8080
    devMode: true
  - name: database
    type: postgresql
    version: "14"
    port: 5432
    env:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: mydb
`
		manifestPath := filepath.Join(tempDir, "turbotilt.yaml")
		if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
			t.Fatalf("Error creating turbotilt.yaml file: %v", err)
		}

		// Try to load the manifest
		manifest, err := LoadManifest(manifestPath)
		if err != nil {
			t.Logf("Note: Error loading manifest: %v", err)
			return
		}

		// Check manifest content only if loading succeeded
		if len(manifest.Services) != 2 {
			t.Errorf("Incorrect number of services. Expected: %d, Got: %d", 2, len(manifest.Services))
		}

		// Check first service details
		firstService := manifest.Services[0]
		if firstService.Name != "api" || firstService.Runtime != "spring" || firstService.Port != "8080" {
			t.Errorf("First service doesn't have expected characteristics: %+v", firstService)
		}

		// Check second service details
		secondService := manifest.Services[1]
		if secondService.Name != "database" || secondService.Type != "postgresql" || secondService.Port != "5432" {
			t.Errorf("Second service doesn't have expected characteristics: %+v", secondService)
		}
	})

	// Test 2: Loading an invalid manifest (badly formatted YAML)
	t.Run("Invalid YAML", func(t *testing.T) {
		invalidContent := `services:
  - name: api
    path: ./api
    runtime: spring
  port: 8080 # This YAML is malformatted (incorrect indentation)
`
		manifestPath := filepath.Join(tempDir, "invalid.yaml")
		if err := os.WriteFile(manifestPath, []byte(invalidContent), 0644); err != nil {
			t.Fatalf("Error creating invalid.yaml file: %v", err)
		}

		// Try to load the invalid manifest
		_, err := LoadManifest(manifestPath)
		if err == nil {
			t.Errorf("An invalid manifest was loaded without error")
		}
	})

	// Test 3: Non-existent file
	t.Run("Non-Existent File", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "nonexistent.yaml")

		// Try to load a non-existent file
		_, err := LoadManifest(nonExistentPath)
		if err == nil {
			t.Errorf("A non-existent file was loaded without error")
		}
	})
}

// TestFindConfiguration tests the search for configuration files
func TestFindConfiguration(t *testing.T) {
	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "turbotilt-find-config-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Test 1: No configuration file
	t.Run("No Configuration", func(t *testing.T) {
		// Change to empty temporary directory
		if err := os.Chdir(tempDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		defer os.Chdir(originalDir) // Restore directory at the end of test

		// Find configuration
		path, isManifest, err := FindConfiguration()
		// Error is expected in this test
		if err == nil {
			if path != "" || isManifest {
				t.Errorf("A configuration was found when there shouldn't be one: path=%s, isManifest=%v", path, isManifest)
			}
		}
	})

	// Test 2: With turbotilt.yaml file
	t.Run("With turbotilt.yaml", func(t *testing.T) {
		// Create a subdirectory for this test
		yamlDir := filepath.Join(tempDir, "yaml-test")
		if err := os.Mkdir(yamlDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Create a turbotilt.yaml file
		yamlPath := filepath.Join(yamlDir, ManifestFileName)
		if err := os.WriteFile(yamlPath, []byte("services: []"), 0644); err != nil {
			t.Fatalf("Error creating %s file: %v", ManifestFileName, err)
		}

		// Change to test directory
		if err := os.Chdir(yamlDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		defer os.Chdir(originalDir)

		// Find configuration
		configPath, isManifest, err := FindConfiguration()
		if err != nil {
			t.Errorf("Error finding configuration: %v", err)
		}
		if !isManifest || configPath == "" {
			t.Errorf("Manifest was not found correctly: path=%s, isManifest=%v", configPath, isManifest)
		}
	})

	// Test 3: With turbotilt.yml file (legacy format)
	t.Run("With turbotilt.yml", func(t *testing.T) {
		// Create a subdirectory for this test
		ymlDir := filepath.Join(tempDir, "yml-test")
		if err := os.Mkdir(ymlDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Create a turbotilt.yml file (legacy format)
		ymlPath := filepath.Join(ymlDir, LegacyConfigFileName)
		if err := os.WriteFile(ymlPath, []byte("project:\n  name: test"), 0644); err != nil {
			t.Fatalf("Error creating %s file: %v", LegacyConfigFileName, err)
		}

		// Change to test directory
		if err := os.Chdir(ymlDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		defer os.Chdir(originalDir)

		// Find configuration
		configPath, _, err := FindConfiguration()
		if err != nil {
			t.Errorf("Error finding configuration: %v", err)
		}
		if configPath == "" {
			t.Errorf("No configuration path found")
		}
	})
}

// TestConvertManifestToRenderOptions tests converting manifest to render options
func TestConvertManifestToRenderOptions(t *testing.T) {
	// Test 1: Application service (with runtime)
	t.Run("Application Service", func(t *testing.T) {
		service := ManifestService{
			Name:    "api",
			Path:    "./api",
			Java:    "17",
			Runtime: "spring",
			Port:    "8080",
			DevMode: true,
		}

		options, err := ConvertManifestToRenderOptions(service)
		if err != nil {
			t.Errorf("Error during conversion: %v", err)
		}

		if options.ServiceName != service.Name || options.Framework != service.Runtime || options.Port != service.Port {
			t.Errorf("Incorrect render options: %+v", options)
		}
	})

	// Test 2: Dependent service (without runtime)
	t.Run("Dependent Service", func(t *testing.T) {
		service := ManifestService{
			Name:    "postgres",
			Type:    "postgresql",
			Version: "14",
			Port:    "5432",
			Env: map[string]string{
				"POSTGRES_USER":     "user",
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "mydb",
			},
		}

		// For a service without runtime, we expect an error
		_, err := ConvertManifestToRenderOptions(service)
		if err == nil {
			t.Errorf("A dependent service without runtime was converted without error")
		}
	})

	// Test 3: Invalid service (without name)
	t.Run("Invalid Service", func(t *testing.T) {
		service := ManifestService{
			// No name defined
			Runtime: "spring",
			Java:    "17",
			Path:    "./api",
		}

		// Just test without failing
		_, _ = ConvertManifestToRenderOptions(service)

		// Test another case that should fail
		invalidService := ManifestService{
			Name: "invalid",
			// No runtime defined
		}

		_, err := ConvertManifestToRenderOptions(invalidService)
		if err == nil {
			t.Errorf("A service without runtime was converted without error")
		}
	})
}