package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"turbotilt/internal/config"
	"turbotilt/internal/render"
	"turbotilt/internal/scan"
)

var (
	// Options pour la commande init
	forceFramework   string
	port             string
	jdkVersion       string
	devMode          bool
	detectServices   bool
	generateManifest bool
	fromManifest     bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scan et génère Tiltfile & Compose",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Initialisation de Turbotilt...")

		// Rechercher un manifeste existant
		configPath, isManifest, _ := config.FindConfiguration()

		// Si --from-manifest est demandé ou si un manifeste existe et --generate-manifest n'est pas demandé
		if fromManifest || (isManifest && !generateManifest) {
			if configPath == "" {
				fmt.Println("❌ Aucun manifeste trouvé. Utilisez --generate-manifest pour en créer un.")
				return
			}

			fmt.Printf("📄 Utilisation du manifeste %s\n", configPath)
			manifest, err := config.LoadManifest(configPath)
			if err != nil {
				fmt.Printf("❌ Erreur lors du chargement du manifeste: %v\n", err)
				return
			}

			fmt.Printf("✅ Manifeste chargé avec %d service(s)\n", len(manifest.Services))

			// Convertir les services du manifeste en options de rendu
			serviceList := render.ServiceList{
				Services: []render.Options{},
			}

			for _, service := range manifest.Services {
				// Ignorer les services dépendants (sans runtime)
				if service.Runtime == "" {
					continue
				}

				opts, err := config.ConvertManifestToRenderOptions(service)
				if err != nil {
					fmt.Printf("⚠️ Avertissement: %v\n", err)
					continue
				}

				serviceList.Services = append(serviceList.Services, *opts)
			}

			// Générer les fichiers pour un projet multi-services
			if len(serviceList.Services) > 0 {
				fmt.Println("🔧 Génération des configurations pour un projet multi-services...")

				if err := render.GenerateMultiServiceCompose(serviceList); err != nil {
					fmt.Printf("❌ Erreur lors de la génération du docker-compose.yml: %v\n", err)
					return
				}

				if err := render.GenerateMultiServiceTiltfile(serviceList); err != nil {
					fmt.Printf("❌ Erreur lors de la génération du Tiltfile: %v\n", err)
					return
				}

				fmt.Println("✨ Configuration Turbotilt terminée!")
				fmt.Println("📋 Fichiers générés à partir du manifeste:")
				fmt.Println("   - docker-compose.yml")
				fmt.Println("   - Tiltfile")
				fmt.Println("\n▶️ Pour lancer l'environnement: turbotilt up")
				return
			}
		}

		// Si on arrive ici, on procède avec l'auto-détection ou les options CLI

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
			ServiceName: appName, // Utiliser le nom pour l'identifier dans un contexte multi-services
			Framework:   framework,
			AppName:     appName,
			Port:        port,
			JDKVersion:  jdkVersion,
			DevMode:     devMode,
			Path:        ".",
			Services:    services,
		}

		// Générer le manifest si demandé
		if generateManifest {
			fmt.Println("📝 Génération du manifeste turbotilt.yaml...")

			// Créer une configuration basée sur les résultats de la détection
			cfg := config.Config{
				Project: config.ProjectConfig{
					Name:        appName,
					Description: "Projet Turbotilt",
					Version:     "1.0.0",
				},
				Framework: config.FrameworkConfig{
					Type:       framework,
					JdkVersion: jdkVersion,
				},
				Docker: config.DockerConfig{
					Port: port,
				},
				Development: config.DevelopmentConfig{
					EnableLiveReload: devMode,
				},
				Services: []config.ServiceConfig{},
			}

			// Convertir les services scan.ServiceConfig en config.ServiceConfig
			for _, svc := range services {
				// Générer un nom basé sur le type
				serviceName := strings.ToLower(string(svc.Type))

				configSvc := config.ServiceConfig{
					Name:        serviceName,
					Type:        string(svc.Type),
					Version:     svc.Version,
					Port:        svc.Port,
					Environment: svc.Credentials,
				}
				cfg.Services = append(cfg.Services, configSvc)
			}

			// Générer le manifeste à partir de la configuration
			manifest := config.GenerateManifestFromConfig(cfg)

			// Enregistrer le manifeste
			if err := config.SaveManifest(manifest, config.ManifestFileName); err != nil {
				fmt.Printf("❌ Erreur lors de l'enregistrement du manifeste: %v\n", err)
			} else {
				fmt.Printf("✅ Manifeste %s généré avec succès!\n", config.ManifestFileName)
			}
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
		if generateManifest {
			fmt.Printf("   - %s\n", config.ManifestFileName)
		}
		fmt.Println("\n▶️ Pour lancer l'environnement: turbotilt up")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags pour la commande init
	initCmd.Flags().StringVarP(&forceFramework, "framework", "f", "", "Spécifier manuellement le framework (spring, quarkus, java)")
	initCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port à exposer pour l'application")
	initCmd.Flags().StringVarP(&jdkVersion, "jdk", "j", "11", "Version du JDK à utiliser")
	initCmd.Flags().BoolVarP(&devMode, "dev", "d", true, "Activer les configurations de développement")
	initCmd.Flags().BoolVarP(&detectServices, "services", "s", true, "Détecter et configurer les services dépendants (MySQL, PostgreSQL, etc.)")
	initCmd.Flags().BoolVarP(&generateManifest, "generate-manifest", "g", false, "Générer un manifeste turbotilt.yaml à partir de la détection")
	initCmd.Flags().BoolVarP(&fromManifest, "from-manifest", "m", false, "Initialiser le projet à partir d'un manifeste existant")
}
