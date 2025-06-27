package scan

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Project struct {
	XMLName     xml.Name `xml:"project"`
	Dependencies struct {
		Dependencies []struct {
			GroupId    string `xml:"groupId"`
			ArtifactId string `xml:"artifactId"`
		} `xml:"dependency"`
	} `xml:"dependencies"`
}

// DetectFramework détecte le framework utilisé dans le projet
func DetectFramework() (string, error) {
	// Vérifier si pom.xml existe
	if _, err := os.Stat("pom.xml"); err == nil {
		return detectMavenFramework()
	}

	// TODO: Ajouter détection gradle, quarkus.*
	return "unknown", nil
}

// detectMavenFramework détecte le framework dans un projet Maven
func detectMavenFramework() (string, error) {
	data, err := os.ReadFile("pom.xml")
	if err != nil {
		return "", fmt.Errorf("erreur lors de la lecture du pom.xml: %w", err)
	}

	var project Project
	if err := xml.Unmarshal(data, &project); err != nil {
		return "", fmt.Errorf("erreur lors du parsing du pom.xml: %w", err)
	}

	// Recherche de dépendances Spring ou Quarkus
	for _, dep := range project.Dependencies.Dependencies {
		if strings.Contains(dep.GroupId, "spring") || strings.Contains(dep.ArtifactId, "spring") {
			return "spring", nil
		}
		if strings.Contains(dep.GroupId, "quarkus") || strings.Contains(dep.ArtifactId, "quarkus") {
			return "quarkus", nil
		}
	}

	// Framework Java non spécifiquement identifié
	return "java", nil
}
