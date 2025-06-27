package render

import (
	"fmt"
	"os"
)

// GenerateDockerfile génère un Dockerfile adapté au framework détecté
func GenerateDockerfile(framework string) error {
	f, err := os.Create("Dockerfile")
	if err != nil {
		return fmt.Errorf("erreur lors de la création du Dockerfile: %w", err)
	}
	defer f.Close()

	switch framework {
	case "spring":
		return writeSpringDockerfile(f)
	case "quarkus":
		return writeQuarkusDockerfile(f)
	case "java":
		return writeJavaDockerfile(f)
	default:
		return writeGenericDockerfile(f)
	}
}

// GenerateCompose génère un docker-compose.yml
func GenerateCompose() error {
	f, err := os.Create("docker-compose.yml")
	if err != nil {
		return fmt.Errorf("erreur lors de la création du docker-compose.yml: %w", err)
	}
	defer f.Close()

	composeContent := `version: '3'
services:
  app:
    build: .
    ports:
      - '8080:8080'
    volumes:
      - './src:/app/src'
    environment:
      - SPRING_PROFILES_ACTIVE=dev
`
	_, err = f.WriteString(composeContent)
	return err
}

// writeSpringDockerfile écrit un Dockerfile pour Spring Boot
func writeSpringDockerfile(f *os.File) error {
	content := `FROM eclipse-temurin:17 AS build
WORKDIR /app
COPY . .
RUN ./mvnw package -DskipTests

FROM eclipse-temurin:17
WORKDIR /app
COPY --from=build /app/target/*.jar app.jar
EXPOSE 8080
CMD ["java", "-jar", "app.jar"]
`
	_, err := f.WriteString(content)
	return err
}

// writeQuarkusDockerfile écrit un Dockerfile pour Quarkus
func writeQuarkusDockerfile(f *os.File) error {
	content := `FROM quay.io/quarkus/ubi-quarkus-native-image:22.3-java17 AS build
WORKDIR /app
COPY . .
RUN ./mvnw package -Pnative -DskipTests

FROM quay.io/quarkus/quarkus-micro-image:2.0
WORKDIR /app
COPY --from=build /app/target/*-runner /app/application
EXPOSE 8080
CMD ["./application", "-Dquarkus.http.host=0.0.0.0"]
`
	_, err := f.WriteString(content)
	return err
}

// writeJavaDockerfile écrit un Dockerfile générique pour Java
func writeJavaDockerfile(f *os.File) error {
	content := `FROM eclipse-temurin:17
WORKDIR /app
COPY . .
RUN ./mvnw package -DskipTests
CMD ["java", "-jar", "target/*.jar"]
`
	_, err := f.WriteString(content)
	return err
}

// writeGenericDockerfile écrit un Dockerfile générique
func writeGenericDockerfile(f *os.File) error {
	content := `# Dockerfile générique
FROM alpine:latest
WORKDIR /app
COPY . .
CMD ["echo", "Veuillez configurer ce Dockerfile pour votre application"]
`
	_, err := f.WriteString(content)
	return err
}
