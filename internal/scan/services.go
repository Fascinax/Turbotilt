package scan

import (
	"os"
	"strings"
)

// ServiceType représente un type de service détecté
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

// ServiceConfig contient la configuration d'un service détecté
type ServiceConfig struct {
	Type        ServiceType
	Version     string
	Port        string
	Credentials map[string]string
}

// DetectServices détecte les services nécessaires pour le projet
func DetectServices() ([]ServiceConfig, error) {
	var services []ServiceConfig

	// Détecter à partir des fichiers de configuration
	configServices, err := detectFromConfigFiles()
	if err != nil {
		return nil, err
	}
	services = append(services, configServices...)

	// Détecter à partir des dépendances (Maven/Gradle)
	depServices, err := detectFromDependencies()
	if err != nil {
		return nil, err
	}
	services = append(services, depServices...)

	return services, nil
}

// detectFromConfigFiles détecte les services à partir des fichiers de configuration
func detectFromConfigFiles() ([]ServiceConfig, error) {
	var services []ServiceConfig

	// Vérifier dans application.properties
	propsServices, err := detectFromPropertiesFile("src/main/resources/application.properties")
	if err == nil {
		services = append(services, propsServices...)
	}

	// Vérifier dans application.yml
	yamlServices, err := detectFromYamlFile("src/main/resources/application.yml")
	if err == nil {
		services = append(services, yamlServices...)
	}

	return services, nil
}

// detectFromPropertiesFile détecte les services à partir d'un fichier properties
func detectFromPropertiesFile(path string) ([]ServiceConfig, error) {
	var services []ServiceConfig

	// Rechercher des patterns comme spring.datasource.url, etc.
	// Exemples:
	// spring.datasource.url=jdbc:mysql://localhost:3306/db
	// spring.data.mongodb.uri=mongodb://localhost:27017/mydb

	// Pour cette version, on implémente une détection basique

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

// detectFromYamlFile détecte les services à partir d'un fichier YAML
func detectFromYamlFile(path string) ([]ServiceConfig, error) {
	var services []ServiceConfig

	// Même logique que pour les fichiers properties mais avec une syntaxe YAML
	// Par exemple:
	// spring:
	//   datasource:
	//     url: jdbc:mysql://localhost:3306/db

	// Pour l'instant, renvoyer une liste vide
	return services, nil
}

// detectFromDependencies détecte les services à partir des dépendances du projet
func detectFromDependencies() ([]ServiceConfig, error) {
	var services []ServiceConfig

	// Analyser pom.xml ou build.gradle pour détecter les dépendances
	// qui indiquent l'utilisation de certains services

	// Par exemple, la présence de mysql-connector-java suggère MySQL
	// Spring Data JPA suggère une base de données SQL, etc.

	// Pour l'instant, renvoyer une liste vide
	return services, nil
}

// hasPropertyPattern vérifie si un fichier contient un certain pattern
func hasPropertyPattern(path string, pattern string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), pattern)
}
