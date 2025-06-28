package scan

import (
	"os"
	"path/filepath"
)

// Service represents a detected microservice
type Service struct {
	Name string
	Path string
	Type string
}

// Detector helps find microservices in a directory structure
type Detector struct {
	// Add any configuration options here
}

// NewDetector creates a new detector instance
func NewDetector() *Detector {
	return &Detector{}
}

// ScanDirectory scans a directory for microservices
func (d *Detector) ScanDirectory(rootPath string) ([]Service, error) {
	services := []Service{}

	// Walk the directory
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == rootPath {
			return nil
		}

		// Only check directories
		if !info.IsDir() {
			return nil
		}

		// Skip hidden directories and "target", "build", "node_modules", etc.
		if info.Name()[0] == '.' ||
			info.Name() == "target" ||
			info.Name() == "build" ||
			info.Name() == "node_modules" ||
			info.Name() == "dist" {
			return filepath.SkipDir
		}

		// Check if this directory contains a service
		serviceType, isService := d.detectService(path)
		if isService {
			// Create a service entry
			service := Service{
				Name: filepath.Base(path),
				Path: path,
				Type: serviceType,
			}
			services = append(services, service)

			// Skip subdirectories of a service
			return filepath.SkipDir
		}

		return nil
	})

	return services, err
}

// detectService checks if a directory contains a microservice
func (d *Detector) detectService(path string) (string, bool) {
	// Check for Maven project (pom.xml)
	if _, err := os.Stat(filepath.Join(path, "pom.xml")); err == nil {
		// Use Scanner to detect specific framework
		scanner := NewScanner(path)
		framework, _, err := scanner.DetectFramework()
		if err == nil && framework != "" {
			return framework, true
		}
		return "java", true
	}

	// Check for Gradle project (build.gradle or build.gradle.kts)
	if _, err := os.Stat(filepath.Join(path, "build.gradle")); err == nil {
		return "gradle", true
	}
	if _, err := os.Stat(filepath.Join(path, "build.gradle.kts")); err == nil {
		return "gradle", true
	}

	// Check for Node.js project (package.json)
	if _, err := os.Stat(filepath.Join(path, "package.json")); err == nil {
		// Check for Angular
		if _, err := os.Stat(filepath.Join(path, "angular.json")); err == nil {
			return "angular", true
		}
		// Check for React
		if _, err := os.Stat(filepath.Join(path, "react-scripts")); err == nil {
			return "react", true
		}
		return "nodejs", true
	}

	// Check for Python project (requirements.txt, setup.py)
	if _, err := os.Stat(filepath.Join(path, "requirements.txt")); err == nil {
		return "python", true
	}
	if _, err := os.Stat(filepath.Join(path, "setup.py")); err == nil {
		return "python", true
	}

	// Check for Go project (go.mod)
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		return "go", true
	}

	// Check for Docker project (Dockerfile or docker-compose.yml)
	if _, err := os.Stat(filepath.Join(path, "Dockerfile")); err == nil {
		return "docker", true
	}
	if _, err := os.Stat(filepath.Join(path, "docker-compose.yml")); err == nil {
		return "docker-compose", true
	}

	// Not a recognized service
	return "", false
}
