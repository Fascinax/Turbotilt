package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"turbotilt/internal/logger"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Arrête l'environnement de développement",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🛑 Arrêt de l'environnement de développement...")

		// Vérifier si Tilt est en cours d'exécution
		tiltRunning := exec.Command("powershell.exe", "-Command", "Get-Process | Where-Object { $_.ProcessName -eq 'tilt' }").Run() == nil

		if tiltRunning {
			fmt.Println("⏳ Arrêt de Tilt...")
			stopTilt := exec.Command("tilt", "down")
			if err := stopTilt.Run(); err != nil {
				fmt.Printf("❌ Erreur lors de l'arrêt de Tilt: %v\n", err)
			} else {
				fmt.Println("✅ Tilt arrêté.")
			}
		} else {
			fmt.Println("⏳ Arrêt de Docker Compose...")
			stopCompose := exec.Command("docker", "compose", "down")
			if err := stopCompose.Run(); err != nil {
				fmt.Printf("❌ Erreur lors de l'arrêt de Docker Compose: %v\n", err)
			} else {
				fmt.Println("✅ Docker Compose arrêté.")
			}
		}

		// Proposer de nettoyer les fichiers temporaires
		if cleanupFlag {
			fmt.Println("⏳ Nettoyage des fichiers temporaires...")
			filesToClean := []string{"Dockerfile", "docker-compose.yml", "Tiltfile"}

			for _, file := range filesToClean {
				cleanFile := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Remove-Item -Path '%s' -ErrorAction SilentlyContinue", file))
				if err := cleanFile.Run(); err != nil {
					logger.Debug("Failed to clean file %s: %v", file, err)
				}
			}

			fmt.Println("✅ Fichiers temporaires nettoyés.")
		}

		fmt.Println("✨ Environnement arrêté.")
	},
}

var cleanupFlag bool

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().BoolVarP(&cleanupFlag, "cleanup", "c", false, "Nettoyer les fichiers générés (Dockerfile, docker-compose.yml, Tiltfile)")
}
