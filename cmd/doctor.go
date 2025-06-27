package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Vérifie installation & config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Vérification de l'installation...")

		// Vérifier Docker
		if err := checkCommand("docker", "--version"); err != nil {
			fmt.Println("❌ Docker n'est pas installé ou n'est pas dans le PATH")
		} else {
			fmt.Println("✅ Docker est installé")
		}

		// Vérifier Tilt
		if err := checkCommand("tilt", "version"); err != nil {
			fmt.Println("❌ Tilt n'est pas installé ou n'est pas dans le PATH")
			fmt.Println("   👉 Pour installer Tilt: https://docs.tilt.dev/install.html")
		} else {
			fmt.Println("✅ Tilt est installé")
		}

		// Vérifier si les fichiers requis existent
		fileCheck := map[string]string{
			"Dockerfile":        "❓ Dockerfile non trouvé. Exécuter `turbotilt init`.",
			"docker-compose.yml": "❓ docker-compose.yml non trouvé. Exécuter `turbotilt init`.",
			"Tiltfile":          "❓ Tiltfile non trouvé. Exécuter `turbotilt init`.",
		}

		for file, message := range fileCheck {
			if _, err := exec.Command("powershell.exe", "-Command", fmt.Sprintf("Test-Path '%s'", file)).Output(); err != nil {
				fmt.Println(message)
			}
		}

		fmt.Println("\n🔧 Pour initialiser le projet: turbotilt init")
		fmt.Println("▶️ Pour démarrer l'environnement: turbotilt up")
	},
}

func checkCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
