package scan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMicronautDetector(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tempDir, err := os.MkdirTemp("", "micronaut-detector-test")
	if err != nil {
		t.Fatalf("Impossible de créer le répertoire temporaire: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Structure de répertoire pour un projet Micronaut
	srcMainJava := filepath.Join(tempDir, "src", "main", "java", "com", "example")
	srcMainResources := filepath.Join(tempDir, "src", "main", "resources")

	// Créer les répertoires
	if err := os.MkdirAll(srcMainJava, 0755); err != nil {
		t.Fatalf("Impossible de créer les répertoires: %v", err)
	}
	if err := os.MkdirAll(srcMainResources, 0755); err != nil {
		t.Fatalf("Impossible de créer les répertoires: %v", err)
	}

	t.Run("Test détection via pom.xml", func(t *testing.T) {
		// Créer un pom.xml factice avec des dépendances Micronaut
		pomXml := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>micronaut-test</artifactId>
    <version>0.1.0</version>
    <parent>
        <groupId>io.micronaut.parent</groupId>
        <artifactId>micronaut-parent</artifactId>
        <version>3.5.0</version>
    </parent>
    <dependencies>
        <dependency>
            <groupId>io.micronaut</groupId>
            <artifactId>micronaut-inject</artifactId>
        </dependency>
    </dependencies>
</project>`

		if err := os.WriteFile(filepath.Join(tempDir, "pom.xml"), []byte(pomXml), 0644); err != nil {
			t.Fatalf("Impossible d'écrire le fichier pom.xml: %v", err)
		}

		detector := MicronautDetector{}
		detected, result, err := detector.Detect(tempDir)

		if err != nil {
			t.Errorf("Erreur lors de la détection: %v", err)
		}

		if !detected {
			t.Errorf("Le détecteur n'a pas détecté Micronaut via pom.xml")
		}

		if result.Framework != "micronaut" {
			t.Errorf("Framework détecté incorrect: %s", result.Framework)
		}

		if result.BuildSystem != "maven" {
			t.Errorf("Système de build détecté incorrect: %s", result.BuildSystem)
		}

		// Supprimer le pom.xml pour le prochain test
		os.Remove(filepath.Join(tempDir, "pom.xml"))
	})

	t.Run("Test détection via build.gradle", func(t *testing.T) {
		// Créer un build.gradle factice avec des dépendances Micronaut
		buildGradle := `plugins {
    id 'io.micronaut.application' version '3.5.0'
}

repositories {
    mavenCentral()
}

dependencies {
    implementation 'io.micronaut:micronaut-http-server-netty'
    implementation 'io.micronaut:micronaut-inject'
}
`

		if err := os.WriteFile(filepath.Join(tempDir, "build.gradle"), []byte(buildGradle), 0644); err != nil {
			t.Fatalf("Impossible d'écrire le fichier build.gradle: %v", err)
		}

		detector := MicronautDetector{}
		detected, result, err := detector.Detect(tempDir)

		if err != nil {
			t.Errorf("Erreur lors de la détection: %v", err)
		}

		if !detected {
			t.Errorf("Le détecteur n'a pas détecté Micronaut via build.gradle")
		}

		if result.Framework != "micronaut" {
			t.Errorf("Framework détecté incorrect: %s", result.Framework)
		}

		if result.BuildSystem != "gradle" {
			t.Errorf("Système de build détecté incorrect: %s", result.BuildSystem)
		}

		// Supprimer le build.gradle pour le prochain test
		os.Remove(filepath.Join(tempDir, "build.gradle"))
	})

	t.Run("Test détection via application.yml", func(t *testing.T) {
		// Créer un application.yml factice
		appYaml := `micronaut:
  application:
    name: test-app
  server:
    port: 8080
`

		if err := os.WriteFile(filepath.Join(srcMainResources, "application.yml"), []byte(appYaml), 0644); err != nil {
			t.Fatalf("Impossible d'écrire le fichier application.yml: %v", err)
		}

		// Ajout d'un build.gradle pour la détection du système de build
		if err := os.WriteFile(filepath.Join(tempDir, "build.gradle"), []byte("// Empty build file"), 0644); err != nil {
			t.Fatalf("Impossible d'écrire le fichier build.gradle: %v", err)
		}

		detector := MicronautDetector{}
		detected, result, err := detector.Detect(tempDir)

		if err != nil {
			t.Errorf("Erreur lors de la détection: %v", err)
		}

		if !detected {
			t.Errorf("Le détecteur n'a pas détecté Micronaut via application.yml")
		}

		if result.Framework != "micronaut" {
			t.Errorf("Framework détecté incorrect: %s", result.Framework)
		}

		if result.BuildSystem != "gradle" {
			t.Errorf("Système de build détecté incorrect: %s", result.BuildSystem)
		}

		// Nettoyage
		os.Remove(filepath.Join(srcMainResources, "application.yml"))
		os.Remove(filepath.Join(tempDir, "build.gradle"))
	})

	t.Run("Test détection via imports dans les fichiers Java", func(t *testing.T) {
		// Créer un contrôleur Java avec des imports Micronaut
		javaFile := `package com.example;

import io.micronaut.http.annotation.Controller;
import io.micronaut.http.annotation.Get;
import io.micronaut.http.HttpResponse;

@Controller("/api")
public class TestController {
    @Get("/")
    public HttpResponse<?> index() {
        return HttpResponse.ok("Hello World");
    }
}
`

		if err := os.WriteFile(filepath.Join(srcMainJava, "TestController.java"), []byte(javaFile), 0644); err != nil {
			t.Fatalf("Impossible d'écrire le fichier Java: %v", err)
		}

		// Ajout d'un pom.xml pour la détection du système de build
		if err := os.WriteFile(filepath.Join(tempDir, "pom.xml"), []byte("<project></project>"), 0644); err != nil {
			t.Fatalf("Impossible d'écrire le fichier pom.xml: %v", err)
		}

		detector := MicronautDetector{}
		detected, result, err := detector.Detect(tempDir)

		if err != nil {
			t.Errorf("Erreur lors de la détection: %v", err)
		}

		if !detected {
			t.Errorf("Le détecteur n'a pas détecté Micronaut via imports Java")
		}

		if result.Framework != "micronaut" {
			t.Errorf("Framework détecté incorrect: %s", result.Framework)
		}

		if result.BuildSystem != "maven" {
			t.Errorf("Système de build détecté incorrect: %s", result.BuildSystem)
		}
	})
}
