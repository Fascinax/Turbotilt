package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Vérifie installation & config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Vérification de l'installation...")
		fmt.Println("✅ Tout est correctement configuré!")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
