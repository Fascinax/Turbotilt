package cmd

import (
	"fmt"
	"os"
	"turbotilt/internal/config"
	"turbotilt/internal/logger"
	"turbotilt/internal/scan"

	"github.com/spf13/cobra"
)

var (
	configPath  string
	projectName string
	projectDesc string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Gérer la configuration Turbotilt",
	Long:  `Gérer la configuration du projet Turbotilt (turbotilt.yml)`,
}

var initConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialiser le fichier de configuration",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Initialisation du fichier de configuration...")

		// Déterminer le framework
		framework := forceFramework
		var err error

		if framework == "" {
			framework, err = scan.DetectFramework()
			if err != nil {
				logger.Error("Erreur lors de la détection du framework: %v", err)
				framework = "unknown"
			}
		}

		// Créer une configuration par défaut
		cfg := config.DefaultConfig(framework)

		// Mettre à jour avec les paramètres de la ligne de commande
		if projectName != "" {
			cfg.Project.Name = projectName
		}

		if projectDesc != "" {
			cfg.Project.Description = projectDesc
		}

		cfg.Docker.Port = port
		cfg.Framework.JdkVersion = jdkVersion
		cfg.Development.EnableLiveReload = devMode

		// Enregistrer la configuration
		if err := config.SaveConfig(cfg, configPath); err != nil {
			logger.Error("Erreur lors de l'enregistrement de la configuration: %v", err)
			fmt.Printf("❌ Erreur: %v\n", err)
			return
		}

		logger.Info("Configuration enregistrée dans %s", configPath)
		fmt.Printf("✅ Configuration enregistrée dans %s\n", configPath)
		fmt.Println("📋 Contenu:")
		fmt.Printf("   - Projet: %s\n", cfg.Project.Name)
		fmt.Printf("   - Framework: %s (JDK %s)\n", cfg.Framework.Type, cfg.Framework.JdkVersion)
		fmt.Printf("   - Port: %s\n", cfg.Docker.Port)
		fmt.Printf("   - Live reload: %v\n", cfg.Development.EnableLiveReload)

		// Proposer la prochaine étape
		fmt.Println("\n▶️ Pour générer les fichiers: turbotilt init")
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Afficher la configuration actuelle",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Affichage de la configuration depuis %s", configPath)

		// Vérifier si le fichier existe
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			logger.Error("Le fichier de configuration n'existe pas: %s", configPath)
			fmt.Printf("❌ Le fichier de configuration n'existe pas: %s\n", configPath)
			fmt.Println("▶️ Pour créer une configuration: turbotilt config init")
			return
		}

		// Charger la configuration
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			logger.Error("Erreur lors du chargement de la configuration: %v", err)
			fmt.Printf("❌ Erreur: %v\n", err)
			return
		}

		fmt.Println("📋 Configuration du projet:")
		fmt.Printf("   - Nom: %s\n", cfg.Project.Name)
		fmt.Printf("   - Description: %s\n", cfg.Project.Description)
		fmt.Printf("   - Version: %s\n", cfg.Project.Version)
		fmt.Println("📦 Framework:")
		fmt.Printf("   - Type: %s\n", cfg.Framework.Type)
		fmt.Printf("   - JDK: %s\n", cfg.Framework.JdkVersion)
		fmt.Println("🐳 Docker:")
		fmt.Printf("   - Port: %s\n", cfg.Docker.Port)
		fmt.Println("🛠️ Développement:")
		fmt.Printf("   - Live reload: %v\n", cfg.Development.EnableLiveReload)
		fmt.Printf("   - Chemin de synchronisation: %s\n", cfg.Development.SyncPath)

		// Afficher les services
		if len(cfg.Services) > 0 {
			fmt.Println("🔌 Services:")
			for _, svc := range cfg.Services {
				fmt.Printf("   - %s (%s:%s)\n", svc.Name, svc.Type, svc.Version)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(initConfigCmd)
	configCmd.AddCommand(showConfigCmd)

	// Options pour config init
	initConfigCmd.Flags().StringVarP(&configPath, "output", "o", "turbotilt.yml", "Chemin du fichier de configuration")
	initConfigCmd.Flags().StringVarP(&projectName, "name", "n", "", "Nom du projet")
	initConfigCmd.Flags().StringVarP(&projectDesc, "description", "D", "", "Description du projet")
	initConfigCmd.Flags().StringVarP(&forceFramework, "framework", "f", "", "Framework (spring, quarkus, java)")
	initConfigCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port à exposer")
	initConfigCmd.Flags().StringVarP(&jdkVersion, "jdk", "j", "17", "Version du JDK")
	initConfigCmd.Flags().BoolVarP(&devMode, "dev", "d", true, "Mode développement")

	// Options pour config show
	showConfigCmd.Flags().StringVarP(&configPath, "file", "f", "turbotilt.yml", "Chemin du fichier de configuration")
}
