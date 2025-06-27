package render

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"turbotilt/internal/scan"
)

// Options contient les paramètres de configuration pour la génération des fichiers
type Options struct {
	ServiceName string               // Nom du service (pour les projets multi-services)
	Framework   string               // Framework (spring, quarkus, etc.)
	AppName     string               // Nom de l'application
	Port        string               // Port exposé
	JDKVersion  string               // Version JDK
	DevMode     bool                 // Mode développement
	Path        string               // Chemin du service (pour les projets multi-services)
	Services    []scan.ServiceConfig // Services dépendants détectés
	EnvFile     string               // Chemin du fichier d'environnement
}

// ServiceList contient la liste des services pour la génération des fichiers multi-services
type ServiceList struct {
	Services []Options // Liste des services applicatifs
}

// DockerfileRenderer définit une interface pour la génération de Dockerfiles
type DockerfileRenderer interface {
	RenderSpringDockerfile(w io.Writer, opts Options) error
	RenderQuarkusDockerfile(w io.Writer, opts Options) error
	RenderMicronautDockerfile(w io.Writer, opts Options) error
	RenderJavaDockerfile(w io.Writer, opts Options) error
	RenderGenericDockerfile(w io.Writer, opts Options) error
}

// GenerateDockerfile génère un Dockerfile adapté au framework détecté
func GenerateDockerfile(opts Options) error {
	f, err := os.Create("Dockerfile")
	if err != nil {
		return fmt.Errorf("erreur lors de la création du Dockerfile: %w", err)
	}
	defer f.Close()

	switch opts.Framework {
	case FrameworkSpring:
		return defaultRenderer.RenderSpringDockerfile(f, opts)
	case FrameworkQuarkus:
		return defaultRenderer.RenderQuarkusDockerfile(f, opts)
	case FrameworkMicronaut:
		return defaultRenderer.RenderMicronautDockerfile(f, opts)
	case FrameworkJava:
		return defaultRenderer.RenderJavaDockerfile(f, opts)
	default:
		return defaultRenderer.RenderGenericDockerfile(f, opts)
	}
}

// GenerateCompose génère un docker-compose.yml
func GenerateCompose(opts Options) error {
	f, err := os.Create("docker-compose.yml")
	if err != nil {
		return fmt.Errorf("erreur lors de la création du docker-compose.yml: %w", err)
	}
	defer f.Close()

	// Vérifier si un fichier d'environnement existe
	envFile := getEnvFilePath(".")

	// Construire le template avec ou sans fichier d'environnement
	var tmplStr string
	if envFile != "" {
		// Ajouter le chemin du fichier d'environnement aux options
		opts.EnvFile = envFile
		tmplStr = ComposeTemplateWithEnvFile
	} else {
		tmplStr = ComposeTemplate
	}

	t, err := template.New("compose").Parse(tmplStr)
	if err != nil {
		return err
	}

	return t.Execute(f, opts)
}
