package render

import (
	"fmt"
	"os"
)

func GenerateDockerfile(framework string) error {
	f, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}
	defer f.Close()
	if framework == "spring" {
		_, err = f.WriteString("FROM eclipse-temurin:17\nCOPY . /app\nWORKDIR /app\nRUN ./mvnw package\nCMD java -jar target/*.jar\n")
	} else {
		_, err = f.WriteString("# TODO: Autres frameworks\n")
	}
	return err
}

func GenerateCompose() error {
	f, err := os.Create("docker-compose.yml")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("version: '3'\nservices:\n  app:\n    build: .\n    ports:\n      - '8080:8080'\n")
	return err
}
