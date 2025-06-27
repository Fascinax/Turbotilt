package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"turbotilt/internal/render"
	"turbotilt/internal/scan"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scan et génère Tiltfile & Compose",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Initialisation de Turbotilt...")

		// Détecter le framework
		framework, err := scan.DetectFramework()
		if err != nil {
			fmt.Printf("❌ Erreur lors de la détection du framework: %v\n", err)
			return
		}

		fmt.Printf("✅ Framework détecté: %s\n", framework)

		// Générer les fichiers
		if err := render.GenerateDockerfile(framework); err != nil {
			fmt.Printf("❌ Erreur lors de la génération du Dockerfile: %v\n", err)
			return
		}

		if err := render.GenerateCompose(); err != nil {
			fmt.Printf("❌ Erreur lors de la génération du docker-compose.yml: %v\n", err)
			return
		}

		if err := render.GenerateTiltfile(framework); err != nil {
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
}
