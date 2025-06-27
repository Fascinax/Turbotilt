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
	Short: "Lance tilt up ou docker compose up",
	Run: func(cmd *cobra.Command, args []string) {
		if debugMode {
			logger.SetLevel(logger.DEBUG)
			logger.Debug("Mode debug activé")
		}
		
		if dryRun {
			fmt.Println("🔍 Mode simulation (dry-run) activé - aucune modification ne sera appliquée")
		}
		
		fmt.Println("🚀 Démarrage de l'environnement de développement...")

		// Définir les options d'exécution
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
	upCmd.Flags().StringVarP(&serviceName, "service", "s", "", "Démarrer un service spécifique du manifeste (compatible avec les projets multi-services)")
	upCmd.Flags().BoolVarP(&detached, "detach", "d", false, "Exécution en arrière-plan (uniquement pour Docker Compose)")
}
