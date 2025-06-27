package scan

// BaseDetector contient les fonctionnalités communes à tous les détecteurs
type BaseDetector struct {
	// Champs communs pour tous les détecteurs
}

// DetectionResult contient les résultats de la détection d'un framework
type DetectionResult struct {
	Framework    string            // Le nom du framework détecté (spring, quarkus, micronaut, etc.)
	BuildSystem  string            // Le système de build (maven, gradle, etc.)
	JavaVersion  string            // La version de Java
	Port         string            // Le port de l'application
	Dependencies []string          // Les dépendances du projet
	Properties   map[string]string // Les propriétés de configuration
	Detected     bool              // Indique si le framework a été détecté
}
