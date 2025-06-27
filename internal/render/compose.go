package render

import (
	"fmt"
	"os"
	"strings"
	"turbotilt/internal/scan"
)

// ComposeServiceDefinition définit un service dans docker-compose.yml
type ComposeServiceDefinition struct {
	Name        string
	Image       string
	Port        string
	Environment map[string]string
	Volumes     []string
	DependsOn   []string
}

// GenerateComposeWithServices génère un docker-compose.yml incluant les services détectés
func GenerateComposeWithServices(opts Options) error {
	// Utiliser les services de l'objet opts
	services := opts.Services
	f, err := os.Create("docker-compose.yml")
	if err != nil {
		return fmt.Errorf("erreur lors de la création du docker-compose.yml: %w", err)
	}
	defer f.Close()

	// Service principal de l'application
	appService := ComposeServiceDefinition{
		Name:  "app",
		Image: "build: .",
		Port:  fmt.Sprintf("%s:%s", opts.Port, opts.Port),
		Volumes: []string{
			"./src:/app/src",
		},
		Environment: make(map[string]string),
	}

	// Configurer l'environnement selon le framework
	switch opts.Framework {
	case "spring":
		if opts.DevMode {
			appService.Environment["SPRING_PROFILES_ACTIVE"] = "dev"
		} else {
			appService.Environment["SPRING_PROFILES_ACTIVE"] = "prod"
		}
	case "quarkus":
		if opts.DevMode {
			appService.Environment["QUARKUS_PROFILE"] = "dev"
		} else {
			appService.Environment["QUARKUS_PROFILE"] = "prod"
		}
	case "micronaut":
		if opts.DevMode {
			appService.Environment["MICRONAUT_ENVIRONMENTS"] = "dev"
		} else {
			appService.Environment["MICRONAUT_ENVIRONMENTS"] = "prod"
		}
	}

	// Liste des services à inclure dans docker-compose.yml
	serviceDefinitions := []ComposeServiceDefinition{appService}
	volumes := make(map[string]bool)

	// Ajouter les services détectés
	for _, service := range services {
		var serviceDefinition ComposeServiceDefinition

		switch service.Type {
		case scan.MySQL:
			serviceDefinition = ComposeServiceDefinition{
				Name:    "mysql",
				Image:   fmt.Sprintf("image: mysql:%s", getOrDefault(service.Version, "latest")),
				Port:    fmt.Sprintf("%s:3306", getOrDefault(service.Port, "3306")),
				Volumes: []string{"mysql_data:/var/lib/mysql"},
				Environment: map[string]string{
					"MYSQL_ROOT_PASSWORD": getFromCredentials(service.Credentials, "MYSQL_ROOT_PASSWORD", "root"),
					"MYSQL_DATABASE":      getFromCredentials(service.Credentials, "MYSQL_DATABASE", "app"),
				},
			}
			volumes["mysql_data"] = true

		case scan.PostgreSQL:
			serviceDefinition = ComposeServiceDefinition{
				Name:    "postgres",
				Image:   fmt.Sprintf("image: postgres:%s", getOrDefault(service.Version, "latest")),
				Port:    fmt.Sprintf("%s:5432", getOrDefault(service.Port, "5432")),
				Volumes: []string{"postgres_data:/var/lib/postgresql/data"},
				Environment: map[string]string{
					"POSTGRES_USER":     getFromCredentials(service.Credentials, "POSTGRES_USER", "postgres"),
					"POSTGRES_PASSWORD": getFromCredentials(service.Credentials, "POSTGRES_PASSWORD", "postgres"),
					"POSTGRES_DB":       getFromCredentials(service.Credentials, "POSTGRES_DB", "app"),
				},
			}
			volumes["postgres_data"] = true

		case scan.MongoDB:
			serviceDefinition = ComposeServiceDefinition{
				Name:    "mongodb",
				Image:   fmt.Sprintf("image: mongo:%s", getOrDefault(service.Version, "latest")),
				Port:    fmt.Sprintf("%s:27017", getOrDefault(service.Port, "27017")),
				Volumes: []string{"mongo_data:/data/db"},
			}
			volumes["mongo_data"] = true

		case scan.Redis:
			serviceDefinition = ComposeServiceDefinition{
				Name:    "redis",
				Image:   fmt.Sprintf("image: redis:%s", getOrDefault(service.Version, "latest")),
				Port:    fmt.Sprintf("%s:6379", getOrDefault(service.Port, "6379")),
				Volumes: []string{"redis_data:/data"},
			}
			volumes["redis_data"] = true

		case scan.Kafka:
			serviceDefinition = ComposeServiceDefinition{
				Name:  "kafka",
				Image: "image: confluentinc/cp-kafka:latest",
				Port:  "9092:9092",
				Environment: map[string]string{
					"KAFKA_ADVERTISED_LISTENERS":             "PLAINTEXT://kafka:9092",
					"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":   "PLAINTEXT:PLAINTEXT",
					"KAFKA_INTER_BROKER_LISTENER_NAME":       "PLAINTEXT",
					"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
				},
				DependsOn: []string{"zookeeper"},
			}

			// Ajouter aussi Zookeeper
			zookeeper := ComposeServiceDefinition{
				Name:  "zookeeper",
				Image: "image: confluentinc/cp-zookeeper:latest",
				Port:  "2181:2181",
				Environment: map[string]string{
					"ZOOKEEPER_CLIENT_PORT": "2181",
				},
			}
			serviceDefinitions = append(serviceDefinitions, zookeeper)

		case scan.RabbitMQ:
			serviceDefinition = ComposeServiceDefinition{
				Name:    "rabbitmq",
				Image:   "image: rabbitmq:3-management",
				Port:    "5672:5672",
				Volumes: []string{"rabbitmq_data:/var/lib/rabbitmq"},
				Environment: map[string]string{
					"RABBITMQ_DEFAULT_USER": getFromCredentials(service.Credentials, "RABBITMQ_DEFAULT_USER", "guest"),
					"RABBITMQ_DEFAULT_PASS": getFromCredentials(service.Credentials, "RABBITMQ_DEFAULT_PASS", "guest"),
				},
			}
			volumes["rabbitmq_data"] = true

		case scan.ElasticSearch:
			serviceDefinition = ComposeServiceDefinition{
				Name:    "elasticsearch",
				Image:   "image: docker.elastic.co/elasticsearch/elasticsearch:7.14.0",
				Port:    "9200:9200",
				Volumes: []string{"es_data:/usr/share/elasticsearch/data"},
				Environment: map[string]string{
					"discovery.type": "single-node",
					"ES_JAVA_OPTS":   "-Xms512m -Xmx512m",
				},
			}
			volumes["es_data"] = true
		}

		if serviceDefinition.Name != "" {
			appService.DependsOn = append(appService.DependsOn, serviceDefinition.Name)
			serviceDefinitions = append(serviceDefinitions, serviceDefinition)
		}
	}

	// Mise à jour de la définition de l'application pour ajouter les dépendances
	serviceDefinitions[0] = appService

	// Construire le contenu du fichier docker-compose.yml
	var sb strings.Builder
	sb.WriteString("version: '3'\n\nservices:\n")

	// Ajouter tous les services
	for _, service := range serviceDefinitions {
		sb.WriteString(fmt.Sprintf("  %s:\n", service.Name))
		sb.WriteString(fmt.Sprintf("    %s\n", service.Image))

		if service.Port != "" {
			sb.WriteString("    ports:\n")
			sb.WriteString(fmt.Sprintf("      - '%s'\n", service.Port))
		}

		if len(service.Environment) > 0 {
			sb.WriteString("    environment:\n")
			for k, v := range service.Environment {
				sb.WriteString(fmt.Sprintf("      - %s=%s\n", k, v))
			}
		}

		if len(service.Volumes) > 0 {
			sb.WriteString("    volumes:\n")
			for _, volume := range service.Volumes {
				sb.WriteString(fmt.Sprintf("      - %s\n", volume))
			}
		}

		if len(service.DependsOn) > 0 {
			sb.WriteString("    depends_on:\n")
			for _, dep := range service.DependsOn {
				sb.WriteString(fmt.Sprintf("      - %s\n", dep))
			}
		}

		sb.WriteString("\n")
	}

	// Ajouter les volumes si nécessaire
	if len(volumes) > 0 {
		sb.WriteString("volumes:\n")
		for volume := range volumes {
			sb.WriteString(fmt.Sprintf("  %s:\n", volume))
		}
	}

	// Écrire le contenu dans le fichier
	_, err = f.WriteString(sb.String())
	return err
}

// getOrDefault retourne la valeur ou une valeur par défaut si vide
func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// getFromCredentials récupère une valeur des credentials ou retourne la valeur par défaut
func getFromCredentials(credentials map[string]string, key, defaultValue string) string {
	if value, ok := credentials[key]; ok && value != "" {
		return value
	}
	return defaultValue
}
