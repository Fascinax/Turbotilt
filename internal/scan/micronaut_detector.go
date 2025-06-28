package scan

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

// MicronautDetector is a specific detector for Micronaut projects
type MicronautDetector struct {
	BaseDetector
}

// Detect detects if the project uses Micronaut
func (d *MicronautDetector) Detect(projectPath string) (bool, DetectionResult, error) {
	result := DetectionResult{
		Framework:   "micronaut",
		BuildSystem: "unknown",
	}

	// Check pom.xml for Maven projects
	pomPath := filepath.Join(projectPath, "pom.xml")
	if _, err := os.Stat(pomPath); err == nil {
		// Read pom.xml file
		data, err := os.ReadFile(pomPath)
		if err != nil {
			return false, result, err
		}

		// Simplified parser for pom.xml
		var pom struct {
			XMLName xml.Name `xml:"project"`
			Parent  struct {
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

		// Check if it's a Micronaut project
		isMicronaut := false

		// Check the parent
		if strings.Contains(pom.Parent.GroupID, "micronaut") ||
			strings.Contains(pom.Parent.ArtifactID, "micronaut") {
			isMicronaut = true
		}

		// Check dependencies
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

	// Check build.gradle for Gradle projects
	gradlePath := filepath.Join(projectPath, "build.gradle")
	if _, err := os.Stat(gradlePath); err == nil {
		// Read the build.gradle file
		data, err := os.ReadFile(gradlePath)
		if err != nil {
			return false, result, err
		}

		content := string(data)

		// Check for Micronaut references in the build.gradle file
		if strings.Contains(content, "micronaut") ||
			strings.Contains(content, "io.micronaut") {
			result.BuildSystem = "gradle"
			return true, result, nil
		}
	}

	// Check for the presence of application.yml in src/main/resources
	appYamlPath := filepath.Join(projectPath, "src", "main", "resources", "application.yml")
	if _, err := os.Stat(appYamlPath); err == nil {
		// Read the application.yml file
		data, err := os.ReadFile(appYamlPath)
		if err != nil {
			return false, result, err
		}

		content := string(data)

		// Check for Micronaut references in the application.yml file
		if strings.Contains(content, "micronaut:") {
			// Determine build system if possible
			if _, err := os.Stat(pomPath); err == nil {
				result.BuildSystem = "maven"
			} else if _, err := os.Stat(gradlePath); err == nil {
				result.BuildSystem = "gradle"
			}
			return true, result, nil
		}
	}

	// Check package structure for typical Micronaut classes
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
				// Check if the file contains Micronaut imports
				data, err := os.ReadFile(path)
				if err != nil {
					return nil // Continue despite the error
				}

				content := string(data)
				if strings.Contains(content, "import io.micronaut.") {
					// Determine build system if possible
					if _, err := os.Stat(pomPath); err == nil {
						result.BuildSystem = "maven"
					} else if _, err := os.Stat(gradlePath); err == nil {
						result.BuildSystem = "gradle"
					}
					result.Detected = true
					return filepath.SkipAll // Exit the Walk function
				}

				// Check for typical filenames
				for _, className := range micronautClasses {
					if strings.HasSuffix(info.Name(), className) {
						data, err := os.ReadFile(path)
						if err != nil {
							continue
						}

						content := string(data)
						if strings.Contains(content, "import io.micronaut.") {
							// Determine the build system if possible
							if _, err := os.Stat(pomPath); err == nil {
								result.BuildSystem = "maven"
							} else if _, err := os.Stat(gradlePath); err == nil {
								result.BuildSystem = "gradle"
							}
							result.Detected = true
							return filepath.SkipAll // Exit the Walk function
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
