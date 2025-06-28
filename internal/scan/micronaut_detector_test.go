package scan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMicronautDetector(t *testing.T) {
	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "micronaut-detector-test")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Directory structure for a Micronaut project
	srcMainJava := filepath.Join(tempDir, "src", "main", "java", "com", "example")
	srcMainResources := filepath.Join(tempDir, "src", "main", "resources")

	// Create the directories
	if err := os.MkdirAll(srcMainJava, 0755); err != nil {
		t.Fatalf("Unable to create directories: %v", err)
	}
	if err := os.MkdirAll(srcMainResources, 0755); err != nil {
		t.Fatalf("Unable to create directories: %v", err)
	}

	t.Run("Test detection via pom.xml", func(t *testing.T) {
		// Create a dummy pom.xml with Micronaut dependencies
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
			t.Fatalf("Unable to write pom.xml file: %v", err)
		}

		detector := MicronautDetector{}
		detected, result, err := detector.Detect(tempDir)

		if err != nil {
			t.Errorf("Error during detection: %v", err)
		}

		if !detected {
			t.Errorf("Detector did not detect Micronaut via pom.xml")
		}

		if result.Framework != "micronaut" {
			t.Errorf("Incorrectly detected framework: %s", result.Framework)
		}

		if result.BuildSystem != "maven" {
			t.Errorf("Incorrectly detected build system: %s", result.BuildSystem)
		}

		// Remove pom.xml for the next test
		os.Remove(filepath.Join(tempDir, "pom.xml"))
	})

	t.Run("Test detection via build.gradle", func(t *testing.T) {
		// Create a dummy build.gradle with Micronaut dependencies
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
			t.Fatalf("Unable to write build.gradle file: %v", err)
		}

		detector := MicronautDetector{}
		detected, result, err := detector.Detect(tempDir)

		if err != nil {
			t.Errorf("Error during detection: %v", err)
		}

		if !detected {
			t.Errorf("Detector did not detect Micronaut via build.gradle")
		}

		if result.Framework != "micronaut" {
			t.Errorf("Incorrectly detected framework: %s", result.Framework)
		}

		if result.BuildSystem != "gradle" {
			t.Errorf("Incorrectly detected build system: %s", result.BuildSystem)
		}

		// Remove build.gradle for the next test
		os.Remove(filepath.Join(tempDir, "build.gradle"))
	})

	t.Run("Test detection via application.yml", func(t *testing.T) {
		// Create a dummy application.yml
		appYaml := `micronaut:
  application:
    name: test-app
  server:
    port: 8080
`

		if err := os.WriteFile(filepath.Join(srcMainResources, "application.yml"), []byte(appYaml), 0644); err != nil {
			t.Fatalf("Unable to write application.yml file: %v", err)
		}

		// Add a build.gradle for build system detection
		if err := os.WriteFile(filepath.Join(tempDir, "build.gradle"), []byte("// Empty build file"), 0644); err != nil {
			t.Fatalf("Unable to write build.gradle file: %v", err)
		}

		detector := MicronautDetector{}
		detected, result, err := detector.Detect(tempDir)

		if err != nil {
			t.Errorf("Error during detection: %v", err)
		}

		if !detected {
			t.Errorf("Detector did not detect Micronaut via application.yml")
		}

		if result.Framework != "micronaut" {
			t.Errorf("Incorrectly detected framework: %s", result.Framework)
		}

		if result.BuildSystem != "gradle" {
			t.Errorf("Incorrectly detected build system: %s", result.BuildSystem)
		}

		// Cleanup
		os.Remove(filepath.Join(srcMainResources, "application.yml"))
		os.Remove(filepath.Join(tempDir, "build.gradle"))
	})

	t.Run("Test detection via imports in Java files", func(t *testing.T) {
		// Create a Java controller with Micronaut imports
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
			t.Fatalf("Unable to write Java file: %v", err)
		}

		// Add a pom.xml for build system detection
		if err := os.WriteFile(filepath.Join(tempDir, "pom.xml"), []byte("<project></project>"), 0644); err != nil {
			t.Fatalf("Unable to write pom.xml file: %v", err)
		}

		detector := MicronautDetector{}
		detected, result, err := detector.Detect(tempDir)

		if err != nil {
			t.Errorf("Error during detection: %v", err)
		}

		if !detected {
			t.Errorf("Detector did not detect Micronaut via Java imports")
		}

		if result.Framework != "micronaut" {
			t.Errorf("Incorrectly detected framework: %s", result.Framework)
		}

		if result.BuildSystem != "maven" {
			t.Errorf("Incorrectly detected build system: %s", result.BuildSystem)
		}
	})
}
