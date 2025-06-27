package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"turbotilt/internal/config"
	"turbotilt/internal/logger"

	"github.com/spf13/cobra"
)

// Structure pour stocker les résultats de diagnostics
type diagResult struct {
	installed bool
	version   string
	detail    string
	weight    int // Importance: 3 = critique, 2 = important, 1 = optionnel
	required  bool
}

var (
	verbose          bool
	logToFile        bool
	showAllInfo      bool
	showSummary      bool
	logFilePath      string
	validateManifest bool
	checkTiltfile    bool
	fixMode          bool
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Vérifie installation & config",
	Long: `Vérifie l'installation et la configuration de l'environnement Turbotilt.
Effectue un diagnostic complet des dépendances et outils nécessaires.
Fournit un score de santé et des recommandations pour réparer les problèmes.

Exemples:
  turbotilt doctor             # Analyse complète
  turbotilt doctor --verbose   # Analyse complète avec détails
  turbotilt doctor --debug     # Mode debug pour information détaillée`,
	Run: func(cmd *cobra.Command, args []string) {
		// Enregistrer le temps de démarrage pour calculer la durée d'exécution
		startTime := time.Now()

		// Configurer le logger
		if verbose {
			logger.SetLevel(logger.DEBUG)
			logger.Debug("Mode verbeux activé")
		}

		// Si l'option validate-manifest est activée, valider uniquement le manifeste sans faire le reste du diagnostic
		if validateManifest {
			// Importer le package config
			validateManifestFile()
			return
		}

		// Configurer le chemin du fichier de log si non spécifié
		if logToFile && logFilePath == "" {
			logFilePath = fmt.Sprintf("turbotilt-doctor-%s.log",
				time.Now().Format("20060102-150405"))
		}

		if logToFile {
			if err := logger.EnableFileLogging(logFilePath); err != nil {
				fmt.Printf("⚠️ Impossible de créer le fichier de log: %v\n", err)
			} else {
				defer logger.DisableFileLogging()
				logger.Info("Log file created: %s", logFilePath)
				if verbose {
					fmt.Printf("📄 Logs enregistrés dans: %s\n\n", logFilePath)
				}
			}
		}

		logger.Info("=== Diagnostic de l'environnement Turbotilt ===")
		logger.Info("Démarré le %s", time.Now().Format("2006-01-02 15:04:05"))
		logger.Info("Système d'exploitation: %s", runtime.GOOS)
		fmt.Println("🔍 Vérification de l'environnement Turbotilt...")

		// On a déjà déclaré startTime, pas besoin de redéclaration

		// Structure pour stocker les résultats
		type diagResult struct {
			installed bool
			version   string
			detail    string
			weight    int // Importance: 3 = critique, 2 = important, 1 = optionnel
			required  bool
		}
		results := make(map[string]diagResult)

		fmt.Println("\n📋 Vérification des dépendances requises:")
		logger.Debug("Checking required dependencies...")

		// Vérifier Docker (critique)
		fmt.Print("⏳ Docker : ")
		logger.Debug("Checking Docker installation...")
		if version, err := execCommand("docker", "--version"); err == nil {
			truncVersion := strings.TrimSpace(version)
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Docker installé: %s", truncVersion)

			// Vérification supplémentaire du démon Docker
			if _, err := execCommand("docker", "info"); err != nil {
				fmt.Println("   ⚠️ Le démon Docker ne semble pas être en cours d'exécution")
				logger.Warning("Le démon Docker ne répond pas")
				results["docker"] = diagResult{true, truncVersion, "Le démon ne répond pas", 3, true}
			} else {
				results["docker"] = diagResult{true, truncVersion, "OK", 3, true}
			}
		} else {
			fmt.Println("❌ Non installé ou inaccessible")
			logger.Warning("Docker non installé ou inaccessible")
			results["docker"] = diagResult{false, "", "Non installé", 3, true}
		}

		// Vérifier Docker Compose (critique)
		fmt.Print("⏳ Docker Compose : ")
		logger.Debug("Checking Docker Compose installation...")
		if version, err := execCommand("docker", "compose", "version"); err == nil {
			truncVersion := strings.Split(version, "\n")[0]
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Docker Compose installé: %s", truncVersion)
			results["docker-compose"] = diagResult{true, truncVersion, "OK", 3, true}
		} else {
			fmt.Println("❌ Non installé ou inaccessible")
			logger.Warning("Docker Compose non installé ou inaccessible")
			results["docker-compose"] = diagResult{false, "", "Non installé", 3, true}
		}

		// Vérifier Tilt (important)
		fmt.Print("⏳ Tilt : ")
		logger.Debug("Checking Tilt installation...")
		if version, err := execCommand("tilt", "version"); err == nil {
			truncVersion := strings.Split(version, "\n")[0]
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Tilt installé: %s", truncVersion)
			results["tilt"] = diagResult{true, truncVersion, "OK", 2, false}
		} else {
			fmt.Println("❌ Non installé ou inaccessible")
			fmt.Println("   👉 Pour installer Tilt: https://docs.tilt.dev/install.html")
			logger.Warning("Tilt non installé ou inaccessible")
			results["tilt"] = diagResult{false, "", "Non installé", 2, false}
		}

		fmt.Println("\n📋 Vérification des outils de développement:")
		logger.Debug("Checking development tools...")

		// Vérifier Java (optionnel)
		fmt.Print("⏳ Java : ")
		logger.Debug("Checking Java installation...")
		if version, err := execCommand("java", "-version"); err == nil {
			truncVersion := strings.Split(version, "\n")[0]
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Java installé: %s", truncVersion)
			results["java"] = diagResult{true, truncVersion, "OK", 1, false}
		} else {
			fmt.Println("❌ Non installé ou inaccessible")
			logger.Warning("Java non installé ou inaccessible")
			results["java"] = diagResult{false, "", "Non installé", 1, false}
		}

		// Vérifier Maven (optionnel)
		logger.Debug("Checking Maven installation...")
		fmt.Print("⏳ Maven : ")
		if version, err := execCommand("mvn", "--version"); err == nil {
			truncVersion := strings.Split(version, "\n")[0]
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Maven installé: %s", truncVersion)
			results["maven"] = diagResult{true, truncVersion, "OK", 1, false}
		} else {
			fmt.Println("⚠️ Non installé (requis pour certains projets)")
			logger.Warning("Maven non installé")
			results["maven"] = diagResult{false, "", "Non installé", 1, false}
		}

		logger.Debug("Checking Gradle installation...")
		// Vérifier Gradle
		fmt.Print("⏳ Gradle : ")
		if version, err := execCommand("gradle", "--version"); err == nil {
			truncVersion := strings.Split(strings.Split(version, "\n")[0], "----")[0]
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Gradle installé: %s", truncVersion)
			results["gradle"] = diagResult{true, truncVersion, "OK", 1, false}
		} else {
			fmt.Println("⚠️ Non installé (requis pour certains projets)")
			logger.Warning("Gradle non installé")
			results["gradle"] = diagResult{false, "", "Non installé", 1, false}
		}

		logger.Debug("Checking project files...")
		// Vérifier le projet courant
		fmt.Println("\n📋 Projet courant :")

		// Vérifier si les fichiers requis existent
		projectFiles := checkProjectFiles()

		// Analyser la santé du projet localement
		projectHealth := func(diagnostics map[string]diagResult) int {
			total := 0
			max := 0

			// Additionner les poids des composants installés
			for tool, result := range diagnostics {
				weight := result.weight
				max += weight

				if result.installed {
					if result.detail == "OK" {
						total += weight
					} else if !result.required {
						// Si pas requis, on compte quand même des points partiels
						total += weight / 2
					}
				}

				logger.Debug("Évaluation de %s: installé=%t, poids=%d, détail=%s, score partiel=%d/%d",
					tool, result.installed, weight, result.detail, total, max)
			}

			// Convertir en pourcentage
			if max == 0 {
				return 0
			}
			return total * 100 / max
		}(results)

		// Afficher les recommandations
		fmt.Println("\n📋 Recommandations :")
		if !results["docker"].installed {
			fmt.Println("❗ Docker est requis : https://docs.docker.com/get-docker/")
			logger.Error("Docker manquant - installation requise")
		} else if results["docker"].detail != "OK" {
			fmt.Println("⚠️ Assurez-vous que le daemon Docker est en cours d'exécution")
			logger.Warning("Problème avec Docker: %s", results["docker"].detail)
		}

		if !results["docker-compose"].installed {
			fmt.Println("❗ Docker Compose est requis : https://docs.docker.com/compose/install/")
			logger.Error("Docker Compose manquant - installation requise")
		}

		if !results["tilt"].installed {
			fmt.Println("❗ Tilt est fortement recommandé : https://docs.tilt.dev/install.html")
			logger.Warning("Tilt manquant - installation recommandée")
		}

		// Recommandations spécifiques pour les développeurs Java
		if len(projectFiles["java"]) > 0 && !results["java"].installed {
			fmt.Println("⚠️ Des fichiers Java ont été détectés mais Java n'est pas installé")
			logger.Warning("Java requis pour ce projet mais non installé")
		}

		if len(projectFiles["maven"]) > 0 && !results["maven"].installed {
			fmt.Println("⚠️ Des fichiers Maven ont été détectés mais Maven n'est pas installé")
			logger.Warning("Maven requis pour ce projet mais non installé")
		}

		if len(projectFiles["gradle"]) > 0 && !results["gradle"].installed {
			fmt.Println("⚠️ Des fichiers Gradle ont été détectés mais Gradle n'est pas installé")
			logger.Warning("Gradle requis pour ce projet mais non installé")
		}

		// Afficher la santé globale du projet
		fmt.Println("\n📊 Santé globale :", healthToEmoji(projectHealth))
		logger.Info("Santé globale du projet: %d%%", projectHealth)

		fmt.Println("\n🔧 Commandes disponibles :")
		fmt.Println("▶️ turbotilt init   : Initialiser un projet")
		fmt.Println("▶️ turbotilt up     : Démarrer l'environnement")
		fmt.Println("▶️ turbotilt stop   : Arrêter l'environnement")
		fmt.Println("▶️ turbotilt doctor : Vérifier la configuration")

		// Calculer et afficher le temps d'exécution
		duration := time.Since(startTime)
		fmt.Printf("\n⏱️ Diagnostic complété en %.2f secondes\n", duration.Seconds())

		if logToFile {
			fmt.Printf("📄 Log enregistré dans: %s\n", logFilePath)
		}

		logger.Info("Diagnostic terminé en %.2f secondes", duration.Seconds())
		logger.Debug("Doctor command completed")
	},
}

// Fonction supprimée car remplacée par une implémentation locale

// healthToEmoji convertit un score de santé en une représentation emoji avec barre de progression
func healthToEmoji(health int) string {
	var emoji, grade, barGraph string

	// Déterminer le grade
	switch {
	case health >= 90:
		emoji = "✅"
		grade = "Excellent"
	case health >= 70:
		emoji = "🟢"
		grade = "Bon"
	case health >= 50:
		emoji = "🟡"
		grade = "Moyen"
	case health >= 30:
		emoji = "🟠"
		grade = "Problématique"
	default:
		emoji = "🔴"
		grade = "Critique"
	}

	// Créer une barre de progression visuelle
	completed := health / 10
	remaining := 10 - completed

	barGraph = strings.Repeat("█", completed) + strings.Repeat("░", remaining)

	return fmt.Sprintf("%s %s (%d%%) %s", emoji, grade, health, barGraph)
}

// execCommand exécute une commande et renvoie sa sortie
func execCommand(command string, args ...string) (string, error) {
	logger.Debug("Exécution de la commande: %s %s", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	if err != nil && verbose {
		logger.Debug("Erreur d'exécution: %v", err)
	}

	// Pour des commandes comme 'java -version' qui écrivent sur stderr
	outputStr := string(output)
	if outputStr == "" && err == nil {
		// Certaines commandes peuvent ne pas produire de sortie mais réussir
		outputStr = "OK (pas de sortie)"
	}

	return outputStr, err
}

func checkProjectFiles() map[string][]string {
	// Map pour stocker les fichiers trouvés par catégorie
	foundFiles := map[string][]string{
		"java":   {},
		"maven":  {},
		"gradle": {},
		"docker": {},
		"tilt":   {},
	}

	fileChecks := []struct {
		path     string
		desc     string
		category string
	}{
		{"pom.xml", "Maven", "maven"},
		{"build.gradle", "Gradle", "gradle"},
		{"build.gradle.kts", "Gradle Kotlin", "gradle"},
		{"Dockerfile", "Docker", "docker"},
		{"docker-compose.yml", "Docker Compose", "docker"},
		{"Tiltfile", "Tilt", "tilt"},
		{"src/main/java", "Java Source", "java"},
		{".mvn", "Maven Wrapper", "maven"},
		{"gradlew", "Gradle Wrapper", "gradle"},
	}

	fmt.Println("Vérification des fichiers de projet:")

	for _, check := range fileChecks {
		if _, err := os.Stat(check.path); err == nil {
			fmt.Printf("  ✅ %s trouvé (%s)\n", check.path, check.desc)
			logger.Info("Fichier trouvé: %s (%s)", check.path, check.desc)
			foundFiles[check.category] = append(foundFiles[check.category], check.path)
		} else if verbose {
			fmt.Printf("  ❌ %s non trouvé\n", check.path)
			logger.Debug("Fichier manquant: %s", check.path)
		}
	}

	if len(foundFiles["maven"]) > 0 {
		fmt.Println("  📄 Projet Maven détecté")
	}
	if len(foundFiles["gradle"]) > 0 {
		fmt.Println("  📄 Projet Gradle détecté")
	}
	if len(foundFiles["docker"]) > 0 {
		fmt.Println("  📦 Configuration Docker détectée")
	}
	if len(foundFiles["tilt"]) > 0 {
		fmt.Println("  🚀 Configuration Tilt détectée")
	}

	return foundFiles
}

// validateManifestFile valide le fichier manifeste turbotilt.yaml
func validateManifestFile() {
	fmt.Println("🔍 Validation du manifeste...")

	// Rechercher le manifeste
	configPath, isManifest, err := config.FindConfiguration()
	if err != nil || !isManifest {
		fmt.Println("❌ Manifeste turbotilt.yaml introuvable")
		return
	}

	fmt.Printf("📄 Manifeste trouvé: %s\n", configPath)

	// Charger et valider le manifeste
	manifest, err := config.LoadManifest(configPath)
	if err != nil {
		fmt.Printf("❌ Erreur de validation: %v\n", err)
		return
	}

	// Afficher un résumé du manifeste
	fmt.Println("✅ Le manifeste est valide!")
	fmt.Printf("📊 Contient %d service(s):\n", len(manifest.Services))

	// Afficher les détails des services
	for i, service := range manifest.Services {
		fmt.Printf("   [%d] %s\n", i+1, service.Name)

		if service.Runtime != "" {
			fmt.Printf("       - Type: Application (%s)\n", service.Runtime)
			fmt.Printf("       - Path: %s\n", service.Path)
			fmt.Printf("       - Port: %s\n", service.Port)
			fmt.Printf("       - Java: %s\n", service.Java)
		} else if service.Type != "" {
			fmt.Printf("       - Type: Service dépendant (%s)\n", service.Type)
			if service.Version != "" {
				fmt.Printf("       - Version: %s\n", service.Version)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(doctorCmd)

	doctorCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Afficher des informations détaillées")
	doctorCmd.Flags().BoolVarP(&logToFile, "log", "l", false, "Enregistrer les résultats dans un fichier log")
	doctorCmd.Flags().StringVar(&logFilePath, "log-file", "", "Chemin du fichier log à utiliser")
	doctorCmd.Flags().BoolVar(&showAllInfo, "all", false, "Afficher toutes les informations")
	doctorCmd.Flags().BoolVar(&showSummary, "summary", false, "Afficher uniquement le résumé")
	doctorCmd.Flags().BoolVar(&validateManifest, "validate-manifest", false, "Valider la syntaxe du manifeste turbotilt.yaml")
}
