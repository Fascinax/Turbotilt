package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"turbotilt/internal/logger"
	"turbotilt/internal/runtime"

	"github.com/spf13/cobra"
)

var tupCmd = &cobra.Command{
	Use:   "tup",
	Short: "Temporary up - run services with temporary configurations",
	Long:  `Generates configuration files, starts services, and cleans up when done.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debugMode {
			logger.SetLevel(logger.DEBUG)
			logger.Debug("Debug mode enabled")
		}

		if dryRun {
			fmt.Println("🔍 Simulation mode (dry-run) enabled - no changes will be applied")
		}

		fmt.Println("🔄 Starting temporary development environment...")

		// 1. Run the init command first to generate config files
		fmt.Println("⏳ Initializing configuration...")
		initProcess := exec.Command(os.Args[0], "init")
		initProcess.Stdout = os.Stdout
		initProcess.Stderr = os.Stderr

		if err := initProcess.Run(); err != nil {
			fmt.Printf("❌ Failed to initialize configuration: %v\n", err)
			return
		}
		fmt.Println("✅ Configuration generated successfully.")

		// 2. Run the up command
		fmt.Println("🚀 Starting services...")

		// Set up a channel to handle Ctrl+C
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Start services using the runtime package
		fmt.Println("⏳ Starting services in development mode...")
		if err := runtime.StartEnvironment(serviceName, useTilt, detached, dryRun); err != nil {
			fmt.Printf("❌ Failed to start environment: %v\n", err)
			// Cleanup even if starting fails
			cleanupFiles()
			return
		}

		// If running in detached mode, cleanup immediately after starting
		if detached {
			fmt.Println("✅ Services started in detached mode.")
			fmt.Println("ℹ️ Use 'turbotilt stop' to stop services when done.")
			fmt.Println("🧹 Cleaning up temporary configuration files...")
			cleanupFiles()
			return
		}

		fmt.Println("✅ Services started successfully.")
		fmt.Println("ℹ️ Press Ctrl+C to stop services and clean up temporary files.")

		// Wait for Ctrl+C
		<-sigChan

		// 3. Stop services and clean up
		fmt.Println("\n🛑 Stopping services...")
		stopProcess := exec.Command(os.Args[0], "stop")
		stopProcess.Stdout = os.Stdout
		stopProcess.Stderr = os.Stderr

		if err := stopProcess.Run(); err != nil {
			fmt.Printf("⚠️ Warning: Failed to stop services cleanly: %v\n", err)
		}

		// 4. Clean up generated files
		fmt.Println("🧹 Cleaning up temporary configuration files...")
		cleanupFiles()

		fmt.Println("✨ Temporary environment shutdown complete.")
	},
}

// cleanupFiles removes temporary configuration files
func cleanupFiles() {
	filesToClean := []string{"Dockerfile", "docker-compose.yml", "Tiltfile", "turbotilt.yaml"}

	for _, file := range filesToClean {
		cleanFile := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Remove-Item -Path '%s' -ErrorAction SilentlyContinue", file))
		if err := cleanFile.Run(); err != nil {
			logger.Debug("Failed to clean file %s: %v", file, err)
		}
	}
}

func init() {
	rootCmd.AddCommand(tupCmd)

	// Reuse the same flags as the 'up' command
	tupCmd.Flags().BoolVar(&useTilt, "tilt", true, "Use Tilt for live reload (default)")
	tupCmd.Flags().BoolVar(&detached, "detached", false, "Run in detached mode")
	tupCmd.Flags().StringVar(&serviceName, "service", "", "Run only a specific service")
}
