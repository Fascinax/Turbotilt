package scan

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	XMLName      xml.Name `xml:"project"`
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

// Scanner struct for detecting frameworks and services
type Scanner struct {
	ProjectPath string
}

// NewScanner creates a new Scanner instance for the given project path
func NewScanner(projectPath string) *Scanner {
	return &Scanner{
		ProjectPath: projectPath,
	}
}

// DetectFramework detects the framework used in the project
func (s *Scanner) DetectFramework() (string, *DetectionResult, error) {
	// Create instances of detectors
	micronautDetector := &MicronautDetector{}

	// First check for Micronaut framework
	detected, result, err := micronautDetector.Detect(s.ProjectPath)
	if err == nil && detected {
		return result.Framework, &result, nil
	}

	// If Micronaut not detected, fall back to simpler detection methods
	// Check for pom.xml (Maven)
	pomPath := filepath.Join(s.ProjectPath, "pom.xml")
	if _, err := os.Stat(pomPath); err == nil {
		data, err := os.ReadFile(pomPath)
		if err == nil {
			content := string(data)
			// Check for Spring Boot
			if strings.Contains(content, "spring-boot") {
				result := DetectionResult{
					Framework:   "spring",
					BuildSystem: "maven",
					Detected:    true,
				}
				return "spring", &result, nil
			}
			// Check for Quarkus
			if strings.Contains(content, "quarkus") {
				result := DetectionResult{
					Framework:   "quarkus",
					BuildSystem: "maven",
					Detected:    true,
				}
				return "quarkus", &result, nil
			}
		}
	}

	// Check for build.gradle (Gradle)
	gradlePath := filepath.Join(s.ProjectPath, "build.gradle")
	if _, err := os.Stat(gradlePath); err == nil {
		data, err := os.ReadFile(gradlePath)
		if err == nil {
			content := string(data)
			// Check for Spring Boot
			if strings.Contains(content, "spring-boot") {
				result := DetectionResult{
					Framework:   "spring",
					BuildSystem: "gradle",
					Detected:    true,
				}
				return "spring", &result, nil
			}
			// Check for Quarkus
			if strings.Contains(content, "quarkus") {
				result := DetectionResult{
					Framework:   "quarkus",
					BuildSystem: "gradle",
					Detected:    true,
				}
				return "quarkus", &result, nil
			}
		}
	}

	// No framework detected
	return "", nil, fmt.Errorf("no framework detected")
}

// DetectFramework detects the framework used in the project
func DetectFramework() (string, error) {
	// Maven detection (pom.xml)
	if _, err := os.Stat("pom.xml"); err == nil {
		return detectMavenFramework()
	}

	// Gradle detection
	gradleFiles := []string{"build.gradle", "build.gradle.kts"}
	for _, file := range gradleFiles {
		if _, err := os.Stat(file); err == nil {
			return detectGradleFramework(file)
		}
	}

	// Quarkus-specific file detection
	quarkusFiles := []string{
		"src/main/resources/application.properties",
		".quarkus",
		"src/main/resources/META-INF/resources/index.html", // Often present in Quarkus projects
	}

	for _, file := range quarkusFiles {
		if _, err := os.Stat(file); err == nil {
			// Check if it contains Quarkus content
			content, err := os.ReadFile(file)
			if err == nil && strings.Contains(string(content), "quarkus") {
				return "quarkus", nil
			}
		}
	}

	// Detection by project structure
	if isSpringProjectStructure() {
		return "spring", nil
	}

	// Framework not detected
	return "unknown", nil
}

// detectMavenFramework detects the framework in a Maven project
func detectMavenFramework() (string, error) {
	data, err := os.ReadFile("pom.xml")
	if err != nil {
		return "", fmt.Errorf("error reading pom.xml: %w", err)
	}

	var project Project
	if err := xml.Unmarshal(data, &project); err != nil {
		return "", fmt.Errorf("error parsing pom.xml: %w", err)
	}

	// Check the parent (common case for Spring Boot)
	if strings.Contains(project.Parent.GroupId, "spring") || strings.Contains(project.Parent.ArtifactId, "spring") {
		return "spring", nil
	}

	// Search for Spring or Quarkus dependencies
	for _, dep := range project.Dependencies.Dependencies {
		if strings.Contains(dep.GroupId, "spring") || strings.Contains(dep.ArtifactId, "spring") {
			return "spring", nil
		}
		if strings.Contains(dep.GroupId, "quarkus") || strings.Contains(dep.ArtifactId, "quarkus") {
			return "quarkus", nil
		}
		// Check for Micronaut or other frameworks
		if strings.Contains(dep.GroupId, "micronaut") || strings.Contains(dep.ArtifactId, "micronaut") {
			return "micronaut", nil
		}
	}

	// Java framework not specifically identified
	return "java", nil
}

// detectGradleFramework detects the framework in a Gradle project
func detectGradleFramework(gradleFile string) (string, error) {
	data, err := os.ReadFile(gradleFile)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %w", gradleFile, err)
	}

	content := string(data)

	// Search for framework-specific indicators
	if strings.Contains(content, "org.springframework") || strings.Contains(content, "spring-boot") {
		return "spring", nil
	}
	if strings.Contains(content, "io.quarkus") || strings.Contains(content, "quarkus-gradle-plugin") {
		return "quarkus", nil
	}
	if strings.Contains(content, "io.micronaut") {
		return "micronaut", nil
	}

	// Java framework not specifically identified
	return "java", nil
}

// isSpringProjectStructure checks if the project structure matches Spring
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
			// Check if the found files contain Spring indicators
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
