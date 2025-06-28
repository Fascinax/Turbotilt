package render

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"turbotilt/internal/scan"
)

// Options contains configuration parameters for file generation
type Options struct {
	ServiceName string               // Service name (for multi-service projects)
	Framework   string               // Framework (spring, quarkus, etc.)
	AppName     string               // Application name
	Port        string               // Exposed port
	JDKVersion  string               // JDK version
	DevMode     bool                 // Development mode
	Path        string               // Service path (for multi-service projects)
	Services    []scan.ServiceConfig // Detected dependent services
	EnvFile     string               // Path to environment file
}

// ServiceList contains the list of services for multi-service file generation
type ServiceList struct {
	Services []Options // List of application services
}

// DockerfileRenderer defines an interface for Dockerfile generation
type DockerfileRenderer interface {
	RenderSpringDockerfile(w io.Writer, opts Options) error
	RenderQuarkusDockerfile(w io.Writer, opts Options) error
	RenderMicronautDockerfile(w io.Writer, opts Options) error
	RenderJavaDockerfile(w io.Writer, opts Options) error
	RenderGenericDockerfile(w io.Writer, opts Options) error
}

// GenerateDockerfile generates a Dockerfile adapted to the detected framework
func GenerateDockerfile(opts Options) error {
	f, err := os.Create("Dockerfile")
	if err != nil {
		return fmt.Errorf("error creating Dockerfile: %w", err)
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

// GenerateCompose generates a docker-compose.yml file
func GenerateCompose(opts Options) error {
	f, err := os.Create("docker-compose.yml")
	if err != nil {
		return fmt.Errorf("error creating docker-compose.yml: %w", err)
	}
	defer f.Close()

	// Check if an environment file exists
	envFile := getEnvFilePath(".")

	// Build the template with or without an environment file
	var tmplStr string
	if envFile != "" {
		// Add the environment file path to the options
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
