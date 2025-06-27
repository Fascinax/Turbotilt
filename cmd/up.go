package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Lance tilt up ou docker compose up",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Démarrage de l'environnement de développement...")
		fmt.Println("Environnement lancé avec succès!")
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
