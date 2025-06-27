package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "turbotilt",
	Short: "Turbotilt CLI",
	Long:  `Turbotilt - Génère et lance des environnements de dev cloud-native.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
