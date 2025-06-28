package scan

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDetectFramework tests the framework detection functionality
func TestDetectFramework(t *testing.T) {
	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "turbotilt-test-*")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get current directory: %v", err)
	}

	// Change to temporary directory for tests
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Unable to change directory: %v", err)
	}
	defer os.Chdir(originalDir) // Restore working directory at the end

	// Test 1: Maven Spring Boot Project
	t.Run("Spring Boot Maven Project", func(t *testing.T) {
		// Create a mock pom.xml for Spring Boot
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
			t.Fatalf("Error creating pom.xml file: %v", err)
		}

		// Test detection
		framework, err := DetectFramework()
		if err != nil {
			t.Errorf("Error detecting framework: %v", err)
		}
		if framework != "spring" {
			t.Errorf("Incorrect detected framework. Expected: %s, Got: %s", "spring", framework)
		}

		// Cleanup
		os.Remove("pom.xml")
	})

	// Test 2: Quarkus Maven Project
	t.Run("Quarkus Maven Project", func(t *testing.T) {
		// Create a mock pom.xml for Quarkus
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
			t.Fatalf("Error creating pom.xml file: %v", err)
		}

		// Test detection
		framework, err := DetectFramework()
		if err != nil {
			t.Errorf("Error detecting framework: %v", err)
		}
		if framework != "quarkus" {
			t.Errorf("Incorrect detected framework. Expected: %s, Got: %s", "quarkus", framework)
		}

		// Cleanup
		os.Remove("pom.xml")
	})

	// Test 3: Gradle Project
	t.Run("Gradle Project", func(t *testing.T) {
		// Create a mock build.gradle for Spring Boot
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
			t.Fatalf("Error creating build.gradle file: %v", err)
		}

		// Test detection
		framework, err := DetectFramework()
		if err != nil {
			t.Errorf("Error detecting framework: %v", err)
		}
		// Based on your current implementation for Gradle
		// We would expect the framework to be spring or java
		expectedFrameworks := map[string]bool{
			"spring": true,
			"java":   true,
		}
		if !expectedFrameworks[framework] {
			t.Errorf("Incorrect detected framework. Expected: spring or java, Got: %s", framework)
		}

		// Cleanup
		os.Remove("build.gradle")
	})

	// Test 4: Unknown Project
	t.Run("Unknown Project", func(t *testing.T) {
		// Don't create any project files

		// Test detection
		framework, err := DetectFramework()
		if err != nil {
			t.Errorf("Error detecting framework: %v", err)
		}
		if framework != "unknown" {
			t.Errorf("Incorrect detected framework. Expected: %s, Got: %s", "unknown", framework)
		}
	})
}

// TestDetectors tests that all detectors are correctly called
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
			// Create a temporary directory
			tempDir, err := os.MkdirTemp("", "detector-test")
			if err != nil {
				t.Fatalf("Error creating temporary directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Create mock files
			for path, content := range tc.mockFiles {
				fullPath := filepath.Join(tempDir, path)
				dirPath := filepath.Dir(fullPath)

				if err := os.MkdirAll(dirPath, 0755); err != nil {
					t.Fatalf("Erreur lors de la création du répertoire %s: %v", dirPath, err)
				}

				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("Error creating file %s: %v", fullPath, err)
				}
			}

			// Run the scanner
			scanner := NewScanner(tempDir)
			framework, result, err := scanner.DetectFramework()

			if err != nil {
				t.Errorf("Error during detection: %v", err)
			}

			if framework != tc.expectedType {
				t.Errorf("Incorrect detected framework: %s, expected: %s", framework, tc.expectedType)
			}

			if result.Framework != tc.expectedType {
				t.Errorf("Framework dans le résultat incorrect: %s, attendu: %s", result.Framework, tc.expectedType)
			}
		})
	}
}
