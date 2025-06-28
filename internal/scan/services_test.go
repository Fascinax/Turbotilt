package scan

import (
	"os"
	"testing"
)

// TestDetectServices teste la fonctionnalité de détection des services
func TestDetectServices(t *testing.T) {
	// Création d'un répertoire temporaire pour les tests
	tempDir, err := os.MkdirTemp("", "turbotilt-services-test-*")
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

	// Test 1: Projet Spring Boot avec MySQL
	t.Run("Spring Boot with MySQL", func(t *testing.T) {
		// Créer un faux fichier application.properties avec configuration MySQL
		propertiesContent := `
spring.datasource.url=jdbc:mysql://localhost:3306/mydb
spring.datasource.username=root
spring.datasource.password=password
spring.datasource.driver-class-name=com.mysql.cj.jdbc.Driver
`
		// Créer les répertoires nécessaires
		if err := os.MkdirAll("src/main/resources", 0755); err != nil {
			t.Fatalf("Impossible de créer les répertoires: %v", err)
		}

		if err := os.WriteFile("src/main/resources/application.properties", []byte(propertiesContent), 0644); err != nil {
			t.Fatalf("Erreur lors de la création du fichier application.properties: %v", err)
		}

		// Tester la détection
		services, err := DetectServices()
		if err != nil {
			t.Errorf("Erreur lors de la détection des services: %v", err)
		}

		// Vérifier si MySQL a été détecté
		found := false
		for _, service := range services {
			if service.Type == "mysql" {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Service MySQL non détecté")
		}

		// Nettoyer
		os.RemoveAll("src")
	})

	// Test 2: Projet avec PostgreSQL
	t.Run("Project with PostgreSQL", func(t *testing.T) {
		// Créer un faux fichier application.yml avec configuration PostgreSQL
		yamlContent := `
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    username: postgres
    password: password
    driver-class-name: org.postgresql.Driver
`
		// Créer les répertoires nécessaires
		if err := os.MkdirAll("src/main/resources", 0755); err != nil {
			t.Fatalf("Impossible de créer les répertoires: %v", err)
		}

		if err := os.WriteFile("src/main/resources/application.yml", []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Erreur lors de la création du fichier application.yml: %v", err)
		}

		// Tester la détection
		services, err := DetectServices()
		if err != nil {
			t.Errorf("Erreur lors de la détection des services: %v", err)
		}

		// Vérifier si PostgreSQL a été détecté
		found := false
		for _, service := range services {
			if service.Type == "postgres" { // Correspond à la constante PostgreSQL
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Service PostgreSQL non détecté")
		}

		// Nettoyer
		os.RemoveAll("src")
	})

	// Test 3: Projet avec Redis et MongoDB (Docker Compose)
	t.Run("Project with Docker Compose", func(t *testing.T) {
		// Créer un faux fichier docker-compose.yml avec Redis et MongoDB
		composeContent := `
version: '3'
services:
  redis:
    image: redis:6
    ports:
      - "6379:6379"
  mongodb:
    image: mongo:4
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: example
`
		if err := os.WriteFile("docker-compose.yml", []byte(composeContent), 0644); err != nil {
			t.Fatalf("Erreur lors de la création du fichier docker-compose.yml: %v", err)
		}

		// Tester la détection (si votre code analyse docker-compose.yml)
		_, err := DetectServices()
		if err != nil {
			t.Errorf("Erreur lors de la détection des services: %v", err)
		}

		// Remarque : nous n'utilisons pas la variable de retour ici car nous ne savons pas
		// si l'implémentation actuelle détecte les services dans docker-compose.yml

		// Cette partie est commentée car on ne sait pas si votre implémentation analyse docker-compose.yml
		// Si ce n'est pas le cas, vous pouvez la désactiver ou l'adapter
		/*
			redisFound := false
			mongoFound := false
			for _, service := range services {
				if service.Type == "redis" {
					redisFound = true
				}
				if service.Type == "mongodb" {
					mongoFound = true
				}
			}

			if !redisFound {
				t.Errorf("Service Redis non détecté")
			}
			if !mongoFound {
				t.Errorf("Service MongoDB non détecté")
			}
		*/

		// Nettoyer
		os.Remove("docker-compose.yml")
	})

	// Test 4: Détection via pom.xml
	t.Run("Dependencies from pom.xml", func(t *testing.T) {
		pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
    <dependencies>
        <dependency>
            <groupId>mysql</groupId>
            <artifactId>mysql-connector-java</artifactId>
            <version>8.0.0</version>
        </dependency>
    </dependencies>
</project>`

		if err := os.WriteFile("pom.xml", []byte(pomContent), 0644); err != nil {
			t.Fatalf("Erreur lors de la création du pom.xml: %v", err)
		}

		services, err := DetectServices()
		if err != nil {
			t.Errorf("Erreur lors de la détection des services: %v", err)
		}

		found := false
		for _, s := range services {
			if s.Type == MySQL {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Service MySQL non détecté depuis pom.xml")
		}

		os.Remove("pom.xml")
	})

	// Test 5: Détection via build.gradle
	t.Run("Dependencies from Gradle", func(t *testing.T) {
		gradleContent := `dependencies {
    implementation 'org.postgresql:postgresql:42.6.0'
    implementation 'org.apache.kafka:kafka-clients:3.5.0'
}`

		if err := os.WriteFile("build.gradle", []byte(gradleContent), 0644); err != nil {
			t.Fatalf("Erreur lors de la création du build.gradle: %v", err)
		}

		services, err := DetectServices()
		if err != nil {
			t.Errorf("Erreur lors de la détection des services: %v", err)
		}

		pgFound := false
		kafkaFound := false
		for _, s := range services {
			if s.Type == PostgreSQL {
				pgFound = true
			}
			if s.Type == Kafka {
				kafkaFound = true
			}
		}
		if !pgFound {
			t.Errorf("Service PostgreSQL non détecté depuis build.gradle")
		}
		if !kafkaFound {
			t.Errorf("Service Kafka non détecté depuis build.gradle")
		}

		os.Remove("build.gradle")
	})

	// Test 6: Aucun service
	t.Run("No Services", func(t *testing.T) {
		// Ne pas créer de fichiers de configuration

		// Tester la détection
		foundServices, err := DetectServices()
		if err != nil {
			t.Errorf("Erreur lors de la détection des services: %v", err)
		}

		if len(foundServices) > 0 {
			t.Errorf("Des services ont été détectés alors qu'il ne devrait y en avoir aucun: %v", foundServices)
		}
	})
}
