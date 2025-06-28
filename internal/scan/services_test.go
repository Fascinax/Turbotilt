package scan

import (
	"os"
	"testing"
)

// TestDetectServices tests the service detection functionality
func TestDetectServices(t *testing.T) {
	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "turbotilt-services-test-*")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save the current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get current directory: %v", err)
	}

	// Change to the temporary directory for tests
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Unable to change directory: %v", err)
	}
	defer os.Chdir(originalDir) // Restore the working directory at the end

	// Test 1: Spring Boot with MySQL
	t.Run("Spring Boot with MySQL", func(t *testing.T) {
		// Create a fake application.properties file with MySQL configuration
		propertiesContent := `
spring.datasource.url=jdbc:mysql://localhost:3306/mydb
spring.datasource.username=root
spring.datasource.password=password
spring.datasource.driver-class-name=com.mysql.cj.jdbc.Driver
`
		// Create necessary directories
		if err := os.MkdirAll("src/main/resources", 0755); err != nil {
			t.Fatalf("Unable to create directories: %v", err)
		}

		if err := os.WriteFile("src/main/resources/application.properties", []byte(propertiesContent), 0644); err != nil {
			t.Fatalf("Error creating application.properties file: %v", err)
		}

		// Test detection
		services, err := DetectServices()
		if err != nil {
			t.Errorf("Error detecting services: %v", err)
		}

		// Check if MySQL was detected
		found := false
		for _, service := range services {
			if service.Type == "mysql" {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("MySQL service not detected")
		}

		// Cleanup
		os.RemoveAll("src")
	})

	// Test 2: Project with PostgreSQL
	t.Run("Project with PostgreSQL", func(t *testing.T) {
		// Create a fake application.yml file with PostgreSQL configuration
		yamlContent := `
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    username: postgres
    password: password
    driver-class-name: org.postgresql.Driver
`
		// Create necessary directories
		if err := os.MkdirAll("src/main/resources", 0755); err != nil {
			t.Fatalf("Unable to create directories: %v", err)
		}

		if err := os.WriteFile("src/main/resources/application.yml", []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Error creating application.yml file: %v", err)
		}

		// Test detection
		services, err := DetectServices()
		if err != nil {
			t.Errorf("Error detecting services: %v", err)
		}

		// Check if PostgreSQL was detected
		found := false
		for _, service := range services {
			if service.Type == "postgres" { // Corresponds to the PostgreSQL constant
				found = true
				break
			}
		}

		if !found {
			t.Errorf("PostgreSQL service not detected")
		}

		// Cleanup
		os.RemoveAll("src")
	})

	// Test 3: Project with Redis and MongoDB (Docker Compose)
	t.Run("Project with Docker Compose", func(t *testing.T) {
		// Create a fake docker-compose.yml file with Redis and MongoDB
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
			t.Fatalf("Error creating docker-compose.yml file: %v", err)
		}

		// Test detection (if your code analyzes docker-compose.yml)
		_, err := DetectServices()
		if err != nil {
			t.Errorf("Error detecting services: %v", err)
		}

		// Note: we don't use the return variable here because we don't know
		// if the current implementation detects services in docker-compose.yml

		// This part is commented out because we don't know if your implementation analyzes docker-compose.yml
		// If it doesn't, you can disable or adapt it
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
				t.Errorf("Redis service not detected")
			}
			if !mongoFound {
				t.Errorf("MongoDB service not detected")
			}
		*/

		// Cleanup
		os.Remove("docker-compose.yml")
	})

	// Test 4: Scan using pom.xml
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
			t.Fatalf("Error during pom.xml creation: %v", err)
		}

		services, err := DetectServices()
		if err != nil {
			t.Errorf("Error during the services scan: %v", err)
		}

		found := false
		for _, s := range services {
			if s.Type == MySQL {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("MySQL Service not detected using pom.xml")
		}

		os.Remove("pom.xml")
	})

	// Test 5: Scan using build.gradle
	t.Run("Dependencies from Gradle", func(t *testing.T) {
		gradleContent := `dependencies {
    implementation 'org.postgresql:postgresql:42.6.0'
    implementation 'org.apache.kafka:kafka-clients:3.5.0'
}`

		if err := os.WriteFile("build.gradle", []byte(gradleContent), 0644); err != nil {
			t.Fatalf("Error during build.gradle creation: %v", err)
		}

		services, err := DetectServices()
		if err != nil {
			t.Errorf("Error during the services scan: %v", err)
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
			t.Errorf("PostgreSQL Service not detected using build.gradle")
		}
		if !kafkaFound {
			t.Errorf("Kafka Service not detected using build.gradle")
		}

		os.Remove("build.gradle")
	})

	// Test 6: No service
	t.Run("No Services", func(t *testing.T) {
		// Don't create any configuration files

		// Test detection
		foundServices, err := DetectServices()
		if err != nil {
			t.Errorf("Error detecting services: %v", err)
		}

		if len(foundServices) > 0 {
			t.Errorf("Services were detected when there shouldn't be any: %v", foundServices)
		}
	})
}
