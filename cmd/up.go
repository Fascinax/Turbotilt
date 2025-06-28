package cmd

import (
	"fmt"
	"turbotilt/internal/logger"
	"turbotilt/internal/runtime"

	"github.com/spf13/cobra"
)

var (
	useTilt     bool
	detached    bool
	serviceName string
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run tilt up or docker compose up",
	Run: func(cmd *cobra.Command, args []string) {
		if debugMode {
			logger.SetLevel(logger.DEBUG)
			logger.Debug("Debug mode enabled")
		}

		if dryRun {
			fmt.Println("🔍 Simulation mode (dry-run) enabled - no changes will be applied")
		}

		fmt.Println("🚀 Starting development environment...")

		// Define runtime options
		opts := runtime.RunOptions{
			UseTilt:     useTilt,
			Detached:    detached,
			TempFiles:   []string{"Dockerfile", "docker-compose.yml", "Tiltfile"},
			ServiceName: serviceName,
			DryRun:      dryRun,
			Debug:       debugMode,
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
}
