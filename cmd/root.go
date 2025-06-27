package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
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
