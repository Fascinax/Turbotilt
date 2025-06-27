package scan

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDetectFramework teste la fonctionnalité de détection de framework
func TestDetectFramework(t *testing.T) {
	// Création d'un répertoire temporaire pour les tests
	tempDir, err := os.MkdirTemp("", "turbotilt-test-*")
	if err != nil {
		t.Fatalf("Impossible de créer le répertoire temporaire: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Sauvegarde du répertoire de travail actuel
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Impossible d'obtenir le répertoire courant: %v", err)
	}

	// Changer au répertoire temporaire pour les tests
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Impossible de changer de répertoire: %v", err)
	}
	defer os.Chdir(originalDir) // Restaurer le répertoire de travail à la fin

	// Test 1: Projet Maven Spring Boot
	t.Run("Spring Boot Maven Project", func(t *testing.T) {
		// Créer un faux pom.xml pour Spring Boot
		pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" 
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>2.7.0</version>
    </parent>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>0.0.1-SNAPSHOT</version>
    <name>demo</name>
    <description>Demo project for Spring Boot</description>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
    </dependencies>
</project>`

		if err := os.WriteFile("pom.xml", []byte(pomContent), 0644); err != nil {
			t.Fatalf("Erreur lors de la création du fichier pom.xml: %v", err)
		}

		// Tester la détection
		framework, err := DetectFramework()
		if err != nil {
			t.Errorf("Erreur lors de la détection du framework: %v", err)
		}
		if framework != "spring" {
			t.Errorf("Framework détecté incorrect. Attendu: %s, Obtenu: %s", "spring", framework)
		}

		// Nettoyer
		os.Remove("pom.xml")
	})

	// Test 2: Projet Quarkus Maven
	t.Run("Quarkus Maven Project", func(t *testing.T) {
		// Créer un faux pom.xml pour Quarkus
		pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>quarkus-demo</artifactId>
    <version>1.0.0-SNAPSHOT</version>
    <dependencies>
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-core</artifactId>
            <version>2.0.0.Final</version>
        </dependency>
    </dependencies>
</project>`

		if err := os.WriteFile("pom.xml", []byte(pomContent), 0644); err != nil {
			t.Fatalf("Erreur lors de la création du fichier pom.xml: %v", err)
		}

		// Tester la détection
		framework, err := DetectFramework()
		if err != nil {
			t.Errorf("Erreur lors de la détection du framework: %v", err)
		}
		if framework != "quarkus" {
			t.Errorf("Framework détecté incorrect. Attendu: %s, Obtenu: %s", "quarkus", framework)
		}

		// Nettoyer
		os.Remove("pom.xml")
	})

	// Test 3: Projet Gradle
	t.Run("Gradle Project", func(t *testing.T) {
		// Créer un faux build.gradle pour Spring Boot
		gradleContent := `
plugins {
    id 'org.springframework.boot' version '2.7.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
    id 'java'
}

group = 'com.example'
version = '0.0.1-SNAPSHOT'
sourceCompatibility = '17'

repositories {
    mavenCentral()
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}
`
		if err := os.WriteFile("build.gradle", []byte(gradleContent), 0644); err != nil {
			t.Fatalf("Erreur lors de la création du fichier build.gradle: %v", err)
		}

		// Tester la détection
		framework, err := DetectFramework()
		if err != nil {
			t.Errorf("Erreur lors de la détection du framework: %v", err)
		}
		// Selon votre implémentation actuelle pour Gradle
		// On s'attendrait à ce que le framework soit spring ou java
		expectedFrameworks := map[string]bool{
			"spring": true,
			"java":   true,
		}
		if !expectedFrameworks[framework] {
			t.Errorf("Framework détecté incorrect. Attendu: spring ou java, Obtenu: %s", framework)
		}

		// Nettoyer
		os.Remove("build.gradle")
	})

	// Test 4: Projet inconnu
	t.Run("Unknown Project", func(t *testing.T) {
		// Ne pas créer de fichiers de projet

		// Tester la détection
		framework, err := DetectFramework()
		if err != nil {
			t.Errorf("Erreur lors de la détection du framework: %v", err)
		}
		if framework != "unknown" {
			t.Errorf("Framework détecté incorrect. Attendu: %s, Obtenu: %s", "unknown", framework)
		}
	})
}

// TestDetectors teste que tous les détecteurs sont correctement appelés
func TestDetectors(t *testing.T) {
	testCases := []struct {
		name         string
		mockFiles    map[string]string
		expectedType string
	}{
		{
			name: "Spring Boot",
			mockFiles: map[string]string{
				"pom.xml": `<project>
					<dependencies>
						<dependency>
							<groupId>org.springframework.boot</groupId>
							<artifactId>spring-boot-starter</artifactId>
						</dependency>
					</dependencies>
				</project>`,
				"src/main/java/com/example/App.java": `
					package com.example;
					import org.springframework.boot.SpringApplication;
					import org.springframework.boot.autoconfigure.SpringBootApplication;
					
					@SpringBootApplication
					public class App {
						public static void main(String[] args) {
							SpringApplication.run(App.class, args);
						}
					}
				`,
			},
			expectedType: "spring",
		},
		{
			name: "Quarkus",
			mockFiles: map[string]string{
				"pom.xml": `<project>
					<dependencies>
						<dependency>
							<groupId>io.quarkus</groupId>
							<artifactId>quarkus-core</artifactId>
						</dependency>
					</dependencies>
				</project>`,
				"src/main/java/com/example/App.java": `
					package com.example;
					import io.quarkus.runtime.Quarkus;
					
					public class App {
						public static void main(String[] args) {
							Quarkus.run(args);
						}
					}
				`,
			},
			expectedType: "quarkus",
		},
		{
			name: "Micronaut",
			mockFiles: map[string]string{
				"pom.xml": `<project>
					<dependencies>
						<dependency>
							<groupId>io.micronaut</groupId>
							<artifactId>micronaut-core</artifactId>
						</dependency>
					</dependencies>
				</project>`,
				"src/main/java/com/example/App.java": `
					package com.example;
					import io.micronaut.runtime.Micronaut;
					
					public class App {
						public static void main(String[] args) {
							Micronaut.run(App.class, args);
						}
					}
				`,
			},
			expectedType: "micronaut",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Créer un répertoire temporaire
			tempDir, err := os.MkdirTemp("", "detector-test")
			if err != nil {
				t.Fatalf("Erreur lors de la création du répertoire temporaire: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Créer les fichiers simulés
			for path, content := range tc.mockFiles {
				fullPath := filepath.Join(tempDir, path)
				dirPath := filepath.Dir(fullPath)

				if err := os.MkdirAll(dirPath, 0755); err != nil {
					t.Fatalf("Erreur lors de la création du répertoire %s: %v", dirPath, err)
				}

				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("Erreur lors de la création du fichier %s: %v", fullPath, err)
				}
			}

			// Exécuter le scanner
			scanner := NewScanner(tempDir)
			framework, result, err := scanner.DetectFramework()

			if err != nil {
				t.Errorf("Erreur lors de la détection: %v", err)
			}

			if framework != tc.expectedType {
				t.Errorf("Framework détecté incorrect: %s, attendu: %s", framework, tc.expectedType)
			}

			if result.Framework != tc.expectedType {
				t.Errorf("Framework dans le résultat incorrect: %s, attendu: %s", result.Framework, tc.expectedType)
			}
		})
	}
}
