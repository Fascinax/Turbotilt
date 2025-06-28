package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"turbotilt/internal/logger"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🛑 Stopping development environment...")

		// Check if Tilt is running
		tiltRunning := exec.Command("powershell.exe", "-Command", "Get-Process | Where-Object { $_.ProcessName -eq 'tilt' }").Run() == nil

		if tiltRunning {
			fmt.Println("⏳ Stopping Tilt...")
			stopTilt := exec.Command("tilt", "down")
			if err := stopTilt.Run(); err != nil {
				fmt.Printf("❌ Error stopping Tilt: %v\n", err)
			} else {
				fmt.Println("✅ Tilt stopped.")
			}
		} else {
			fmt.Println("⏳ Stopping Docker Compose...")
			stopCompose := exec.Command("docker", "compose", "down")
			if err := stopCompose.Run(); err != nil {
				fmt.Printf("❌ Error stopping Docker Compose: %v\n", err)
			} else {
				fmt.Println("✅ Docker Compose stopped.")
			}
		}

		// Offer to clean temporary files
		if cleanupFlag {
			fmt.Println("⏳ Cleaning temporary files...")
			filesToClean := []string{"Dockerfile", "docker-compose.yml", "Tiltfile"}

			for _, file := range filesToClean {
				cleanFile := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Remove-Item -Path '%s' -ErrorAction SilentlyContinue", file))
				if err := cleanFile.Run(); err != nil {
					logger.Debug("Failed to clean file %s: %v", file, err)
				}
			}

			fmt.Println("✅ Temporary files cleaned.")
		}

		fmt.Println("✨ Environment stopped.")
	},
}

var cleanupFlag bool

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().BoolVarP(&cleanupFlag, "cleanup", "c", false, "Clean generated files (Dockerfile, docker-compose.yml, Tiltfile)")
}
