package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"turbotilt/internal/config"
	"turbotilt/internal/logger"
	"turbotilt/internal/scan"
)

var (
	configPath  string
	projectName string
	projectDesc string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Turbotilt configuration",
	Long:  `Manage the Turbotilt project configuration (turbotilt.yml)`,
}

var initConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Initializing configuration file...")

		// Determine framework
		framework := forceFramework
		var err error

		if framework == "" {
			framework, err = scan.DetectFramework()
			if err != nil {
				logger.Error("Error detecting framework: %v", err)
				framework = "unknown"
			}
		}

		// Create default configuration
		cfg := config.DefaultConfig(framework)

		// Update with command line parameters
		if projectName != "" {
			cfg.Project.Name = projectName
		}

		if projectDesc != "" {
			cfg.Project.Description = projectDesc
		}

		cfg.Docker.Port = port
		cfg.Framework.JdkVersion = jdkVersion
		cfg.Development.EnableLiveReload = devMode

		// Save configuration
		if err := config.SaveConfig(cfg, configPath); err != nil {
			logger.Error("Error saving configuration: %v", err)
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		logger.Info("Configuration saved to %s", configPath)
		fmt.Printf("✅ Configuration saved to %s\n", configPath)
		fmt.Println("📋 Content:")
		fmt.Printf("   - Project: %s\n", cfg.Project.Name)
		fmt.Printf("   - Framework: %s (JDK %s)\n", cfg.Framework.Type, cfg.Framework.JdkVersion)
		fmt.Printf("   - Port: %s\n", cfg.Docker.Port)
		fmt.Printf("   - Live reload: %v\n", cfg.Development.EnableLiveReload)

		// Suggest next step
		fmt.Println("\n▶️ To generate files: turbotilt init")
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Displaying configuration from %s", configPath)

		// Check if file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			logger.Error("Configuration file does not exist: %s", configPath)
			fmt.Printf("❌ Configuration file does not exist: %s\n", configPath)
			fmt.Println("▶️ To create a configuration: turbotilt config init")
			return
		}

		// Load configuration
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			logger.Error("Error loading configuration: %v", err)
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Println("📋 Project Configuration:")
		fmt.Printf("   - Name: %s\n", cfg.Project.Name)
		fmt.Printf("   - Description: %s\n", cfg.Project.Description)
		fmt.Printf("   - Version: %s\n", cfg.Project.Version)
		fmt.Println("📦 Framework:")
		fmt.Printf("   - Type: %s\n", cfg.Framework.Type)
		fmt.Printf("   - JDK: %s\n", cfg.Framework.JdkVersion)
		fmt.Println("🐳 Docker:")
		fmt.Printf("   - Port: %s\n", cfg.Docker.Port)
		fmt.Println("🛠️ Development:")
		fmt.Printf("   - Live reload: %v\n", cfg.Development.EnableLiveReload)
		fmt.Printf("   - Sync path: %s\n", cfg.Development.SyncPath)

		// Display services
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

	// Options for config init
	initConfigCmd.Flags().StringVarP(&configPath, "output", "o", "turbotilt.yml", "Configuration file path")
	initConfigCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
	initConfigCmd.Flags().StringVarP(&projectDesc, "description", "D", "", "Project description")
	initConfigCmd.Flags().StringVarP(&forceFramework, "framework", "f", "", "Framework (spring, quarkus, java)")
	initConfigCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to expose")
	initConfigCmd.Flags().StringVarP(&jdkVersion, "jdk", "j", "17", "JDK version")
	initConfigCmd.Flags().BoolVarP(&devMode, "dev", "d", true, "Development mode")

	// Options for config show
	showConfigCmd.Flags().StringVarP(&configPath, "file", "f", "turbotilt.yml", "Configuration file path")
}
