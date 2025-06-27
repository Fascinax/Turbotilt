package render

import (
	"fmt"
	"os"
	"text/template"
	"turbotilt/internal/scan"
)

// Options contient les paramètres de configuration pour la génération des fichiers
type Options struct {
	Framework   string
	Port        string
	JDKVersion  string
	DevMode     bool
	Services    []scan.ServiceConfig // Services dépendants détectés
}

// GenerateDockerfile génère un Dockerfile adapté au framework détecté
func GenerateDockerfile(opts Options) error {
	f, err := os.Create("Dockerfile")
	if err != nil {
		return fmt.Errorf("erreur lors de la création du Dockerfile: %w", err)
	}
	defer f.Close()

	switch opts.Framework {
	case "spring":
		return writeSpringDockerfile(f, opts)
	case "quarkus":
		return writeQuarkusDockerfile(f, opts)
	case "micronaut":
		return writeMicronautDockerfile(f, opts)
	case "java":
		return writeJavaDockerfile(f, opts)
	default:
		return writeGenericDockerfile(f, opts)
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

// writeSpringDockerfile écrit un Dockerfile pour Spring Boot
func writeSpringDockerfile(f *os.File, opts Options) error {
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
	t, err := template.New("dockerfile").Parse(tmpl)
	if err != nil {
		return err
	}
	
	return t.Execute(f, opts)
}

// writeQuarkusDockerfile écrit un Dockerfile pour Quarkus
func writeQuarkusDockerfile(f *os.File, opts Options) error {
	tmpl := `FROM quay.io/quarkus/ubi-quarkus-native-image:22.3-java{{.JDKVersion}} AS build
WORKDIR /app
COPY . .
RUN ./mvnw package -Pnative {{if .DevMode}}-DskipTests{{end}}

FROM quay.io/quarkus/quarkus-micro-image:2.0
WORKDIR /app
COPY --from=build /app/target/*-runner /app/application
EXPOSE {{.Port}}
CMD ["./application", "-Dquarkus.http.host=0.0.0.0"]
`
	t, err := template.New("dockerfile").Parse(tmpl)
	if err != nil {
		return err
	}
	
	return t.Execute(f, opts)
}

// writeMicronautDockerfile écrit un Dockerfile pour Micronaut
func writeMicronautDockerfile(f *os.File, opts Options) error {
	tmpl := `FROM eclipse-temurin:{{.JDKVersion}} AS build
WORKDIR /app
COPY . .
RUN ./gradlew build {{if .DevMode}}-x test{{end}}

FROM eclipse-temurin:{{.JDKVersion}}
WORKDIR /app
COPY --from=build /app/build/libs/*.jar app.jar
EXPOSE {{.Port}}
CMD ["java", "-jar", "app.jar"]
`
	t, err := template.New("dockerfile").Parse(tmpl)
	if err != nil {
		return err
	}
	
	return t.Execute(f, opts)
}

// writeJavaDockerfile écrit un Dockerfile générique pour Java
func writeJavaDockerfile(f *os.File, opts Options) error {
	tmpl := `FROM eclipse-temurin:{{.JDKVersion}}
WORKDIR /app
COPY . .
RUN ./mvnw package {{if .DevMode}}-DskipTests{{end}}
EXPOSE {{.Port}}
CMD ["java", "-jar", "target/*.jar"]
`
	t, err := template.New("dockerfile").Parse(tmpl)
	if err != nil {
		return err
	}
	
	return t.Execute(f, opts)
}

// writeGenericDockerfile écrit un Dockerfile générique
func writeGenericDockerfile(f *os.File, opts Options) error {
	tmpl := `# Dockerfile générique
FROM alpine:latest
WORKDIR /app
COPY . .
EXPOSE {{.Port}}
CMD ["echo", "Veuillez configurer ce Dockerfile pour votre application"]
`
	t, err := template.New("dockerfile").Parse(tmpl)
	if err != nil {
		return err
	}
	
	return t.Execute(f, opts)
}
