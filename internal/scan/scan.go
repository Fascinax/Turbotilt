package scan

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
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
	Parent struct {
		GroupId    string `xml:"groupId"`
		ArtifactId string `xml:"artifactId"`
	} `xml:"parent"`
}

type GradleBuild struct {
	Content string
}

// DetectFramework détecte le framework utilisé dans le projet
func DetectFramework() (string, error) {
	// Détection Maven (pom.xml)
	if _, err := os.Stat("pom.xml"); err == nil {
		return detectMavenFramework()
	}

	// Détection Gradle
	gradleFiles := []string{"build.gradle", "build.gradle.kts"}
	for _, file := range gradleFiles {
		if _, err := os.Stat(file); err == nil {
			return detectGradleFramework(file)
		}
	}

	// Détection fichiers Quarkus spécifiques
	quarkusFiles := []string{
		"src/main/resources/application.properties",
		".quarkus",
		"src/main/resources/META-INF/resources/index.html", // Souvent présent dans les projets Quarkus
	}

	for _, file := range quarkusFiles {
		if _, err := os.Stat(file); err == nil {
			// Vérifier si contient du contenu Quarkus
			content, err := os.ReadFile(file)
			if err == nil && strings.Contains(string(content), "quarkus") {
				return "quarkus", nil
			}
		}
	}

	// Détection par structure de projet
	if isSpringProjectStructure() {
		return "spring", nil
	}

	// Framework non détecté
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

	// Vérifier le parent (cas courant pour Spring Boot)
	if strings.Contains(project.Parent.GroupId, "spring") || strings.Contains(project.Parent.ArtifactId, "spring") {
		return "spring", nil
	}

	// Recherche de dépendances Spring ou Quarkus
	for _, dep := range project.Dependencies.Dependencies {
		if strings.Contains(dep.GroupId, "spring") || strings.Contains(dep.ArtifactId, "spring") {
			return "spring", nil
		}
		if strings.Contains(dep.GroupId, "quarkus") || strings.Contains(dep.ArtifactId, "quarkus") {
			return "quarkus", nil
		}
		// Vérification pour Micronaut ou d'autres frameworks
		if strings.Contains(dep.GroupId, "micronaut") || strings.Contains(dep.ArtifactId, "micronaut") {
			return "micronaut", nil
		}
	}

	// Framework Java non spécifiquement identifié
	return "java", nil
}

// detectGradleFramework détecte le framework dans un projet Gradle
func detectGradleFramework(gradleFile string) (string, error) {
	data, err := os.ReadFile(gradleFile)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la lecture de %s: %w", gradleFile, err)
	}

	content := string(data)

	// Recherche d'indices spécifiques au framework
	if strings.Contains(content, "org.springframework") || strings.Contains(content, "spring-boot") {
		return "spring", nil
	}
	if strings.Contains(content, "io.quarkus") || strings.Contains(content, "quarkus-gradle-plugin") {
		return "quarkus", nil
	}
	if strings.Contains(content, "io.micronaut") {
		return "micronaut", nil
	}

	// Framework Java non spécifiquement identifié
	return "java", nil
}

// isSpringProjectStructure vérifie si la structure du projet correspond à Spring
func isSpringProjectStructure() bool {
	springPatterns := []string{
		"src/main/java/*/Application.java",
		"src/main/java/*/SpringBootApplication.java",
		"src/main/java/**/*Controller.java",
		"src/main/resources/application.yml",
		"src/main/resources/application.properties",
	}

	for _, pattern := range springPatterns {
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			// Vérifier si les fichiers trouvés contiennent des indications Spring
			for _, file := range matches {
				content, err := os.ReadFile(file)
				if err == nil {
					if strings.Contains(string(content), "springframework") ||
						strings.Contains(string(content), "@SpringBootApplication") ||
						strings.Contains(string(content), "@RestController") {
						return true
					}
				}
			}
		}
	}

	return false
}
