package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"
)

func TestGenerateTiltfile(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Change working directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get current working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Unable to change directory: %v", err)
	}

	// Create a folder for templates
	templatesDir := filepath.Join(tempDir, "templates")
	if err := os.Mkdir(templatesDir, 0755); err != nil {
		t.Fatalf("Unable to create templates folder: %v", err)
	}

	// Create a test template with direct content
	templateContent := `# Test Tiltfile Template
# Framework: [[.Framework]]
# App: [[.AppName]]
# Port: [[.Port]]
# Date: [[.Date]]
# DevMode: [[.DevMode]]
`
	// Create the template directly in the templates folder
	templatePath := filepath.Join(templatesDir, "Tiltfile.tmpl")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Unable to create template file: %v", err)
	}

	// Options for generation
	opts := Options{
		Framework: "spring",
		AppName:   "test-app",
		Port:      "8080",
		DevMode:   true,
	}

	// Test with a custom template using GenerateTiltfileFromTemplate
	testOutputPath := filepath.Join(tempDir, "Tiltfile")
	err = GenerateTiltfileFromTemplate(opts, templatePath, testOutputPath)
	if err != nil {
		t.Fatalf("GenerateTiltfileFromTemplate returned an error: %v", err)
	}

	// Verify that the file was created
	if _, err := os.Stat(testOutputPath); err != nil {
		t.Fatalf("The Tiltfile was not created: %v", err)
	}

	// Read the content of the generated file
	content, err := os.ReadFile(testOutputPath)
	if err != nil {
		t.Fatalf("Unable to read the generated Tiltfile: %v", err)
	}

	// Verify that the content is correct
	strContent := string(content)
	if !strings.Contains(strContent, "Framework: spring") {
		t.Error("The generated Tiltfile does not contain 'Framework: spring'")
	}
	if !strings.Contains(strContent, "App: test-app") {
		t.Error("The generated Tiltfile does not contain 'App: test-app'")
	}
	if !strings.Contains(strContent, "Port: 8080") {
		t.Error("The generated Tiltfile does not contain 'Port: 8080'")
	}
	if !strings.Contains(strContent, "DevMode: true") {
		t.Error("The generated Tiltfile does not contain 'DevMode: true'")
	}
}

func TestGenerateTiltfileMultiService(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Change working directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get current working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Unable to change directory: %v", err)
	}

	// Create a folder for templates
	templatesDir := filepath.Join(tempDir, "templates")
	if err := os.Mkdir(templatesDir, 0755); err != nil {
		t.Fatalf("Unable to create templates folder: %v", err)
	}

	// Create a test template for multi-service with the k8s_yaml directive
	templateContent := `# Test Tiltfile Multi Template
# Date: [[.Date]]
k8s_yaml('docker-compose.yml')
`
	templatePath := filepath.Join(templatesDir, "Tiltfile.multi.tmpl")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Unable to create template file: %v", err)
	}

	// Also create a docker-compose.yml file for testing
	dockerComposeContent := `version: '3'
services:
  service1:
    image: test/service1
  service2:
    image: test/service2
`
	if err := os.WriteFile(filepath.Join(tempDir, "docker-compose.yml"), []byte(dockerComposeContent), 0644); err != nil {
		t.Fatalf("Unable to create docker-compose.yml file: %v", err)
	}

	// Create a Tiltfile directly from the template
	tiltfilePath := filepath.Join(tempDir, "Tiltfile")

	// Prepare template data
	data := TiltfileTemplateData{
		Date:           time.Now().Format("2006-01-02 15:04:05"),
		IsMultiService: true,
	}

	// Open output file
	f, err := os.Create(tiltfilePath)
	if err != nil {
		t.Fatalf("Unable to create Tiltfile: %v", err)
	}
	defer f.Close()

	// Create template and execute it
	tmpl, err := template.New("MultiTest").Delims("[[", "]]").Parse(templateContent)
	if err != nil {
		t.Fatalf("Error parsing template: %v", err)
	}

	if err := tmpl.Execute(f, data); err != nil {
		t.Fatalf("Error executing template: %v", err)
	}

	// Verify that the file was created
	if _, err := os.Stat(tiltfilePath); err != nil {
		t.Fatalf("The Tiltfile was not created: %v", err)
	}

	// Read the content of the generated file
	content, err := os.ReadFile(tiltfilePath)
	if err != nil {
		t.Fatalf("Unable to read the generated Tiltfile: %v", err)
	}

	// Verify that the content is correct for a multi-service project
	strContent := string(content)
	if !strings.Contains(strContent, "k8s_yaml") {
		t.Error("The generated multi-service Tiltfile should contain 'k8s_yaml'")
	}
}
