package cmd

import (
	"fmt"
	"turbotilt/internal/runtime"

	"github.com/spf13/cobra"
)

var (
	useTilt  bool
	detached bool
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Lance tilt up ou docker compose up",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Démarrage de l'environnement de développement...")

		// Définir les options d'exécution
		opts := runtime.RunOptions{
			UseTilt:   useTilt,
			Detached:  detached,
			TempFiles: []string{"Dockerfile", "docker-compose.yml", "Tiltfile"},
		}

		var err error
		if useTilt {
			err = runtime.TiltUp(opts)
		} else {
			err = runtime.ComposeUp(opts)
		}

		if err != nil {
			fmt.Printf("❌ Erreur lors du démarrage: %v\n", err)
			return
		}

		if detached {
			fmt.Println("✅ Environnement démarré en arrière-plan.")
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Flags pour la commande up
	upCmd.Flags().BoolVarP(&useTilt, "tilt", "t", true, "Utiliser Tilt (si false, utilise Docker Compose)")
	upCmd.Flags().BoolVarP(&detached, "detach", "d", false, "Exécution en arrière-plan (uniquement pour Docker Compose)")
}
