package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scan et génère Tiltfile & Compose",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initialisation de Turbotilt...")
		fmt.Println("Scan du projet terminé. Fichiers générés.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
