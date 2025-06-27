package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	shortVersion bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Affiche les informations de version",
	Long:  `Affiche les informations détaillées sur la version de Turbotilt, incluant la version, le commit git, et les informations de build.`,
	Run: func(cmd *cobra.Command, args []string) {
		if shortVersion {
			fmt.Printf("Turbotilt %s\n", Version)
			return
		}

		// Utiliser tabwriter pour un formatage propre
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "🚀 Turbotilt - CLI pour environnements dev cloud-native")
		fmt.Fprintln(w, strings.Repeat("─", 50))
		fmt.Fprintf(w, "Version:\t%s\n", Version)
		fmt.Fprintf(w, "Build time:\t%s\n", BuildTime)
		fmt.Fprintf(w, "Git commit:\t%s\n", GitCommit)
		fmt.Fprintf(w, "Go version:\t%s\n", runtime.Version())
		fmt.Fprintf(w, "OS/Arch:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Fprintln(w, strings.Repeat("─", 50))
		fmt.Fprintln(w, "📦 Homepage: https://github.com/Fascinax/turbotilt")
		fmt.Fprintln(w, "")

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&shortVersion, "short", "s", false, "Affiche uniquement le numéro de version")
}
