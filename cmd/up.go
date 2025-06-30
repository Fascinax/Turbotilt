package cmd

import (
	"fmt"
	"turbotilt/internal/config"
	"turbotilt/internal/i18n"
	"turbotilt/internal/logger"
	"turbotilt/internal/runtime"

	"github.com/spf13/cobra"
)

var (
	useTilt     bool
	detached    bool
	serviceName string
	configFile  string
	useMemory   bool
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run tilt up or docker compose up",
	Run: func(cmd *cobra.Command, args []string) {
		t := i18n.GetTranslator()
		log := logger.GetLogger()
		
		if debugMode {
			logger.SetLevel(logger.DEBUG)
			log.Debug("Debug mode enabled")
		}

		if dryRun {
			log.Info("🔍 Simulation mode (dry-run) enabled - no changes will be applied")
		}

		log.Info("🚀 Starting development environment...")
		
		// Vérifier si on doit utiliser la configuration en mémoire
		if useMemory {
			// Vérifier si des services sont stockés en mémoire
			memoryStore := config.GetMemoryStore()
			if !memoryStore.HasSelectedServices() {
				log.Error(t.Tr("No services found in memory. Run 'turbotilt select' first or specify a configuration file."))
				return
			}
			log.Info(t.Tr("Using services from memory selected with 'select' command"))
		}

		// Define runtime options
		opts := runtime.RunOptions{
			UseTilt:     useTilt,
			Detached:    detached,
			TempFiles:   []string{"Dockerfile", "docker-compose.yml", "Tiltfile"},
			ServiceName: serviceName,
			DryRun:      dryRun,
			Debug:       debugMode,
			ConfigFile:  configFile,
			UseMemory:   useMemory,
		}

		var err error
		if useTilt {
			err = runtime.TiltUp(opts)
		} else {
			err = runtime.ComposeUp(opts)
		}

		if err != nil {
			fmt.Printf("❌ Error starting environment: %v\n", err)
			return
		}

		if detached {
			fmt.Println("✅ Environment started in background.")
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Flags for the up command
	upCmd.Flags().BoolVarP(&useTilt, "tilt", "t", true, "Use Tilt (if false, uses Docker Compose)")
	upCmd.Flags().StringVarP(&serviceName, "service", "s", "", "Start a specific service from the manifest (compatible with multi-service projects)")
	upCmd.Flags().BoolVarP(&detached, "detach", "d", false, "Run in background (only for Docker Compose)")
	upCmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to the configuration file (if not specified, uses turbotilt.yaml or memory)")
	upCmd.Flags().BoolVarP(&useMemory, "memory", "m", true, "Use services selected with the select command stored in memory")
}
