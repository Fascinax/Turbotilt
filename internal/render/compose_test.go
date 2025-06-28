package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"turbotilt/internal/scan"
)

func TestGenerateComposeWithServices(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Change the working directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get current working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Unable to change directory: %v", err)
	}

	// Options for generation
	opts := Options{
		Framework:   "spring",
		ServiceName: "api",
		Port:        "8080",
		Path:        ".",
		Services: []scan.ServiceConfig{
			{
				Type:    "mysql",
				Version: "8.0",
				Credentials: map[string]string{
					"username": "root",
					"password": "password",
				},
				Port: "3306",
			},
		},
	}

	// Generate docker-compose.yml
	err = GenerateComposeWithServices(opts)
	if err != nil {
		t.Fatalf("GenerateComposeWithServices returned an error: %v", err)
	}

	// Verify the file was created
	composePath := filepath.Join(tempDir, "docker-compose.yml")
	if _, err := os.Stat(composePath); err != nil {
		t.Fatalf("docker-compose.yml was not created: %v", err)
	}

	// Read the content of the generated file
	content, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("Unable to read the generated docker-compose.yml: %v", err)
	}

	// Verify the content is correct
	strContent := string(content)

	// Verify the main service is present
	if !strings.Contains(strContent, "api:") {
		t.Error("The generated docker-compose.yml does not contain the 'api' service")
	}

	// Verify the port
	if !strings.Contains(strContent, "8080:8080") {
		t.Error("The generated docker-compose.yml does not contain the port mapping '8080:8080'")
	}

	// Verify the MySQL service is present
	if !strings.Contains(strContent, "mysql:") {
		t.Error("The generated docker-compose.yml does not contain the 'mysql' service")
	}
}

func TestGenerateComposeMultiService(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Change the working directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get current working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Unable to change directory: %v", err)
	}

	// Create a list of services
	serviceList := ServiceList{
		Services: []Options{
			{
				ServiceName: "service1",
				Path:        "./service1",
				Framework:   "spring",
				Port:        "8080",
				DevMode:     true,
			},
			{
				ServiceName: "service2",
				Path:        "./service2",
				Framework:   "quarkus",
				Port:        "8081",
				DevMode:     false,
			},
		},
	}

	// PostgreSQL database service
	pgService := scan.ServiceConfig{
		Type:    scan.PostgreSQL,
		Version: "14",
		Port:    "5432",
	}

	// Add database service to the configuration
	serviceList.Services[0].Services = []scan.ServiceConfig{pgService}

	// Generate docker-compose.yml for multiple services
	err = GenerateComposeMultiService(serviceList)
	if err != nil {
		t.Fatalf("GenerateComposeMultiService returned an error: %v", err)
	}

	// Verify the file was created
	composePath := filepath.Join(tempDir, "docker-compose.yml")
	if _, err := os.Stat(composePath); err != nil {
		t.Fatalf("docker-compose.yml was not created: %v", err)
	}

	// Read the content of the generated file
	content, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("Unable to read the generated docker-compose.yml: %v", err)
	}

	// Verify the content is correct
	strContent := string(content)

	// Verify the services are present
	if !strings.Contains(strContent, "service1:") {
		t.Error("The generated docker-compose.yml does not contain the 'service1' service")
	}

	if !strings.Contains(strContent, "service2:") {
		t.Error("The generated docker-compose.yml does not contain the 'service2' service")
	}

	if !strings.Contains(strContent, "postgres:") {
		t.Error("The generated docker-compose.yml does not contain the 'postgres' service")
	}
}
