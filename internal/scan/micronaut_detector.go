package scan

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

// MicronautDetector est un détecteur spécifique pour les projets Micronaut
type MicronautDetector struct {
	BaseDetector
}

// Detect détecte si le projet utilise Micronaut
func (d *MicronautDetector) Detect(projectPath string) (bool, DetectionResult, error) {
	result := DetectionResult{
		Framework:   "micronaut",
		BuildSystem: "unknown",
	}

	// Vérifier pom.xml pour les projets Maven
	pomPath := filepath.Join(projectPath, "pom.xml")
	if _, err := os.Stat(pomPath); err == nil {
		// Lire le fichier pom.xml
		data, err := os.ReadFile(pomPath)
		if err != nil {
			return false, result, err
		}

		// Parseur simplifié pour pom.xml
		var pom struct {
			XMLName    xml.Name `xml:"project"`
			Parent     struct {
				GroupID    string `xml:"groupId"`
				ArtifactID string `xml:"artifactId"`
			} `xml:"parent"`
			Dependencies struct {
				Dependency []struct {
					GroupID    string `xml:"groupId"`
					ArtifactID string `xml:"artifactId"`
				} `xml:"dependency"`
			} `xml:"dependencies"`
		}

		if err := xml.Unmarshal(data, &pom); err != nil {
			return false, result, err
		}

		// Vérifier si c'est un projet Micronaut
		isMicronaut := false
		
		// Vérifier le parent
		if strings.Contains(pom.Parent.GroupID, "micronaut") || 
		   strings.Contains(pom.Parent.ArtifactID, "micronaut") {
			isMicronaut = true
		}
		
		// Vérifier les dépendances
		for _, dep := range pom.Dependencies.Dependency {
			if strings.Contains(dep.GroupID, "micronaut") {
				isMicronaut = true
				break
			}
		}

		if isMicronaut {
			result.BuildSystem = "maven"
			return true, result, nil
		}
	}

	// Vérifier build.gradle pour les projets Gradle
	gradlePath := filepath.Join(projectPath, "build.gradle")
	if _, err := os.Stat(gradlePath); err == nil {
		// Lire le fichier build.gradle
		data, err := os.ReadFile(gradlePath)
		if err != nil {
			return false, result, err
		}

		content := string(data)
		
		// Vérifier les références à Micronaut dans le fichier build.gradle
		if strings.Contains(content, "micronaut") || 
		   strings.Contains(content, "io.micronaut") {
			result.BuildSystem = "gradle"
			return true, result, nil
		}
	}

	// Vérifier la présence du fichier application.yml dans src/main/resources
	appYamlPath := filepath.Join(projectPath, "src", "main", "resources", "application.yml")
	if _, err := os.Stat(appYamlPath); err == nil {
		// Lire le fichier application.yml
		data, err := os.ReadFile(appYamlPath)
		if err != nil {
			return false, result, err
		}

		content := string(data)
		
		// Vérifier les références à Micronaut dans le fichier application.yml
		if strings.Contains(content, "micronaut:") {
			// Déterminer le système de build si possible
			if _, err := os.Stat(pomPath); err == nil {
				result.BuildSystem = "maven"
			} else if _, err := os.Stat(gradlePath); err == nil {
				result.BuildSystem = "gradle"
			}
			return true, result, nil
		}
	}

	// Vérifier la structure du package pour les classes typiques de Micronaut
	micronautClasses := []string{
		"Controller.java",
		"Resource.java",
		"Factory.java",
		"Client.java",
	}

	srcPath := filepath.Join(projectPath, "src", "main", "java")
	if _, err := os.Stat(srcPath); err == nil {
		err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".java") {
				// Vérifier si le fichier contient des imports Micronaut
				data, err := os.ReadFile(path)
				if err != nil {
					return nil // Continuer malgré l'erreur
				}
				
				content := string(data)
				if strings.Contains(content, "import io.micronaut.") {
					// Déterminer le système de build si possible
					if _, err := os.Stat(pomPath); err == nil {
						result.BuildSystem = "maven"
					} else if _, err := os.Stat(gradlePath); err == nil {
						result.BuildSystem = "gradle"
					}
					result.Detected = true
					return filepath.SkipAll // Sortir de la fonction Walk
				}
				
				// Vérifier les noms de fichiers typiques
				for _, className := range micronautClasses {
					if strings.HasSuffix(info.Name(), className) {
						data, err := os.ReadFile(path)
						if err != nil {
							continue
						}
						
						content := string(data)
						if strings.Contains(content, "import io.micronaut.") {
							// Déterminer le système de build si possible
							if _, err := os.Stat(pomPath); err == nil {
								result.BuildSystem = "maven"
							} else if _, err := os.Stat(gradlePath); err == nil {
								result.BuildSystem = "gradle"
							}
							result.Detected = true
							return filepath.SkipAll // Sortir de la fonction Walk
						}
					}
				}
			}
			return nil
		})
		
		if err != nil {
			return false, result, err
		}
		
		if result.Detected {
			return true, result, nil
		}
	}

	return false, result, nil
}
