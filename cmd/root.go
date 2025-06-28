package cmd

import (
	"fmt"
	"os"
	"turbotilt/internal/update"

	"github.com/spf13/cobra"
)

// Version information, set at build time via ldflags
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// Command-line flags
var (
	dryRun    bool
	debugMode bool
	noUpdate  bool // Flag to disable update checks
)

var rootCmd = &cobra.Command{
	Use:     "turbotilt",
	Short:   "Turbotilt CLI",
	Long:    `Turbotilt - Generate and run cloud-native dev environments.`,
	Version: Version,
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Simulate execution without making changes")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug mode with verbose output")
	rootCmd.PersistentFlags().BoolVar(&noUpdate, "no-update", false, "Disable automatic update checks")

	// Custom version template
	rootCmd.SetVersionTemplate(`Turbotilt {{.Version}}
Build: {{if .GitCommit}}{{.GitCommit}} ({{.BuildTime}}){{else}}unknown{{end}}
`)
}

func Execute() {
	// Set version template variables
	rootCmd.SetVersionTemplate(fmt.Sprintf(`Turbotilt %s
Commit: %s
Built: %s
`, Version, GitCommit, BuildTime))

	// Check for updates unless disabled
	if !noUpdate {
		checkForUpdates()
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// checkForUpdates checks if a newer version is available and prints a message
func checkForUpdates() {
	release, hasUpdate := update.CheckForUpdates(Version)
	if hasUpdate && release != nil {
		fmt.Printf("\n📦 A new version of Turbotilt is available: %s (current: %s)\n",
			release.TagName, Version)
		fmt.Printf("Download it at: %s\n\n", release.URL)
	}
}
