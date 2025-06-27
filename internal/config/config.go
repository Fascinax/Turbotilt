package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Config représente la configuration de Turbotilt
type Config struct {
	Project     ProjectConfig     `yaml:"project"`
	Framework   FrameworkConfig   `yaml:"framework"`
	Docker      DockerConfig      `yaml:"docker"`
	Development DevelopmentConfig `yaml:"development"`
	Services    []ServiceConfig   `yaml:"services,omitempty"`
}

// ProjectConfig contient les informations générales du projet
type ProjectConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Version     string `yaml:"version,omitempty"`
}

// FrameworkConfig contient les informations sur le framework utilisé
type FrameworkConfig struct {
	Type    string `yaml:"type"`    // spring, quarkus, micronaut, etc.
	Version string `yaml:"version,omitempty"`
	JdkVersion string `yaml:"jdk_version,omitempty"`
}

// DockerConfig contient les paramètres Docker
type DockerConfig struct {
	Port      string `yaml:"port"`
	BuildArgs map[string]string `yaml:"build_args,omitempty"`
}

// DevelopmentConfig contient les paramètres de développement
type DevelopmentConfig struct {
	EnableLiveReload bool   `yaml:"enable_live_reload"`
	SyncPath         string `yaml:"sync_path,omitempty"`
	WatchPaths       []string `yaml:"watch_paths,omitempty"`
}

// ServiceConfig représente un service dépendant
type ServiceConfig struct {
	Name        string            `yaml:"name"`
	Type        string            `yaml:"type"` // mysql, postgres, etc.
	Version     string            `yaml:"version,omitempty"`
	Port        string            `yaml:"port,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
}

// DefaultConfig crée une configuration par défaut
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

// LoadConfig charge la configuration à partir d'un fichier
func LoadConfig(path string) (Config, error) {
	var config Config
	
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	
	err = yaml.Unmarshal(data, &config)
	return config, err
}

// SaveConfig enregistre la configuration dans un fichier
func SaveConfig(config Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// currentDir renvoie le répertoire courant
func currentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}
