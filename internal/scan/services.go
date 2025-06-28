package scan

import (
	"os"
	"strings"
)

// ServiceType represents a detected service type
type ServiceType string

const (
	MySQL         ServiceType = "mysql"
	PostgreSQL    ServiceType = "postgres"
	MongoDB       ServiceType = "mongodb"
	Redis         ServiceType = "redis"
	Kafka         ServiceType = "kafka"
	RabbitMQ      ServiceType = "rabbitmq"
	ElasticSearch ServiceType = "elasticsearch"
)

// ServiceConfig contains the configuration of a detected service
type ServiceConfig struct {
	Type        ServiceType
	Version     string
	Port        string
	Credentials map[string]string
}

// DetectServices detects the services required for the project
func DetectServices() ([]ServiceConfig, error) {
	var services []ServiceConfig

	// Detect from configuration files
	configServices, err := detectFromConfigFiles()
	if err != nil {
		return nil, err
	}
	services = append(services, configServices...)

	// Detect from dependencies (Maven/Gradle)
	depServices, err := detectFromDependencies()
	if err != nil {
		return nil, err
	}
	services = append(services, depServices...)

	return services, nil
}

// detectFromConfigFiles detects services from configuration files
func detectFromConfigFiles() ([]ServiceConfig, error) {
	var services []ServiceConfig

	// Check in application.properties
	propsServices, err := detectFromPropertiesFile("src/main/resources/application.properties")
	if err == nil {
		services = append(services, propsServices...)
	}

	// Check in application.yml
	yamlServices, err := detectFromYamlFile("src/main/resources/application.yml")
	if err == nil {
		services = append(services, yamlServices...)
	}

	return services, nil
}

// detectFromPropertiesFile detects services from a properties file
func detectFromPropertiesFile(path string) ([]ServiceConfig, error) {
	var services []ServiceConfig

	// Look for patterns like spring.datasource.url, etc.
	// Examples:
	// spring.datasource.url=jdbc:mysql://localhost:3306/db
	// spring.data.mongodb.uri=mongodb://localhost:27017/mydb

	// For this version, we implement basic detection

	// MySQL
	if hasPropertyPattern(path, "jdbc:mysql") {
		services = append(services, ServiceConfig{
			Type:    MySQL,
			Version: "latest",
			Port:    "3306",
			Credentials: map[string]string{
				"MYSQL_ROOT_PASSWORD": "root",
				"MYSQL_DATABASE":      "app",
			},
		})
	}

	// PostgreSQL
	if hasPropertyPattern(path, "jdbc:postgresql") {
		services = append(services, ServiceConfig{
			Type:    PostgreSQL,
			Version: "latest",
			Port:    "5432",
			Credentials: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "app",
			},
		})
	}

	// MongoDB
	if hasPropertyPattern(path, "mongodb://") {
		services = append(services, ServiceConfig{
			Type:    MongoDB,
			Version: "latest",
			Port:    "27017",
		})
	}

	// Redis
	if hasPropertyPattern(path, "redis://") || hasPropertyPattern(path, "spring.redis") {
		services = append(services, ServiceConfig{
			Type:    Redis,
			Version: "latest",
			Port:    "6379",
		})
	}

	return services, nil
}

// detectFromYamlFile detects services from a YAML file
func detectFromYamlFile(path string) ([]ServiceConfig, error) {
	var services []ServiceConfig

	// Simplification: we'll just search for strings in the YAML file
	// as a first implementation, even if it's not robust
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	// PostgreSQL detection
	if strings.Contains(contentStr, "jdbc:postgresql") {
		services = append(services, ServiceConfig{
			Type:    PostgreSQL,
			Version: "latest",
			Port:    "5432",
			Credentials: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "app",
			},
		})
	}

	// MySQL detection
	if strings.Contains(contentStr, "jdbc:mysql") {
		services = append(services, ServiceConfig{
			Type:    MySQL,
			Version: "latest",
			Port:    "3306",
			Credentials: map[string]string{
				"MYSQL_ROOT_PASSWORD": "root",
				"MYSQL_DATABASE":      "app",
			},
		})
	}

	// MongoDB detection
	if strings.Contains(contentStr, "mongodb://") {
		services = append(services, ServiceConfig{
			Type:    MongoDB,
			Version: "latest",
			Port:    "27017",
		})
	}

	// Redis detection
	if strings.Contains(contentStr, "redis://") || strings.Contains(contentStr, "spring.redis") {
		services = append(services, ServiceConfig{
			Type:    Redis,
			Version: "latest",
			Port:    "6379",
		})
	}

	return services, nil
}

// detectFromDependencies detects services from project dependencies
func detectFromDependencies() ([]ServiceConfig, error) {
	var services []ServiceConfig

	// Liste of build files to check
	buildFiles := []string{"pom.xml", "build.gradle", "build.gradle.kts"}

	// Utils function to check if a service is already in the list
	contains := func(t ServiceType) bool {
		for _, s := range services {
			if s.Type == t {
				return true
			}
		}
		return false
	}

	// Check each build file for service patterns
	for _, file := range buildFiles {
		if _, err := os.Stat(file); err != nil {
			continue
		}

		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		content := strings.ToLower(string(data))

		if strings.Contains(content, "mysql") && !contains(MySQL) {
			services = append(services, ServiceConfig{
				Type:    MySQL,
				Version: "latest",
				Port:    "3306",
				Credentials: map[string]string{
					"MYSQL_ROOT_PASSWORD": "root",
					"MYSQL_DATABASE":      "app",
				},
			})
		}

		if (strings.Contains(content, "postgresql") || strings.Contains(content, "postgres")) && !contains(PostgreSQL) {
			services = append(services, ServiceConfig{
				Type:    PostgreSQL,
				Version: "latest",
				Port:    "5432",
				Credentials: map[string]string{
					"POSTGRES_USER":     "postgres",
					"POSTGRES_PASSWORD": "postgres",
					"POSTGRES_DB":       "app",
				},
			})
		}

		if strings.Contains(content, "mongodb") && !contains(MongoDB) {
			services = append(services, ServiceConfig{
				Type:    MongoDB,
				Version: "latest",
				Port:    "27017",
			})
		}

		if strings.Contains(content, "redis") && !contains(Redis) {
			services = append(services, ServiceConfig{
				Type:    Redis,
				Version: "latest",
				Port:    "6379",
			})
		}

		if strings.Contains(content, "kafka") && !contains(Kafka) {
			services = append(services, ServiceConfig{
				Type:    Kafka,
				Version: "latest",
				Port:    "9092",
			})
		}

		if strings.Contains(content, "rabbitmq") && !contains(RabbitMQ) {
			services = append(services, ServiceConfig{
				Type:    RabbitMQ,
				Version: "latest",
				Port:    "5672",
			})
		}

		if strings.Contains(content, "elasticsearch") && !contains(ElasticSearch) {
			services = append(services, ServiceConfig{
				Type:    ElasticSearch,
				Version: "latest",
				Port:    "9200",
			})
		}
	}

	return services, nil
}

// hasPropertyPattern checks if a file contains a certain pattern
func hasPropertyPattern(path string, pattern string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), pattern)
}
