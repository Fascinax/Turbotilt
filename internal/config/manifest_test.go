package config

import (
	"os"
	"path/filepath"
	"testing"

	"turbotilt/internal/render"
)

func TestManifestServiceToRenderOptions(t *testing.T) {
	tests := []struct {
		name    string
		service ManifestService
		want    render.Options
	}{
		{
			name: "Spring service",
			service: ManifestService{
				Name:       "test-app",
				Path:       "./app",
				Java:       "17",
				Build:      "maven",
				Runtime:    "spring",
				Port:       "8080",
				DevMode:    true,
				Env:        map[string]string{"SPRING_PROFILES_ACTIVE": "dev"},
				WatchPaths: []string{"src/main/java", "src/main/resources"},
			},
			want: render.Options{
				AppName:    "test-app",
				Framework:  "spring",
				JDKVersion: "17",
				Port:       "8080",
				Path:       "./app",
				DevMode:    true,
			},
		},
		{
			name: "Quarkus service",
			service: ManifestService{
				Name:    "quarkus-app",
				Path:    "./quarkus-app",
				Java:    "21",
				Build:   "gradle",
				Runtime: "quarkus",
				Port:    "8081",
				DevMode: true,
			},
			want: render.Options{
				AppName:    "quarkus-app",
				Framework:  "quarkus",
				JDKVersion: "21",
				Port:       "8081",
				Path:       "./quarkus-app",
				DevMode:    true,
			},
		},
		{
			name: "Micronaut service",
			service: ManifestService{
				Name:    "micronaut-app",
				Path:    "./micronaut-app",
				Java:    "17",
				Build:   "gradle",
				Runtime: "micronaut",
				Port:    "8082",
				DevMode: true,
			},
			want: render.Options{
				AppName:    "micronaut-app",
				Framework:  "micronaut",
				JDKVersion: "17",
				Port:       "8082",
				Path:       "./micronaut-app",
				DevMode:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertManifestToRenderOptions(tt.service)
			if err != nil {
				t.Errorf("ConvertManifestToRenderOptions() error = %v", err)
				return
			}

			if got.AppName != tt.want.AppName {
				t.Errorf("AppName = %v, want %v", got.AppName, tt.want.AppName)
			}
			if got.Framework != tt.want.Framework {
				t.Errorf("Framework = %v, want %v", got.Framework, tt.want.Framework)
			}
			if got.JDKVersion != tt.want.JDKVersion {
				t.Errorf("JDKVersion = %v, want %v", got.JDKVersion, tt.want.JDKVersion)
			}
			if got.Port != tt.want.Port {
				t.Errorf("Port = %v, want %v", got.Port, tt.want.Port)
			}
			if got.Path != tt.want.Path {
				t.Errorf("Path = %v, want %v", got.Path, tt.want.Path)
			}
			if got.DevMode != tt.want.DevMode {
				t.Errorf("DevMode = %v, want %v", got.DevMode, tt.want.DevMode)
			}
		})
	}
}

func TestValidateManifest(t *testing.T) {
	tests := []struct {
		name       string
		manifest   string
		wantErrors bool
	}{
		{
			name: "Valid manifest",
			manifest: `services:
  - name: test-app
    path: ./app
    java: "17"
    build: maven
    runtime: spring
    port: "8080"
    devMode: true`,
			wantErrors: false,
		},
		{
			name: "Missing required fields",
			manifest: `services:
  - name: test-app
    runtime: spring`,
			wantErrors: true,
		},
		{
			name: "Invalid framework",
			manifest: `services:
  - name: test-app
    path: ./app
    java: "17"
    build: maven
    runtime: invalid
    port: "8080"`,
			wantErrors: true,
		},
		{
			name: "Valid with dependent services",
			manifest: `services:
  - name: test-app
    path: ./app
    java: "17"
    build: maven
    runtime: spring
    port: "8080"
  - name: mysql
    type: mysql
    version: "8.0"
    path: ./mysql`,
			wantErrors: false,
		},
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "turbotilt-test-")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary manifest file
			manifestPath := filepath.Join(tempDir, "turbotilt.yaml")
			err := os.WriteFile(manifestPath, []byte(tt.manifest), 0644)
			if err != nil {
				t.Fatalf("Error writing file: %v", err)
			}

			// Load and validate the manifest
			manifest, err := LoadManifest(manifestPath)
			if err != nil && !tt.wantErrors {
				t.Errorf("LoadManifest() error = %v", err)
				return
			}

			if err == nil && tt.wantErrors {
				t.Errorf("LoadManifest() a réussi alors qu'une erreur était attendue")
			}

			if err == nil {
				if len(manifest.Services) == 0 {
					t.Errorf("LoadManifest() n'a chargé aucun service")
				}
			}
		})
	}
}
