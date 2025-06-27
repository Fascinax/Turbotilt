package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"turbotilt/internal/scan"
)

func TestGenerateComposeWithServices(t *testing.T) {
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

	// Options pour la génération
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

	// Générer le docker-compose.yml
	err = GenerateComposeWithServices(opts)
	if err != nil {
		t.Fatalf("GenerateComposeWithServices a retourné une erreur: %v", err)
	}

	// Vérifier que le fichier a été créé
	composePath := filepath.Join(tempDir, "docker-compose.yml")
	if _, err := os.Stat(composePath); err != nil {
		t.Fatalf("Le docker-compose.yml n'a pas été créé: %v", err)
	}

	// Lire le contenu du fichier généré
	content, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("Impossible de lire le docker-compose.yml généré: %v", err)
	}

	// Vérifier que le contenu est correct
	strContent := string(content)

	// Vérifier la présence du service principal
	if !strings.Contains(strContent, "api:") {
		t.Error("Le docker-compose.yml généré ne contient pas le service 'api'")
	}

	// Vérifier le port
	if !strings.Contains(strContent, "8080:8080") {
		t.Error("Le docker-compose.yml généré ne contient pas le mapping de port '8080:8080'")
	}

	// Vérifier la présence du service MySQL
	if !strings.Contains(strContent, "mysql:") {
		t.Error("Le docker-compose.yml généré ne contient pas le service 'mysql'")
	}
}

func TestGenerateComposeMultiService(t *testing.T) {
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

	// Créer une liste de services
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

	// Service base de données PostgreSQL
	pgService := scan.ServiceConfig{
		Type:    scan.PostgreSQL,
		Version: "14",
		Port:    "5432",
	}

	// Ajouter le service de base de données à la configuration
	serviceList.Services[0].Services = []scan.ServiceConfig{pgService}

	// Générer le docker-compose.yml pour les services multiples
	err = GenerateComposeMultiService(serviceList)
	if err != nil {
		t.Fatalf("GenerateComposeMultiService a retourné une erreur: %v", err)
	}

	// Vérifier que le fichier a été créé
	composePath := filepath.Join(tempDir, "docker-compose.yml")
	if _, err := os.Stat(composePath); err != nil {
		t.Fatalf("Le docker-compose.yml n'a pas été créé: %v", err)
	}

	// Lire le contenu du fichier généré
	content, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("Impossible de lire le docker-compose.yml généré: %v", err)
	}

	// Vérifier que le contenu est correct
	strContent := string(content)

	// Vérifier la présence des services
	if !strings.Contains(strContent, "service1:") {
		t.Error("Le docker-compose.yml généré ne contient pas le service 'service1'")
	}

	if !strings.Contains(strContent, "service2:") {
		t.Error("Le docker-compose.yml généré ne contient pas le service 'service2'")
	}

	if !strings.Contains(strContent, "postgres:") {
		t.Error("Le docker-compose.yml généré ne contient pas le service 'postgres'")
	}
}
