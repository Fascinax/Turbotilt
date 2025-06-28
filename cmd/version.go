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
	Short: "Display version information",
	Long:  `Display detailed information about the Turbotilt version, including version number, git commit, and build information.`,
	Run: func(cmd *cobra.Command, args []string) {
		if shortVersion {
			fmt.Printf("Turbotilt %s\n", Version)
			return
		}

		// Use tabwriter for clean formatting
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "🚀 Turbotilt - CLI for cloud-native dev environments")
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
	versionCmd.Flags().BoolVarP(&shortVersion, "short", "s", false, "Display only the version number")
}
