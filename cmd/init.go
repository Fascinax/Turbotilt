package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"turbotilt/internal/render"
	"turbotilt/internal/scan"

	"github.com/spf13/cobra"
)

var (
	// Options pour la commande init
	forceFramework string
	port           string
	jdkVersion     string
	devMode        bool
	detectServices bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scan et génère Tiltfile & Compose",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Initialisation de Turbotilt...")

		// Détecter le framework ou utiliser celui spécifié
		framework := forceFramework
		var err error

		if framework == "" {
			framework, err = scan.DetectFramework()
			if err != nil {
				fmt.Printf("❌ Erreur lors de la détection du framework: %v\n", err)
				return
			}
		}

		fmt.Printf("✅ Framework détecté/sélectionné: %s\n", framework)

		// Détecter les services si demandé
		var services []scan.ServiceConfig
		if detectServices {
			fmt.Println("🔍 Détection des services dépendants...")
			services, err = scan.DetectServices()
			if err != nil {
				fmt.Printf("⚠️ Avertissement lors de la détection des services: %v\n", err)
			}

			// Afficher les services détectés
			if len(services) > 0 {
				fmt.Println("✅ Services détectés:")
				for _, service := range services {
					fmt.Printf("   - %s\n", service.Type)
				}
			} else {
				fmt.Println("ℹ️ Aucun service dépendant détecté")
			}
		}

		// Déterminer le nom de l'application (dossier courant par défaut)
		appName := "app"
		cwd, err := os.Getwd()
		if err == nil {
			appName = filepath.Base(cwd)
		}

		// Préparer les options de rendu
		renderOpts := render.Options{
			Framework:  framework,
			AppName:    appName,
			Port:       port,
			JDKVersion: jdkVersion,
			DevMode:    devMode,
			Services:   services,
		}

		// Générer les fichiers
		if err := render.GenerateDockerfile(renderOpts); err != nil {
			fmt.Printf("❌ Erreur lors de la génération du Dockerfile: %v\n", err)
			return
		}

		// Utiliser le nouveau générateur de docker-compose avec support des services
		if len(services) > 0 {
			if err := render.GenerateComposeWithServices(renderOpts); err != nil {
				fmt.Printf("❌ Erreur lors de la génération du docker-compose.yml: %v\n", err)
				return
			}
		} else {
			if err := render.GenerateCompose(renderOpts); err != nil {
				fmt.Printf("❌ Erreur lors de la génération du docker-compose.yml: %v\n", err)
				return
			}
		}

		if err := render.GenerateTiltfile(renderOpts); err != nil {
			fmt.Printf("❌ Erreur lors de la génération du Tiltfile: %v\n", err)
			return
		}

		fmt.Println("✨ Configuration Turbotilt terminée!")
		fmt.Println("📋 Fichiers générés:")
		fmt.Println("   - Dockerfile")
		fmt.Println("   - docker-compose.yml")
		fmt.Println("   - Tiltfile")
		fmt.Println("\n▶️ Pour lancer l'environnement: turbotilt up")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags pour la commande init
	initCmd.Flags().StringVarP(&forceFramework, "framework", "f", "", "Spécifier manuellement le framework (spring, quarkus, java)")
	initCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port à exposer pour l'application")
	initCmd.Flags().StringVarP(&jdkVersion, "jdk", "j", "17", "Version du JDK à utiliser")
	initCmd.Flags().BoolVarP(&devMode, "dev", "d", true, "Activer les configurations de développement")
	initCmd.Flags().BoolVarP(&detectServices, "services", "s", true, "Détecter et configurer les services dépendants (MySQL, PostgreSQL, etc.)")
}
