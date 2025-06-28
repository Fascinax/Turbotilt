package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"turbotilt/internal/render"
)

// Constants for file names
const (
	LegacyConfigFileName = "turbotilt.yml"
	ManifestFileName     = "turbotilt.yaml"
)

// Config represents the Turbotilt configuration
type Config struct {
	Project     ProjectConfig     `yaml:"project,omitempty"`
	Framework   FrameworkConfig   `yaml:"framework,omitempty"`
	Docker      DockerConfig      `yaml:"docker,omitempty"`
	Development DevelopmentConfig `yaml:"development,omitempty"`
	Services    []ServiceConfig   `yaml:"services,omitempty"`
}

// ProjectConfig contains the general project information
type ProjectConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Version     string `yaml:"version,omitempty"`
}

// FrameworkConfig contains information about the framework used
type FrameworkConfig struct {
	Type       string `yaml:"type"` // spring, quarkus, micronaut, etc.
	Version    string `yaml:"version,omitempty"`
	JdkVersion string `yaml:"jdk_version,omitempty"`
}

// DockerConfig contains Docker parameters
type DockerConfig struct {
	Port      string            `yaml:"port"`
	BuildArgs map[string]string `yaml:"build_args,omitempty"`
}

// DevelopmentConfig contains development parameters
type DevelopmentConfig struct {
	EnableLiveReload bool     `yaml:"enable_live_reload"`
	SyncPath         string   `yaml:"sync_path,omitempty"`
	WatchPaths       []string `yaml:"watch_paths,omitempty"`
}

// ServiceConfig represents a dependent service in the old structure
type ServiceConfig struct {
	Name        string            `yaml:"name"`
	Type        string            `yaml:"type"` // mysql, postgres, etc.
	Version     string            `yaml:"version,omitempty"`
	Port        string            `yaml:"port,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
}

// Manifest represents the new declarative structure of the turbotilt.yaml file
type Manifest struct {
	Services []ManifestService `yaml:"services"`
}

// ManifestService represents a service in the declarative manifest
type ManifestService struct {
	Name       string            `yaml:"name"`
	Path       string            `yaml:"path"`
	Java       string            `yaml:"java,omitempty"`
	Build      string            `yaml:"build,omitempty"`   // maven, gradle
	Runtime    string            `yaml:"runtime,omitempty"` // spring, quarkus, micronaut
	Port       string            `yaml:"port,omitempty"`
	DevMode    bool              `yaml:"devMode,omitempty"`
	Type       string            `yaml:"type,omitempty"`       // For dependent services: mysql, postgres, etc.
	Version    string            `yaml:"version,omitempty"`    // For dependent services
	Env        map[string]string `yaml:"env,omitempty"`        // Environment variables
	Volumes    []string          `yaml:"volumes,omitempty"`    // Volume mounts
	WatchPaths []string          `yaml:"watchPaths,omitempty"` // Paths to watch for live reload
}

// DefaultConfig creates a default configuration
func DefaultConfig(frameworkType string) Config {
	config := Config{
		Project: ProjectConfig{
			Name:        filepath.Base(currentDir()),
			Description: "Turbotilt project",
			Version:     "1.0.0",
		},
		Framework: FrameworkConfig{
			Type:       frameworkType,
			JdkVersion: "17",
		},
		Docker: DockerConfig{
			Port:      "8080",
			BuildArgs: map[string]string{},
		},
		Development: DevelopmentConfig{
			EnableLiveReload: true,
			SyncPath:         "./src",
		},
	}

	return config
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}

// LoadManifest loads the declarative manifest from a file
func LoadManifest(path string) (Manifest, error) {
	var manifest Manifest

	data, err := os.ReadFile(path)
	if err != nil {
		return manifest, err
	}

	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		return manifest, err
	}

	// Validate the manifest schema
	if err = ValidateManifestSchema(manifest); err != nil {
		return manifest, fmt.Errorf("manifest validation failed: %w", err)
	}

	return manifest, nil
}

// ValidateManifestSchema validates that the manifest conforms to the expected schema
func ValidateManifestSchema(manifest Manifest) error {
	// Basic check: services must have a name and a path
	if len(manifest.Services) == 0 {
		return fmt.Errorf("the manifest must contain at least one service")
	}

	for i, service := range manifest.Services {
		if service.Name == "" {
			return fmt.Errorf("service #%d: name is required", i+1)
		}
		if service.Path == "" {
			return fmt.Errorf("service '%s': path is required", service.Name)
		}

		// Check if it's an application service (runtime) or a dependent service (type)
		if service.Runtime == "" && service.Type == "" {
			return fmt.Errorf("service '%s': either 'runtime' or 'type' must be specified", service.Name)
		}

		// Validate known runtime or service types
		if service.Runtime != "" && !isValidRuntime(service.Runtime) {
			return fmt.Errorf("service '%s': runtime '%s' not supported", service.Name, service.Runtime)
		}

		if service.Type != "" && !isValidServiceType(service.Type) {
			return fmt.Errorf("service '%s': type '%s' not supported", service.Name, service.Type)
		}
	}

	return nil
}

// isValidRuntime checks if the specified runtime is supported
func isValidRuntime(runtime string) bool {
	validRuntimes := map[string]bool{
		"spring":    true,
		"quarkus":   true,
		"micronaut": true,
		"java":      true,
	}
	return validRuntimes[strings.ToLower(runtime)]
}

// isValidServiceType checks if the specified service type is supported
func isValidServiceType(serviceType string) bool {
	validTypes := map[string]bool{
		"mysql":         true,
		"postgresql":    true,
		"postgres":      true,
		"redis":         true,
		"mongodb":       true,
		"kafka":         true,
		"rabbitmq":      true,
		"elasticsearch": true,
	}
	return validTypes[strings.ToLower(serviceType)]
}

// GenerateManifestFromConfig generates a declarative manifest from a config
func GenerateManifestFromConfig(config Config) Manifest {
	manifest := Manifest{
		Services: []ManifestService{},
	}

	// Create a main service from the configuration
	mainService := ManifestService{
		Name:    config.Project.Name,
		Path:    ".",
		Java:    config.Framework.JdkVersion,
		Runtime: config.Framework.Type,
		Port:    config.Docker.Port,
		DevMode: config.Development.EnableLiveReload,
	}

	manifest.Services = append(manifest.Services, mainService)

	// Add dependent services
	for _, svc := range config.Services {
		manifestService := ManifestService{
			Name:    svc.Name,
			Path:    "", // No path for dependent services
			Type:    svc.Type,
			Version: svc.Version,
			Port:    svc.Port,
			Env:     svc.Environment,
			Volumes: svc.Volumes,
		}
		manifest.Services = append(manifest.Services, manifestService)
	}

	return manifest
}

// ConvertManifestToRenderOptions converts a service from the manifest to rendering options
func ConvertManifestToRenderOptions(service ManifestService) (*render.Options, error) {
	if service.Runtime == "" {
		return nil, fmt.Errorf("service '%s' is not an application service (no runtime specified)", service.Name)
	}

	opts := &render.Options{
		ServiceName: service.Name,
		AppName:     service.Name,
		Framework:   service.Runtime,
		Port:        service.Port,
		Path:        service.Path,
		DevMode:     service.DevMode,
	}

	// Set default values if not specified
	if opts.Port == "" {
		opts.Port = "8080"
	}

	if service.Java != "" {
		opts.JDKVersion = service.Java
	} else {
		opts.JDKVersion = "17" // Default value
	}

	return opts, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// SaveManifest saves the manifest to a file
func SaveManifest(manifest Manifest, path string) error {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// FindConfiguration searches for the configuration or manifest file in the current directory
// returns the path of the found file and a boolean indicating whether it is a declarative manifest
func FindConfiguration() (string, bool, error) {
	// First, look for the new format (turbotilt.yaml)
	if _, err := os.Stat(ManifestFileName); err == nil {
		return ManifestFileName, true, nil
	}

	// Then, look for the old format (turbotilt.yml)
	if _, err := os.Stat(LegacyConfigFileName); err == nil {
		return LegacyConfigFileName, false, nil
	}

	return "", false, fmt.Errorf("no configuration file found (neither %s nor %s)", ManifestFileName, LegacyConfigFileName)
}

// currentDir returns the current directory
func currentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}
