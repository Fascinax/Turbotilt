package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"turbotilt/internal/render"
)

// Constantes pour les noms de fichier
const (
	LegacyConfigFileName = "turbotilt.yml"
	ManifestFileName     = "turbotilt.yaml"
)

// Config représente la configuration de Turbotilt
type Config struct {
	Project     ProjectConfig     `yaml:"project,omitempty"`
	Framework   FrameworkConfig   `yaml:"framework,omitempty"`
	Docker      DockerConfig      `yaml:"docker,omitempty"`
	Development DevelopmentConfig `yaml:"development,omitempty"`
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
	Type       string `yaml:"type"` // spring, quarkus, micronaut, etc.
	Version    string `yaml:"version,omitempty"`
	JdkVersion string `yaml:"jdk_version,omitempty"`
}

// DockerConfig contient les paramètres Docker
type DockerConfig struct {
	Port      string            `yaml:"port"`
	BuildArgs map[string]string `yaml:"build_args,omitempty"`
}

// DevelopmentConfig contient les paramètres de développement
type DevelopmentConfig struct {
	EnableLiveReload bool     `yaml:"enable_live_reload"`
	SyncPath         string   `yaml:"sync_path,omitempty"`
	WatchPaths       []string `yaml:"watch_paths,omitempty"`
}

// ServiceConfig représente un service dépendant dans l'ancienne structure
type ServiceConfig struct {
	Name        string            `yaml:"name"`
	Type        string            `yaml:"type"` // mysql, postgres, etc.
	Version     string            `yaml:"version,omitempty"`
	Port        string            `yaml:"port,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
}

// Manifest représente la nouvelle structure déclarative du fichier turbotilt.yaml
type Manifest struct {
	Services []ManifestService `yaml:"services"`
}

// ManifestService représente un service dans le manifeste déclaratif
type ManifestService struct {
	Name       string            `yaml:"name"`
	Path       string            `yaml:"path"`
	Java       string            `yaml:"java,omitempty"`
	Build      string            `yaml:"build,omitempty"`   // maven, gradle
	Runtime    string            `yaml:"runtime,omitempty"` // spring, quarkus, micronaut
	Port       string            `yaml:"port,omitempty"`
	DevMode    bool              `yaml:"devMode,omitempty"`
	Type       string            `yaml:"type,omitempty"`       // Pour les services dépendants: mysql, postgres, etc.
	Version    string            `yaml:"version,omitempty"`    // Pour les services dépendants
	Env        map[string]string `yaml:"env,omitempty"`        // Variables d'environnement
	Volumes    []string          `yaml:"volumes,omitempty"`    // Montages de volumes
	WatchPaths []string          `yaml:"watchPaths,omitempty"` // Chemins à surveiller pour le live reload
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

// LoadManifest charge le manifeste déclaratif à partir d'un fichier
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

	// Valider le schéma du manifeste
	if err = ValidateManifestSchema(manifest); err != nil {
		return manifest, fmt.Errorf("validation du manifeste échouée: %w", err)
	}

	return manifest, nil
}

// ValidateManifestSchema valide que le manifeste est conforme au schéma attendu
func ValidateManifestSchema(manifest Manifest) error {
	// Vérification basique: les services doivent avoir un nom et un chemin
	if len(manifest.Services) == 0 {
		return errors.New("le manifeste doit contenir au moins un service")
	}

	for i, service := range manifest.Services {
		if service.Name == "" {
			return fmt.Errorf("service #%d: le nom est obligatoire", i+1)
		}
		if service.Path == "" {
			return fmt.Errorf("service '%s': le chemin est obligatoire", service.Name)
		}

		// Vérifier si c'est un service applicatif (runtime) ou dépendant (type)
		if service.Runtime == "" && service.Type == "" {
			return fmt.Errorf("service '%s': soit 'runtime' soit 'type' doit être spécifié", service.Name)
		}

		// Valider les types de runtime ou service connus
		if service.Runtime != "" && !isValidRuntime(service.Runtime) {
			return fmt.Errorf("service '%s': runtime '%s' non supporté", service.Name, service.Runtime)
		}

		if service.Type != "" && !isValidServiceType(service.Type) {
			return fmt.Errorf("service '%s': type '%s' non supporté", service.Name, service.Type)
		}
	}

	return nil
}

// isValidRuntime vérifie si le runtime spécifié est supporté
func isValidRuntime(runtime string) bool {
	validRuntimes := map[string]bool{
		"spring":    true,
		"quarkus":   true,
		"micronaut": true,
		"java":      true,
	}
	return validRuntimes[strings.ToLower(runtime)]
}

// isValidServiceType vérifie si le type de service spécifié est supporté
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

// GenerateManifestFromConfig génère un manifeste déclaratif à partir d'une config
func GenerateManifestFromConfig(config Config) Manifest {
	manifest := Manifest{
		Services: []ManifestService{},
	}

	// Créer un service principal à partir de la configuration
	mainService := ManifestService{
		Name:    config.Project.Name,
		Path:    ".",
		Java:    config.Framework.JdkVersion,
		Runtime: config.Framework.Type,
		Port:    config.Docker.Port,
		DevMode: config.Development.EnableLiveReload,
	}

	manifest.Services = append(manifest.Services, mainService)

	// Ajouter les services dépendants
	for _, svc := range config.Services {
		manifestService := ManifestService{
			Name:    svc.Name,
			Path:    "", // Pas de path pour les services dépendants
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

// ConvertManifestToRenderOptions convertit un service du manifeste en options de rendu
func ConvertManifestToRenderOptions(service ManifestService) (*render.Options, error) {
	if service.Runtime == "" {
		return nil, fmt.Errorf("service '%s' n'est pas un service applicatif (pas de runtime spécifié)", service.Name)
	}

	opts := &render.Options{
		ServiceName: service.Name,
		AppName:     service.Name,
		Framework:   service.Runtime,
		Port:        service.Port,
		Path:        service.Path,
		DevMode:     service.DevMode,
	}

	// Définir des valeurs par défaut si non spécifiées
	if opts.Port == "" {
		opts.Port = "8080"
	}

	if service.Java != "" {
		opts.JDKVersion = service.Java
	} else {
		opts.JDKVersion = "17" // Valeur par défaut
	}

	return opts, nil
}

// SaveConfig enregistre la configuration dans un fichier
func SaveConfig(config Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// SaveManifest enregistre le manifeste dans un fichier
func SaveManifest(manifest Manifest, path string) error {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// FindConfiguration recherche le fichier de configuration ou manifeste dans le répertoire courant
// retourne le chemin du fichier trouvé et un booléen indiquant s'il s'agit d'un manifeste déclaratif
func FindConfiguration() (string, bool, error) {
	// Chercher d'abord le nouveau format (turbotilt.yaml)
	if _, err := os.Stat(ManifestFileName); err == nil {
		return ManifestFileName, true, nil
	}

	// Chercher ensuite l'ancien format (turbotilt.yml)
	if _, err := os.Stat(LegacyConfigFileName); err == nil {
		return LegacyConfigFileName, false, nil
	}

	return "", false, fmt.Errorf("aucun fichier de configuration trouvé (ni %s ni %s)", ManifestFileName, LegacyConfigFileName)
}

// currentDir renvoie le répertoire courant
func currentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}
