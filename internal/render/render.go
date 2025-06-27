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

// DefaultDockerfileRenderer est l'implémentation par défaut de DockerfileRenderer
type DefaultDockerfileRenderer struct{}

// RenderSpringDockerfile écrit un Dockerfile pour Spring Boot
func (r *DefaultDockerfileRenderer) RenderSpringDockerfile(w io.Writer, opts Options) error {
	tmpl := `FROM eclipse-temurin:{{.JDKVersion}} AS build
WORKDIR /app
COPY . .
RUN ./mvnw package {{if .DevMode}}-DskipTests{{end}}

FROM eclipse-temurin:{{.JDKVersion}}
WORKDIR /app
COPY --from=build /app/target/*.jar app.jar
EXPOSE {{.Port}}
CMD ["java", "-jar", "app.jar"]
`
	t, err := template.New("spring").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, opts)
}

// RenderQuarkusDockerfile écrit un Dockerfile pour Quarkus
func (r *DefaultDockerfileRenderer) RenderQuarkusDockerfile(w io.Writer, opts Options) error {
	tmpl := `FROM registry.access.redhat.com/ubi8/openjdk-{{.JDKVersion}}:latest AS build
WORKDIR /app
COPY . .
RUN ./mvnw package {{if .DevMode}}-DskipTests{{end}} -Dquarkus.package.type=jar

FROM registry.access.redhat.com/ubi8/openjdk-{{.JDKVersion}}:latest
WORKDIR /app
COPY --from=build /app/target/quarkus-app/lib/ /deployments/lib/
COPY --from=build /app/target/quarkus-app/*.jar /deployments/
COPY --from=build /app/target/quarkus-app/app/ /deployments/app/
COPY --from=build /app/target/quarkus-app/quarkus/ /deployments/quarkus/
EXPOSE {{.Port}}
CMD ["java", "-jar", "/deployments/quarkus-run.jar"]
`
	t, err := template.New("quarkus").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, opts)
}

// RenderMicronautDockerfile écrit un Dockerfile pour Micronaut
func (r *DefaultDockerfileRenderer) RenderMicronautDockerfile(w io.Writer, opts Options) error {
	tmpl := `FROM eclipse-temurin:{{.JDKVersion}} AS build
WORKDIR /app
COPY . .
RUN ./gradlew build {{if .DevMode}}-x test{{end}}

FROM eclipse-temurin:{{.JDKVersion}}
WORKDIR /app
COPY --from=build /app/build/libs/*-all.jar app.jar
EXPOSE {{.Port}}
CMD ["java", "-jar", "app.jar"]
`
	t, err := template.New("micronaut").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, opts)
}

// RenderJavaDockerfile écrit un Dockerfile pour une application Java générique
func (r *DefaultDockerfileRenderer) RenderJavaDockerfile(w io.Writer, opts Options) error {
	tmpl := `FROM eclipse-temurin:{{.JDKVersion}}
WORKDIR /app
COPY . .
RUN javac Main.java
EXPOSE {{.Port}}
CMD ["java", "Main"]
`
	t, err := template.New("java").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, opts)
}

// RenderGenericDockerfile écrit un Dockerfile générique pour les autres types d'applications
func (r *DefaultDockerfileRenderer) RenderGenericDockerfile(w io.Writer, opts Options) error {
	tmpl := `FROM alpine:latest
WORKDIR /app
COPY . .
EXPOSE {{.Port}}
CMD ["sh", "start.sh"]
`
	t, err := template.New("generic").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, opts)
}

// defaultRenderer est l'instance par défaut du rendu
var defaultRenderer DockerfileRenderer = &DefaultDockerfileRenderer{}

// GenerateDockerfile génère un Dockerfile adapté au framework détecté
func GenerateDockerfile(opts Options) error {
	f, err := os.Create("Dockerfile")
	if err != nil {
		return fmt.Errorf("erreur lors de la création du Dockerfile: %w", err)
	}
	defer f.Close()

	switch opts.Framework {
	case "spring":
		return defaultRenderer.RenderSpringDockerfile(f, opts)
	case "quarkus":
		return defaultRenderer.RenderQuarkusDockerfile(f, opts)
	case "micronaut":
		return defaultRenderer.RenderMicronautDockerfile(f, opts)
	case "java":
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

	tmpl := `version: '3'
services:
  app:
    build: .
    ports:
      - '{{.Port}}:{{.Port}}'
    volumes:
      - './src:/app/src'
    environment:
{{if eq .Framework "spring"}}      - SPRING_PROFILES_ACTIVE={{if .DevMode}}dev{{else}}prod{{end}}
{{else if eq .Framework "quarkus"}}      - QUARKUS_PROFILE={{if .DevMode}}dev{{else}}prod{{end}}
{{else if eq .Framework "micronaut"}}      - MICRONAUT_ENVIRONMENTS={{if .DevMode}}dev{{else}}prod{{end}}
{{else}}      # Ajoutez vos variables d'environnement spécifiques ici
{{end}}
{{range .Services}}
  {{.Name}}:
    image: {{.Image}}
    ports:
      - '{{.Port}}'
    environment:
      - SPRING_PROFILES_ACTIVE={{if $.DevMode}}dev{{else}}prod{{end}}
{{end}}
`
	t, err := template.New("compose").Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(f, opts)
}
