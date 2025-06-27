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
	// Créer un répertoire temporaire pour le test
	tempDir := t.TempDir()

	// Changer le répertoire de travail
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Impossible d'obtenir le répertoire de travail actuel: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Impossible de changer de répertoire: %v", err)
	}

	// Créer un dossier pour les templates
	templatesDir := filepath.Join(tempDir, "templates")
	if err := os.Mkdir(templatesDir, 0755); err != nil {
		t.Fatalf("Impossible de créer le dossier templates: %v", err)
	}

	// Créer un template de test avec un contenu direct
	templateContent := `# Test Tiltfile Template
# Framework: [[.Framework]]
# App: [[.AppName]]
# Port: [[.Port]]
# Date: [[.Date]]
# DevMode: [[.DevMode]]
`
	// Créer le template directement dans le dossier templates
	templatePath := filepath.Join(templatesDir, "Tiltfile.tmpl")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Impossible de créer le fichier de template: %v", err)
	}

	// Options pour la génération
	opts := Options{
		Framework: "spring",
		AppName:   "test-app",
		Port:      "8080",
		DevMode:   true,
	}

	// Test avec un template personnalisé en utilisant GenerateTiltfileFromTemplate
	testOutputPath := filepath.Join(tempDir, "Tiltfile")
	err = GenerateTiltfileFromTemplate(opts, templatePath, testOutputPath)
	if err != nil {
		t.Fatalf("GenerateTiltfileFromTemplate a retourné une erreur: %v", err)
	}

	// Vérifier que le fichier a été créé
	if _, err := os.Stat(testOutputPath); err != nil {
		t.Fatalf("Le Tiltfile n'a pas été créé: %v", err)
	}

	// Lire le contenu du fichier généré
	content, err := os.ReadFile(testOutputPath)
	if err != nil {
		t.Fatalf("Impossible de lire le Tiltfile généré: %v", err)
	}

	// Vérifier que le contenu est correct
	strContent := string(content)
	if !strings.Contains(strContent, "Framework: spring") {
		t.Error("Le Tiltfile généré ne contient pas 'Framework: spring'")
	}
	if !strings.Contains(strContent, "App: test-app") {
		t.Error("Le Tiltfile généré ne contient pas 'App: test-app'")
	}
	if !strings.Contains(strContent, "Port: 8080") {
		t.Error("Le Tiltfile généré ne contient pas 'Port: 8080'")
	}
	if !strings.Contains(strContent, "DevMode: true") {
		t.Error("Le Tiltfile généré ne contient pas 'DevMode: true'")
	}
}

func TestGenerateTiltfileMultiService(t *testing.T) {
	// Créer un répertoire temporaire pour le test
	tempDir := t.TempDir()

	// Changer le répertoire de travail
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Impossible d'obtenir le répertoire de travail actuel: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Impossible de changer de répertoire: %v", err)
	}

	// Créer un dossier pour les templates
	templatesDir := filepath.Join(tempDir, "templates")
	if err := os.Mkdir(templatesDir, 0755); err != nil {
		t.Fatalf("Impossible de créer le dossier templates: %v", err)
	}

	// Créer un template de test pour multi-service avec la directive k8s_yaml
	templateContent := `# Test Tiltfile Multi Template
# Date: [[.Date]]
k8s_yaml('docker-compose.yml')
`
	templatePath := filepath.Join(templatesDir, "Tiltfile.multi.tmpl")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Impossible de créer le fichier de template: %v", err)
	}

	// Créer également un fichier docker-compose.yml pour le test
	dockerComposeContent := `version: '3'
services:
  service1:
    image: test/service1
  service2:
    image: test/service2
`
	if err := os.WriteFile(filepath.Join(tempDir, "docker-compose.yml"), []byte(dockerComposeContent), 0644); err != nil {
		t.Fatalf("Impossible de créer le fichier docker-compose.yml: %v", err)
	}

	// Créer un Tiltfile directement depuis le template
	tiltfilePath := filepath.Join(tempDir, "Tiltfile")

	// Préparer les données du template
	data := TiltfileTemplateData{
		Date:           time.Now().Format("2006-01-02 15:04:05"),
		IsMultiService: true,
	}

	// Ouvrir le fichier de sortie
	f, err := os.Create(tiltfilePath)
	if err != nil {
		t.Fatalf("Impossible de créer le Tiltfile: %v", err)
	}
	defer f.Close()

	// Créer le template et l'exécuter
	tmpl, err := template.New("MultiTest").Delims("[[", "]]").Parse(templateContent)
	if err != nil {
		t.Fatalf("Erreur lors du parsing du template: %v", err)
	}

	if err := tmpl.Execute(f, data); err != nil {
		t.Fatalf("Erreur lors de l'exécution du template: %v", err)
	}

	// Vérifier que le fichier a été créé
	if _, err := os.Stat(tiltfilePath); err != nil {
		t.Fatalf("Le Tiltfile n'a pas été créé: %v", err)
	}

	// Lire le contenu du fichier généré
	content, err := os.ReadFile(tiltfilePath)
	if err != nil {
		t.Fatalf("Impossible de lire le Tiltfile généré: %v", err)
	}

	// Vérifier que le contenu est correct pour un projet multi-service
	strContent := string(content)
	if !strings.Contains(strContent, "k8s_yaml") {
		t.Error("Le Tiltfile multi-service généré devrait contenir 'k8s_yaml'")
	}
}
